#!/bin/bash
ROOT_PATH=$(cd $(dirname $0) && pwd)

echo "Installing to $ROOT_PATH/local/"

mkdir -p "$ROOT_PATH/local"

curl -s https://en-gb.wordpress.org/latest-en_GB.tar.gz -o "$ROOT_PATH/local/wordpress.tgz"
tar xvzf "$ROOT_PATH/local/wordpress.tgz" -C "$ROOT_PATH/local/"

cp "$ROOT_PATH/wp-config.php" "$ROOT_PATH/local/wordpress/"
