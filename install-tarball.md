# Installing a binary tarball release

# Download
Download a [gzip-compressed tar archive release from
Github](https://github.com/hockeypuck/hockeypuck/releases) for your operating
system and architecture. Generally, Hockeypuck can be built for any Unix-like
platform that the Go language compiler and linker supports.

# Install

## Extract into '/'
The archive can be extracted into '/'. This will preserve the path references
in the archived files.

## Or chroot, run in a container, etc.
For added security, you could extract into an arbitrary path and chroot the
Hockeypuck process, or extract it over an LXC rootfs. If you do this, consider
the implications for a local UNIX domain socket connection to PostgreSQL.

# Packaging
The Hockeypuck binary archive distributions could be a useful starting point to
build packages for other operating system distributions. Contributions to the [packaging](https://github.com/hockeypuck/packaging) project would be welcome for RPMs, BSD ports, etc.

# Next steps

* [Configure](configuration.md) the Hockeypuck server.
* Run Hockeypuck with `/path/to/hockeypuck -config /path/to/hockeypuck.conf`.

