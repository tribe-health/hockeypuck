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

package openpgp

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errgo"
	"github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func NewDB() (db *DB, err error) {
	db = new(DB)
	db.DB, err = sqlx.Connect(Config().Driver(), Config().DSN())
	if err != nil {
		err = errgo.Mask(err)
	}
	return db, err
}

func (db *DB) CreateSchema() (err error) {
	if err = db.CreateTables(); err != nil {
		return err
	}
	return db.CreateConstraints()
}

func (db *DB) CreateTables() (err error) {
	for _, crSql := range CreateTablesSql {
		logger.Tracef(crSql)
		_, err = db.Exec(crSql)
		if err != nil {
			return errgo.NoteMask(err, crSql)
		}
	}
	return
}

func (db *DB) DeleteDuplicates() (err error) {
	for _, sql := range DeleteDuplicatesSql {
		logger.Tracef(sql)
		if _, err = db.Exec(sql); err != nil {
			return errgo.NoteMask(err, sql)
		}
	}
	return
}

func isDuplicate(err error) bool {
	if pgerr, is := err.(pq.PGError); is {
		switch pgerr.Get('C') {
		case "23000":
			return true
		case "23505":
			return true
		}
	}
	return false
}

func isDuplicateConstraint(err error) bool {
	if pgerr, is := err.(pq.PGError); is {
		switch pgerr.Get('C') {
		case "42P16":
			return true
		case "42P07":
			return true
		case "42P10":
			return true
		case "42710":
			return true
		}
	}
	return false
}

func (db *DB) CreateConstraints() (err error) {
	for _, crSqls := range CreateConstraintsSql {
		for _, crSql := range crSqls {
			if _, err = db.Exec(crSql); err != nil {
				if isDuplicateConstraint(err) {
					err = nil
				} else {
					return errgo.NoteMask(err, crSql)
				}
			}
		}
	}
	return
}

func (db *DB) DropConstraints() (err error) {
	for _, drSqls := range DropConstraintsSql {
		for _, drSql := range drSqls {
			if _, err := db.Exec(drSql); err != nil {
				logger.Warningf("%s: %v", drSql, err)
			}
		}
	}
	return nil
}
