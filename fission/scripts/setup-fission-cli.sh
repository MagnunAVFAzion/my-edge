#!/bin/bash

tmp_dir=$(mktemp -d)
cd "$tmp_dir" || exit

curl -Lo fission https://github.com/fission/fission/releases/download/v1.20.4/fission-v1.20.4-darwin-amd64
chmod +x fission
sudo mv fission /usr/local/bin/

cd - || exit
rm -rf "$tmp_dir"