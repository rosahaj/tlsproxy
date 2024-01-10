#!/bin/bash
set -ex
# generate CA's  key
openssl genrsa -aes256 -passout pass:1 -out key.pem 4096
openssl rsa -passin pass:1 -in key.pem -out key.pem

openssl req -config certs/openssl.cnf -key key.pem -new -x509 -days 7300 -sha256 -extensions v3_ca -out cert.pem
