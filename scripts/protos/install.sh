#!/usr/bin/env bash

protoc_version=26.1

if [[ "$(echo "$(protoc --version)" | cut -d' ' -f2)" == "$protoc_version" ]]; then
  echo "protoc version $protoc_version is already installed"
  exit 0
fi

os_type=$(uname -s | tr '[:upper:]' '[:lower:]')
os_arch=$(uname -m)

if [[ "$os_type" == "darwin" ]]; then
  os_type="osx"
fi

if [[ "$os_arch" == "arm64" ]]; then
  os_arch="aarch_64"
fi

curl \
-Lo /var/tmp/protoc.zip \
"https://github.com/protocolbuffers/protobuf/releases/download/v$protoc_version/protoc-$protoc_version-$os_type-$os_arch.zip" \

unzip /var/tmp/protoc.zip -d /var/tmp/protoc && \
sudo mv /var/tmp/protoc/bin/protoc /usr/local/bin/protoc && \
sudo mv /var/tmp/protoc/include/* /usr/local/include/ && \
rm -rf /var/tmp/protoc /var/tmp/protoc.zip