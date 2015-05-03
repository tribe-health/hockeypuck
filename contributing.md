# Contributing

## Where to contribute

Hockeypuck has been split into several Go package projects.

In general, all Hockeypuck projects are maintained under the
[hockeypuck](https://github.com/hockeypuck) organization, in the following
Github projects:

### Hockeypuck Packages

Hockeypuck is composed of several small Go packages, that attempt to do one thing well.

* [conflux](https://gopkg.in/hockeypuck/conflux.v2) Reconciliation protocol used for peering.
* [hkp](https://gopkg.in/hockeypuck/hkp.v1) HKP protocol handler.
* [mgohkp](https://gopkg.in/hockeypuck/mgohkp.v1) MongoDB storage driver.
* [openpgp](https://gopkg.in/hockeypuck/openpgp.v1) OpenPGP public key data model & processing. 
* [pghkp](https://gopkg.in/hockeypuck/pghkp.v1) PostgreSQL JSONB storage driver.

### Hockeypuck Server

The Hockeypuck server at https://github.com/hockeypuck/server integrates the
above packages with a server configuration model, logging, server and
maintenance utility binaries.

### Hockeypuck Front-End

https://github.com/hockeypuck/webroot is a fork of Matt Rude's
[pgpkeyserver-lite](https://github.com/mattrude/pgpkeyserver-lite) front end.
It's included in the Hockeypuck release for convenience.

### Hockeypuck Packaging

https://github.com/hockeypuck/packaging is a collection of release management
scripts. These scripts fetch the above Hockeypuck source packages and all their
dependencies at known compatible and working versions. Debian packaging as well
as cross-compiled tarballs are supported for release distribution.

### Hockeypuck documentation

This project, https://github.com/hockeypuck/hockeypuck used to be the entire
Hockeypuck source tree, but since it has been broken up into separate package
projects, it is used primarily for project documentation.

### Hockeypuck testing

https://github.com/hockeypuck/testing contains testdata used by some of the
other projects, such as OpenPGP key material examples used in
[openpgp](https://gopkg.in/hockeypuck/openpgp) test cases. It also contains
Ansible playbooks used to coordinate integration tests. These might be a useful
starting point, if you would like to automate your Hockeypuck deployment with
Ansible.

## Tools

### gopkg

https://gopkg.in is used to version Hockeypuck APIs. This is a distinct concern
from dependency revision management. Versioned APIs provide certain guarantees
with regard to compatibility.

Hockeypuck has vendored [logrus](https://gopkg.in/hockeypuck/logrus) to guard against upstream API changes.

### godeps

The [packaging](https://github.com/hockeypuck/packaging) project uses [godeps](https://launchpad.net/godeps) for dependency manangement. I've tried many dependency management solutions for Go packages, and I find godeps to be simple and easy to work with.

### Travis CI

Many of the Hockeypuck projects use Travis CI to build and test the projects on
commit. I've found Travis to be more useful than not for small, simple projects
with short tests. Which is really where you want to be. It's also nice to test
on many versions of Go, which is something I would not bother to do on every
commit otherwise.

## Pull request guidelines

Github pull requests will be reviewed on merit of relevance to the project
goals. Significant feature development should be discussed first on the
[mailing list](https://groups.google.com/forum/#!forum/hockeypuck-devel).

In code reviews, I'll look & ask for:

* Correctness. Does it do what it claims to?
* Succinct and appropriate naming.
* General Go style (Effective Go, etc.) and codebase fit. When in doubt, blend in.
* Godoc comments on public symbols (an area I'd like to improve)
* Adequate test coverage.

