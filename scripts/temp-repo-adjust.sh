#!/usr/bin/env bash

echo "==> Temporarily adjusting repo directories ..."

export ORIG_TRAVIS_BUILD_DIR=${TRAVIS_BUILD_DIR}
mkdir ${TRAVIS_HOME}/gopath/src/github.com/skytap
cd ${TRAVIS_HOME}/gopath/src/github.com/skytap
mv ${ORIG_TRAVIS_BUILD_DIR} ${TRAVIS_HOME}/gopath/src/github.com/snidell/terraform-provider-techzoneskytap/skytap
export TRAVIS_BUILD_DIR=${TRAVIS_HOME}/gopath/src/github.com/snidell/terraform-provider-techzoneskytap/skytap
cd ${TRAVIS_HOME}/gopath/src/github.com/snidell/terraform-provider-techzoneskytap/skytap

exit 0
