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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/errgo.v1"

	"github.com/hockeypuck/hockeypuck"
	"github.com/hockeypuck/hockeypuck/openpgp"
)

var (
	configFile = flag.String("config", "", "config file")
	outputDir  = flag.String("path", ".", "output path")
	count      = flag.Int("count", 15000, "keys per file")
)

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, errgo.Details(err))
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	flag.Parse()

	if *configFile != "" {
		if err := hockeypuck.LoadConfigFile(*configFile); err != nil {
			die(err)
		}
	} else {
		// Fall back on default empty config
		hockeypuck.SetConfig("")
	}

	err := dump()
	die(err)
}

func dump() error {
	var db *openpgp.DB
	var err error
	if db, err = openpgp.NewDB(); err != nil {
		return errgo.Mask(err)
	}

	num := 0
	i := 0
	f, err := openNew(num)
	if err != nil {
		return errgo.Mask(err)
	}

	w := &openpgp.Worker{Loader: openpgp.NewLoader(db, false)}
	rows, err := db.Queryx(`SELECT uuid FROM openpgp_pubkey`)
	if err != nil {
		return errgo.Mask(err)
	}
	for rows.Next() {
		var uuid string
		err = rows.Scan(&uuid)
		if err != nil {
			return errgo.Mask(err)
		}

		pubkey, err := w.FetchKey(uuid)
		if err != nil {
			return errgo.Mask(err)
		}

		err = openpgp.WritePackets(f, pubkey)
		if err != nil {
			return errgo.Mask(err)
		}

		i++
		if i >= *count {
			f.Close()
			num++
			f, err = openNew(num)
			if err != nil {
				return errgo.Mask(err)
			}
			i = 0
		}
	}
	f.Close()
	return nil
}

func openNew(num int) (*os.File, error) {
	return os.Create(filepath.Join(*outputDir, fmt.Sprintf("hkp-dump-%04d.pgp", num)))
}
