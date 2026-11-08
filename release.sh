#!/bin/sh
set -xe

OS="${1:-}"
VERSION="${2:-}"
FLAG="${3:-}"

if test -z "$OS" || test -z "$VERSION"; then
  echo "Usage: $0 <windows|linux|darwin> <version> [ -arm ]"
  exit 1
fi

GOARCH="amd64"
if test "$FLAG" = "-arm"; then
  GOARCH="arm64"
fi

EXT=""
if test "$OS" = "windows"; then
  EXT=".exe"
fi

export GOOS="$OS"
export GOARCH="$GOARCH"
export CGO_ENABLED="0"

OUT="bench-cmp-${GOOS}-${GOARCH}-${VERSION}${EXT}"
echo "Building $OUT (GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO_ENABLED)"

go build -o "$OUT" ./cmd
echo "Created $OUT"
