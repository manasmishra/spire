#!/bin/bash

set -e

COVERPROFILE=
if [ -n "${COVERALLS_TOKEN}" ]; then
    COVERPROFILE=profile.cov
    go install github.com/mattn/goveralls@v0.0.7
fi

COVERPROFILE="${COVERPROFILE}" make test

if [ -n "${COVERALLS_TOKEN}" ]; then
    "$(go env GOPATH)"/bin/goveralls -coverprofile="${COVERPROFILE}" \
            -service=github
fi
