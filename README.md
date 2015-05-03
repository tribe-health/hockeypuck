# Hockeypuck

Hockeypuck is an OpenPGP [public keyserver](https://en.wikipedia.org/wiki/Key_server_(cryptographic)).

# Latest Release

The latest release of Hockeypuck is 2.0-b2, which is a significant redesign of
the original Hockeypuck 1.x server.

# Features

## OpenPGP Standards

Hockeypuck implements the
[HKP draft protocol specification](http://ietfreport.isoc.org/idref/draft-shaw-openpgp-hkp/)
as well as several extensions to the protocol supported by [SKS](http://sks-keyservers.net/).

Public key material conforming to [RFC
4880](https://tools.ietf.org/html/rfc4880) is supported by the keyserver, as
are [RFC 6637](https://tools.ietf.org/html/rfc6637>) ECC keys. As-of-yet unsupported key
material, such as recent Ed25519 signing keys, may be distributed by
Hockeypuck, however Hockeypuck is not able to validate them yet.

## SKS Reconciliation Protocol
Hockeypuck is the first keyserver besides SKS that is capable of synchronizing public
key material with peers using the SKS reconciliation protocol. Recon protocol
support is provided with the [Conflux](https://gopkg.in/hockeypuck/conflux.v2) package.

## Modular storage backend

Hockeypuck may use either MongoDB or PostgreSQL 9.4 for storing indexed key
material. The architecture supports additional storage backends, which must
implement a simple Go interface.

## Flexible rendering

Hockeypuck internally represents key material with an arbitrary document model
that can be used with web applications by rendering it to JSON in responses.
Classic standard HKP responses are also supported, as are customizable HTML
template responses that operate on the document model.

# Install

Several options are available:
* [Build and Install from source](install-source.md)
* [Install from binaries](install-tarball.md)
* [Install from Ubuntu archives](install-ubuntu.md)
* [Deploy with Juju](juju.md)

# Next Steps

* [Configuring](configuration.md) a Hockeypuck server.
* [Peering](peering.md) with other compatible public keyservers.
* [Contributing](contributing.md) to Hockeypuck.
* [Community](community.md) support for Hockeypuck, for issues and new feature requests.

