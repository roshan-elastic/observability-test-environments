#!/bin/env bash
# rootCA
openssl genrsa -out rootCA.key 2048
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3600 -out rootCA.pem -subj  '/O=Elastic/OU=Obs/CN=oblt-cli-CA/emailAddress=observability-robots@elastic.co'
openssl x509 -in rootCA.pem -text -noout

# Server certificate
openssl genrsa -out tls.key 2048

openssl req -new -key tls.key -out tls.csr -subj '/O=Elastic/OU=Obs/CN=vpn/emailAddress=observability-robots@elastic.co'
openssl req -in rootCA.pem -text -noout

cat > openssl.cnf <<EOF
# Extensions to add to a certificate request
basicConstraints       = CA:FALSE
authorityKeyIdentifier = keyid:always, issuer:always
keyUsage               = nonRepudiation, digitalSignature, keyEncipherment, dataEncipherment
subjectAltName         = @alt_names
[ alt_names ]
DNS.1 = localhost.aviatrix.com
EOF

openssl x509 -req \
    -in tls.csr \
    -CA rootCA.pem \
    -CAkey rootCA.key \
    -CAcreateserial \
    -out tls.crt \
    -days 3650 \
    -sha256 \
    -extfile openssl.cnf

sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain rootCA.pem
