# Build and install from source

# Prerequisites

## Go
Install Go 1.2 or newer from [golang.org](http://golang.org/doc/install).

## DVCS Clients
Go will need these DVCS clients installed in order to fetch all of Hockeypuck's
package dependencies:

* Bazaar
* Git
* Mercurial

On Ubuntu:

`sudo apt-get install bzr git mercurial`

# Fetch sources with the packaging scripts

```
git clone https://github.com/hockeypuck/packaging
cd packaging
./prepare.bash
```

# Build hockeypuck executables

```
export GOPATH=$(pwd)
go install github.com/hockeypuck/server/cmd/hockeypuck
go install github.com/hockeypuck/server/cmd/hockeypuck-load
go install github.com/hockeypuck/server/cmd/hockeypuck-pbuild
```

# Install

Copy the executables into the desired target location. The files under
`instroot` in the packaging project may also be useful.

# Next steps

* [Configure](configuration.md) the Hockeypuck server.
* Run Hockeypuck with `/path/to/hockeypuck -config /path/to/hockeypuck.conf`.

