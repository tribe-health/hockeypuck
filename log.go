/*
   Hockeypuck - OpenPGP key server
   Copyright (C) 2012-2014  Casey Marshall

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, version 3.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package hockeypuck

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("hockeypuck")

// Logfile option
func (s *Settings) LogFile() string {
	return s.GetString("hockeypuck.logfile")
}

func (s *Settings) LogSpec() string {
	return s.GetStringDefault("hockeypuck.logspec", `<root>=INFO`)
}

// InitLog initializes the logging output to the globally configured settings.
// It also registers SIGHUP, SIGUSR1 and SIGUSR2 to close and reopen the log file
// for logrotate(8) support.
//
// BUG: If InitLog is called before the application is properly configured, it will automatically
// configure the application with an empty TOML (accept all defaults).
func InitLog() {
	if Config() == nil {
		SetConfig("")
	}
	if Config().LogFile() != "" {
		// Handle signals for log rotation
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)
		go func() {
			for {
				select {
				case _ = <-sigChan:
					openLog()
					logger.Infof("reopened logfile")
				}
			}
		}()
	}
	// Open the log
	openLog()
}

type logWriter struct {
	writer    io.Writer
	formatter loggo.Formatter
}

func (lw *logWriter) Write(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) {
	logLine := lw.formatter.Format(level, module, filename, line, timestamp, message)
	fmt.Fprintln(lw.writer, logLine)
}

func (lw *logWriter) Close() error {
	if cl, ok := lw.writer.(io.WriteCloser); ok {
		return cl.Close()
	}
	return nil
}

func setOutput(w io.Writer) error {
	newWriter := &logWriter{w, &loggo.DefaultFormatter{}}
	prevWriter, err := loggo.ReplaceDefaultWriter(newWriter)
	if err != nil {
		return err
	}
	if lw, ok := prevWriter.(*logWriter); ok {
		lw.Close()
	}
	return nil
}

func openLog() {
	if Config().LogFile() != "" {
		logOut, err := os.OpenFile(Config().LogFile(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			logger.Errorf("failed to open logfile: %v", err)
			logOut = os.Stderr
		} else {
			err = setOutput(logOut)
			if err != nil {
				logger.Errorf("failed to set logfile output: %v", err)
				logOut = os.Stderr
			}
		}
	} else {
		err := setOutput(os.Stderr)
		if err != nil {
			logger.Errorf("failed to set logfile output: %v", err)
		}
	}
	err := loggo.ConfigureLoggers(Config().LogSpec())
	if err != nil {
		logger.Errorf("invalid logger spec %q: %v", Config().LogSpec(), err)
	}
}
