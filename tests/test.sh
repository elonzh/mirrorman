#!/usr/bin/env bash
set -eux
export HTTP_PROXY=http://127.0.0.1:8081
export HTTPS_PROXY=http://127.0.0.1:8081
export NO_PROXY=localhost,example.com
echo $HTTP_PROXY $HTTPS_PROXY
cache_dir="tests/.cache"
downloads_dir="tests/downloads"
rm -rf ${cache_dir} ${downloads_dir}
mkdir -p ${downloads_dir}
cd ${downloads_dir}
# TODO: using unit tests
urls=(
#  "https://github.com/gohugoio/hugo/releases/download/v0.73.0/hugo_0.73.0_Linux-64bit.deb"
#  "https://github.com/gohugoio/hugo/archive/v0.73.0.tar.gz"
#  "https://mirrors.huaweicloud.com/helm/v3.2.4/helm-v3.2.4-linux-amd64.tar.gz"
# "https://github.com/rancher/k3s/releases/download/v1.18.6-rc1%2Bk3s1/k3s"
# "https://github.com/rancher/k3s/releases/download/v1.18.6-rc1+k3s1/k3s"
  "http://127.0.0.1:8080/proxy/https://mirrors.huaweicloud.com/helm/v3.2.4/helm-v3.2.4-linux-amd64.tar.gz"
#  "http://127.0.0.1:8080/proxy/https://github.com/rancher/k3s/releases/download/v1.18.6-rc1+k3s1/k3s"
#  "https://mirrors.aliyun.com/ubuntu-releases/19.10/ubuntu-19.10-live-server-amd64.iso"
#  "https://storage.googleapis.com/golang/go1.13.4.linux-amd64.tar.gz"
)
for i in "${!urls[@]}"; do
  url=${urls[$i]}
  curl --config ../.curlrc --dump-header "${i}_remote_header.txt" --url "${url}"
  curl --config ../.curlrc --dump-header "${i}_proxy_header.txt" --url "${url}"
done
