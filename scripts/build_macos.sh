#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-$(git describe --tags --always)}"
DIST_DIR="dist"

mkdir -p "$DIST_DIR"

build_target() {
  local arch="$1"
  local out="${DIST_DIR}/brewls_${VERSION}_darwin_${arch}"
  local tarball="${out}.tar.gz"

  echo "Building darwin/${arch}..."
  GOOS=darwin GOARCH="$arch" go build -o "$out" ./cmd/brewls
  tar -czf "$tarball" -C "$DIST_DIR" "brewls_${VERSION}_darwin_${arch}"
  shasum -a 256 "$tarball" > "${tarball}.sha256"
  rm -f "$out"
}

build_target amd64
build_target arm64

echo "Artifacts written to ${DIST_DIR}/"
