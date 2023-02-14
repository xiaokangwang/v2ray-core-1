#!/bin/bash

set -x -e

download() {
  curl -L "https://github.com/v2fly/v2ray-core/releases/download/$1/$2" >"$2"
}

downloadAndUnzip() {
  download "$1" "$2"
  unzip -n -d "${2%\.zip}" "$2"
}
mkdir -p assets

pushd assets
downloadAndUnzip "$1" "v2ray-linux-32.zip"
downloadAndUnzip "$1" "v2ray-linux-64.zip"
downloadAndUnzip "$1" "v2ray-linux-arm32-v6.zip"
downloadAndUnzip "$1" "v2ray-linux-arm32-v7a.zip"
downloadAndUnzip "$1" "v2ray-linux-arm64-v8a.zip"
downloadAndUnzip "$1" "v2ray-extra.zip"
popd

placeFile() {
  mkdir -p "context/$2"
  cp -R "assets/$1/$3" "context/$2/$3"
}

function generateStandardVersion() {
  placeFile "$1" "$2/bin" "v2ray"
}

function generateExtraVersion() {
  generateStandardVersion "$1" "$2"
  placeFile "$1" "$2/share" "geosite.dat"
  placeFile "$1" "$2/share" "geoip.dat"
  placeFile "$1" "$2/etc" "config.json"
  placeFile "v2ray-extra" "$2/share" "browserforwarder"
}

generateStandardVersion "v2ray-linux-32" "linux/386/std"
generateStandardVersion "v2ray-linux-64" "linux/amd64/std"
generateStandardVersion "v2ray-linux-arm32-v6" "linux/arm/v6/std"
generateStandardVersion "v2ray-linux-arm32-v7a" "linux/arm/v7/std"
generateStandardVersion "v2ray-linux-arm64-v8a" "linux/arm64/std"
generateStandardVersion "v2ray-linux-arm64-v8a" "linux/arm64/v8/std"

generateExtraVersion "v2ray-linux-32" "linux/386/extra"
generateExtraVersion "v2ray-linux-64" "linux/amd64/extra"
generateExtraVersion "v2ray-linux-arm32-v6" "linux/arm/v6/extra"
generateExtraVersion "v2ray-linux-arm32-v7a" "linux/arm/v7/extra"
generateExtraVersion "v2ray-linux-arm64-v8a" "linux/arm64/extra"
generateExtraVersion "v2ray-linux-arm64-v8a" "linux/arm64/v8/extra"