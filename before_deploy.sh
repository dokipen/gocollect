#!/bin/bash
set -e

VERSION="0.0.1"
NAME="gocollect"

if [[ -n "${TRAVIS_BUILD_NUMBER}" ]]; then
    VERSION="${VERSION}-r${TRAVIS_BUILD_NUMBER}"
fi

if [[ "${TRAVIS_BRANCH}" != "production" ]]; then
    NAME="${TRAVIS_BRANCH}-${NAME}"
fi

mkdir debroot -p
cp debian/* debroot -r
mkdir debroot/usr/bin -p
cp "${GOPATH}/bin/gocollect" debroot/usr/bin
gem install fpm

fpm \
  -s dir \
  -t deb \
  --force \
  -C "debroot" \
  -a "amd64" \
  -n "${NAME}" \
  --version "${VERSION}" \
  .
