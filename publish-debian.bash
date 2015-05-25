#!/bin/bash -xe

REPO=unstable
LTS_SERIES="precise trusty"
PACKAGE_VERSION=$(cat version-release)

for SERIES in $LTS_SERIES; do
	CHANGES=../hockeypuck1-migration-tools_${PACKAGE_VERSION}~${SERIES}_source.changes
	dput ppa:hockeypuck/${REPO} $CHANGES
done
