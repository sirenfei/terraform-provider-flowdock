#!/usr/bin/env bash

VERSION=$(cat version)
echo "building terraform-provider-flowdock_${VERSION}"

go build -o terraform-provider-flowdock_${VERSION}