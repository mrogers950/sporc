#!/bin/bash

create_ca() {
  echo "==== Generating CA $1 ===="
  openssl genrsa -out CA$1.key 2048
  openssl req -new -sha256 -key CA$1.key -out CA$1.csr -subj "/CN=$1" -config ../openssl_ca.cnf
  openssl x509 -signkey CA$1.key -in CA$1.csr -req -days 3650 -out CA$1.crt -extensions v3_req -extfile ../openssl_ca.cnf
}

create_leaf_cert_signed_by_ca() {
  echo "==== Generating cert $2 signed by CA $1 ===="
  openssl genrsa -out Leaf$2_signed_by_CA$1.key 2048
  openssl req -new -sha256 -key Leaf$2_signed_by_CA$1.key -out Leaf$2_signed_by_CA$1.csr -subj "/CN=leaf-$2-by-$1" -config ../openssl_leaf.cnf
  openssl x509 -req -in Leaf$2_signed_by_CA$1.csr -CA CA$1.crt -CAkey CA$1.key -CAcreateserial -out Leaf$2_signed_by_CA$1.crt -days 3650 -sha256 -extensions v3_req -extfile ../openssl_leaf.cnf
  openssl x509 -text -noout -in Leaf$2_signed_by_CA$1.crt
}

create_crl() {
  echo "==== Generating CRL by $1 ===="
  openssl ca -gencrl -keyfile CA$1.key -cert CA$1.crt -out $1.crl.pem -config ../openssl_ca.cnf 
  openssl crl -inform PEM -in $1.crl.pem -outform PEM -out $1.crl
}

revoke_cert() {
  echo "==== for CA $1 revoking server cert $2 ===="
  openssl ca -config ../openssl_ca.cnf -revoke Leaf$2_signed_by_CA$1.crt -keyfile CA$1.key -cert CA$1.crt
}

rm -rf generated
mkdir -p generated
pushd generated
touch index.txt
echo 0 > crlnumber

create_ca "norev"
create_ca "tworev"
create_leaf_cert_signed_by_ca "norev" "leaf1"
create_leaf_cert_signed_by_ca "norev" "leaf2"
create_leaf_cert_signed_by_ca "tworev" "leaf1"
create_leaf_cert_signed_by_ca "tworev" "leaf2"
create_leaf_cert_signed_by_ca "tworev" "leaf3"
create_crl "norev"
create_crl "tworev"
revoke_cert "tworev" "leaf1"
revoke_cert "tworev" "leaf2"
# update CRL
create_crl "tworev"

rm -f *.csr
rm -f *.srl

popd
