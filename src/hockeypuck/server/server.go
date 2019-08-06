package server

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/carbocation/interpose"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/errgo.v1"
	"gopkg.in/tomb.v2"

	"hockeypuck/hkp"
	"hockeypuck/hkp/sks"
	"hockeypuck/hkp/storage"
	log "hockeypuck/logrus"
	"hockeypuck/metrics"
	"hockeypuck/mgohkp"
	"hockeypuck/pghkp"
)

type Server struct {
	settings        *Settings
	st              storage.Storage
	middle          *interpose.Middleware
	r               *httprouter.Router
	sksPeer         *sks.Peer
	logWriter       io.WriteCloser
	metricsListener *metrics.Metrics

	t                 tomb.Tomb
	hkpAddr, hkpsAddr string
}

type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewStatusCodeResponseWriter(w http.ResponseWriter) *statusCodeResponseWriter {
	// WriteHeader is not called if our response implicitly
	// returns 200 OK, so we default to that status code.
	return &statusCodeResponseWriter{w, http.StatusOK}
}

func (scrw *statusCodeResponseWriter) WriteHeader(code int) {
	scrw.statusCode = code
	scrw.ResponseWriter.WriteHeader(code)
}

func NewServer(settings *Settings) (*Server, error) {
	if settings == nil {
		defaults := DefaultSettings()
		settings = &defaults
	}
	s := &Server{
		settings: settings,
		r:        httprouter.New(),
	}

	var err error
	s.st, err = DialStorage(settings)
	if err != nil {
		return nil, err
	}

	s.middle = interpose.New()
	s.middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			start := time.Now()
			scrw := NewStatusCodeResponseWriter(rw)
			next.ServeHTTP(scrw, req)
			duration := time.Since(start)
			fields := log.Fields{
				req.Method:    req.URL.String(),
				"duration":    duration.String(),
				"from":        req.RemoteAddr,
				"host":        req.Host,
				"status-code": scrw.statusCode,
				"user-agent":  req.UserAgent(),
			}
			proxyHeaders := []string{
				"x-forwarded-for",
				"x-forwarded-host",
				"x-forwarded-server",
			}
			for _, ph := range proxyHeaders {
				if v := req.Header.Get(ph); v != "" {
					fields[ph] = v
				}
			}
			log.WithFields(fields).Info()
			recordHTTPRequestDuration(req.Method, scrw.statusCode, duration)
		})
	})
	s.middle.UseHandler(s.r)

	s.sksPeer, err = sks.NewPeer(s.st, settings.Conflux.Recon.LevelDB.Path, &settings.Conflux.Recon.Settings)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	s.metricsListener = metrics.NewMetrics(&settings.Metrics)

	options := []hkp.HandlerOption{hkp.StatsFunc(s.stats)}
	if settings.IndexTemplate != "" {
		options = append(options, hkp.IndexTemplate(settings.IndexTemplate))
	}
	if settings.VIndexTemplate != "" {
		options = append(options, hkp.VIndexTemplate(settings.VIndexTemplate))
	}
	if settings.StatsTemplate != "" {
		options = append(options, hkp.StatsTemplate(settings.StatsTemplate))
	}
	h, err := hkp.NewHandler(s.st, options...)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	h.Register(s.r)

	if settings.Webroot != "" {
		err := s.registerWebroot(settings.Webroot)
		if err != nil {
			return nil, errgo.Mask(err)
		}
	}

	registerMetrics()
	s.st.Subscribe(metricsStorageNotifier)

	return s, nil
}

func DialStorage(settings *Settings) (storage.Storage, error) {
	switch settings.OpenPGP.DB.Driver {
	case "mongo":
		var options []mgohkp.Option
		if settings.OpenPGP.DB.Mongo != nil {
			if settings.OpenPGP.DB.Mongo.DB != "" {
				options = append(options, mgohkp.DBName(settings.OpenPGP.DB.Mongo.DB))
			}
			if settings.OpenPGP.DB.Mongo.Collection != "" {
				options = append(options, mgohkp.CollectionName(settings.OpenPGP.DB.Mongo.Collection))
			}
		}
		return mgohkp.Dial(settings.OpenPGP.DB.DSN, options...)
	case "postgres-jsonb":
		return pghkp.Dial(settings.OpenPGP.DB.DSN)
	}
	return nil, errgo.Newf("storage driver %q not supported", settings.OpenPGP.DB.Driver)
}

type stats struct {
	Now       string      `json:"now"`
	Version   string      `json:"version"`
	Hostname  string      `json:"hostname"`
	Nodename  string      `json:"nodename"`
	Contact   string      `json:"server_contact"`
	HTTPAddr  string      `json:"httpAddr"`
	ReconAddr string      `json:"reconAddr"`
	Software  string      `json:"software"`
	Peers     []statsPeer `json:"peers"`

	Total  int
	Hourly []loadStat
	Daily  []loadStat
}

type loadStat struct {
	*sks.LoadStat
	Time time.Time
}

type loadStats []loadStat

func (s loadStats) Len() int           { return len(s) }
func (s loadStats) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s loadStats) Less(i, j int) bool { return s[i].Time.Before(s[j].Time) }

type statsPeer struct {
	Name      string
	HTTPAddr  string `json:"httpAddr"`
	ReconAddr string `json:"reconAddr"`
}

type statsPeers []statsPeer

func (s statsPeers) Len() int           { return len(s) }
func (s statsPeers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s statsPeers) Less(i, j int) bool { return s[i].Name < s[j].Name }

