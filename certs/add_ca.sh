# NOTE: only tested in ubuntu
set -ex

#goproxy_version="$(go list -m all | grep -E "^github.com/elazarl/goproxy\s" | awk '{print $2}')"
#mod_path="$(go env GOPATH)/pkg/mod/github.com/elazarl/goproxy@${goproxy_version}"
#echo "mod path:" "${mod_path}"
ca_path="/usr/share/ca-certificates/extra"
sudo mkdir -vp "${ca_path}"
sudo openssl x509 -in "ca.pem" -inform PEM -out "${ca_path}/mirrorman.crt"
sudo dpkg-reconfigure ca-certificates
sudo update-ca-certificates
