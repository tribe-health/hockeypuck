#!/bin/bash -ex

. prepare.bash

### Build source package for each supported series.

cd ${GOPATH}

# Build for each supported Ubuntu version
for SERIES in $LTS_SERIES; do
	cat >debian/changelog <<EOF
hockeypuck1-migration-tools (${PACKAGE_VERSION}~${SERIES}) ${SERIES}; urgency=medium

  * Migration tools for dumping keys on Hockeypuck 1.x server, for upgrading to
    2.x.

 -- $DEBFULLNAME <$DEBEMAIL>  $(date -u -R)
EOF

	dpkg-buildpackage -rfakeroot -d -S -k0x879CF8AA8DDA301A
done