func (s *Server) stats() (interface{}, error) {
	sksStats := s.sksPeer.Stats()

	result := &stats{
		Now:       time.Now().UTC().Format(time.RFC3339),
		Version:   s.settings.Version,
		Contact:   s.settings.Contact,
		HTTPAddr:  s.settings.HKP.Bind,
		ReconAddr: s.settings.Conflux.Recon.Settings.ReconAddr,
		Software:  s.settings.Software,

		Total: sksStats.Total,
	}

	nodename, err := os.Hostname()
	if err != nil {
		log.Warningf("cannot determine local hostname: %v", err)
	} else {
		result.Nodename = nodename
	}

	if s.settings.Hostname != "" {
		result.Hostname = s.settings.Hostname
	} else if nodename != "" {
		result.Hostname = nodename
	}

	for k, v := range sksStats.Hourly {
		result.Hourly = append(result.Hourly, loadStat{LoadStat: v, Time: k})
	}
	sort.Sort(loadStats(result.Hourly))
	for k, v := range sksStats.Daily {
		result.Daily = append(result.Daily, loadStat{LoadStat: v, Time: k})
	}
	sort.Sort(loadStats(result.Daily))
	for k, v := range s.settings.Conflux.Recon.Settings.Partners {
		result.Peers = append(result.Peers, statsPeer{
			Name:      k,
			HTTPAddr:  v.HTTPAddr,
			ReconAddr: v.ReconAddr,
		})
	}
	sort.Sort(statsPeers(result.Peers))
	return result, nil
}

func (s *Server) registerWebroot(webroot string) error {
	fileServer := http.FileServer(http.Dir(webroot))
	d, err := os.Open(webroot)
	if os.IsNotExist(err) {
		log.Errorf("webroot %q not found", webroot)
		// non-fatal error
		return nil
	} else if err != nil {
		return errgo.Mask(err)
	}
	defer d.Close()
	files, err := d.Readdir(0)
	if err != nil {
		return errgo.Mask(err)
	}

	s.r.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		fileServer.ServeHTTP(w, req)
	})
	// httprouter needs explicit paths, so we need to set up a route for each
	// path. This will panic if there are any paths that conflict with
	// previously registered routes.
	for _, fi := range files {
		name := fi.Name()
		if !fi.IsDir() {
			s.r.GET("/"+name, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
				req.URL.Path = "/" + name
				fileServer.ServeHTTP(w, req)
			})
		} else {
			s.r.GET("/"+name+"/*filepath", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
				req.URL.Path = "/" + name + ps.ByName("filepath")
				fileServer.ServeHTTP(w, req)
			})
		}
	}
	return nil
}

func (s *Server) Start() error {
	s.openLog()

	s.t.Go(s.listenAndServeHKP)
	if s.settings.HKPS != nil {
		s.t.Go(s.listenAndServeHKPS)
	}

	if s.sksPeer != nil {
		s.sksPeer.Start()
	}

	if s.metricsListener != nil {
		s.metricsListener.Start()
	}

	return nil
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func (s *Server) openLog() {
	defer func() {
		level, err := log.ParseLevel(strings.ToLower(s.settings.LogLevel))
		if err != nil {
			log.Warningf("invalid LogLevel=%q: %v", s.settings.LogLevel, err)
			return
		}
		log.SetLevel(level)
	}()

	s.logWriter = nopCloser{os.Stderr}
	if s.settings.LogFile != "" {
		f, err := os.OpenFile(s.settings.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Errorf("failed to open LogFile=%q: %v", s.settings.LogFile, err)
		}
		s.logWriter = f
	}
	log.SetOutput(s.logWriter)
	log.Debug("log opened")
}

func (s *Server) closeLog() {
	log.SetOutput(os.Stderr)
	s.logWriter.Close()
}

func (s *Server) LogRotate() {
	w := s.logWriter
	s.openLog()
	w.Close()
}

func (s *Server) Wait() error {
	return s.t.Wait()
}

func (s *Server) Stop() {
	defer s.closeLog()

	if s.sksPeer != nil {
		s.sksPeer.Stop()
	}
	s.t.Kill(nil)
	s.t.Wait()
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by listenAndServe and listenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept implements net.Listener.
func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

var newListener = (*Server).newListener

func (s *Server) newListener(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	s.t.Go(func() error {
		<-s.t.Dying()
		return ln.Close()
	})
	return tcpKeepAliveListener{ln.(*net.TCPListener)}, nil
}

func (s *Server) listenAndServeHKP() error {
	ln, err := newListener(s, s.settings.HKP.Bind)
	if err != nil {
		return err
	}
	s.hkpAddr = ln.Addr().String()
	return http.Serve(ln, s.middle)
}

func (s *Server) listenAndServeHKPS() error {
	config := &tls.Config{
		NextProtos: []string{"http/1.1"},
	}
	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(s.settings.HKPS.Cert, s.settings.HKPS.Key)
	if err != nil {
		return errgo.Notef(err, "failed to load HKPS certificate=%q key=%q", s.settings.HKPS.Cert, s.settings.HKPS.Key)
	}

	ln, err := newListener(s, s.settings.HKP.Bind)
	if err != nil {
		return errgo.Mask(err)
	}
	s.hkpsAddr = ln.Addr().String()
	ln = tls.NewListener(ln, config)
	return http.Serve(ln, s.middle)
}
