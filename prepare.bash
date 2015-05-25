#!/bin/bash -ex

export DEBEMAIL="cmars@cmarstech.com"
export DEBFULLNAME="Casey Marshall"

export RELEASE_VERSION=1.0.1~rc1
export BUILD_PACKAGE=github.com/hockeypuck/hockeypuck
export BUILD_BRANCH=1.0.1-dump

### Set up GOPATH
export GOPATH=$(pwd)

for pkg in launchpad.net/godeps github.com/mitchellh/gox; do
	go get ${pkg}
	go install ${pkg}
done

mkdir -p ${GOPATH}/src/github.com/hockeypuck
cd ${GOPATH}/src/github.com/hockeypuck
if [ ! -d "hockeypuck" ]; then
	git clone git@github.com:hockeypuck/hockeypuck.git
else
	(cd hockeypuck; git fetch origin)
fi
(cd hockeypuck; git checkout ${BUILD_BRANCH})

cd ${GOPATH}/src/${BUILD_PACKAGE}
${GOPATH}/bin/godeps -u dependencies.tsv

export SHORTHASH=$(git log -1 --pretty=format:%h)
export LONGHASH=$(git log -1 --pretty=format:%H)
export HEXDATE=$(date +%s)

# Get our current and last built revision
export LTS_SERIES="precise trusty"
export PACKAGE_VERSION="${RELEASE_VERSION}~${HEXDATE}+${SHORTHASH}"

cd ${GOPATH}
echo "$LONGHASH" > version-git-commit
echo "$PACKAGE_VERSION" > version-release

