#!/bin/bash
set -ex
# generate CA's  key
pass=9876
openssl genrsa -aes256 -passout pass:${pass} -out ca.key.pem 4096
openssl rsa -passin pass:${pass} -in ca.key.pem -out ca.key.pem.tmp
mv ca.key.pem.tmp ca.key.pem

openssl req -batch -config openssl.cnf -key ca.key.pem -new -x509 -days 7300 -sha256 -extensions v3_ca -out ca.pem
