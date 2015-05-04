#!/bin/bash -e

if [ ! -d "publish" ]; then
	echo "run 'represent' to publish documentation."
	exit 1
fi

cd publish
rm -rf .git
git init
#git checkout -b gh-pages
#git remote add origin git@github.com:cmars/hockeypuck.git
git remote add origin git@github.com:hockeypuck/hockeypuck.github.io.git
git add *
git commit -m "Generated Hockeypuck documentation at $(date '+%y%m%d-%H%M%S')"
git checkout -b gh-pages
git push -f origin gh-pages
git push -f origin master
