#!/usr/bin/env bash
set -eux
export HTTP_PROXY=http://127.0.0.1:9876
export HTTPS_PROXY=http://127.0.0.1:9876
export NO_PROXY=localhost,example.com
echo $HTTP_PROXY $HTTPS_PROXY
cache_dir="tests/.cache"
downloads_dir="tests/downloads"
rm -rf ${cache_dir} ${downloads_dir}
mkdir -p ${downloads_dir}
cd ${downloads_dir}
urls=(
#  "https://github.com/gohugoio/hugo/releases/download/v0.73.0/hugo_0.73.0_Linux-64bit.deb"
#  "https://github.com/gohugoio/hugo/archive/v0.73.0.tar.gz"
#  "https://mirrors.huaweicloud.com/helm/v3.2.4/helm-v3.2.4-linux-amd64.tar.gz"
  "https://storage.googleapis.com/golang/go1.13.4.linux-amd64.tar.gz"
)
for i in "${!urls[@]}"; do
  url=${urls[$i]}
  curl --config ../.curlrc --dump-header "${i}_remote_header.txt" --url "${url}"
  curl --config ../.curlrc --dump-header "${i}_proxy_header.txt" --url "${url}"
done
