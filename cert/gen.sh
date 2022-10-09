#!/usr/bin/env bash

rm -rf *.pem *.srl *.cnf

openssl req -x509 -newkey rsa:4096 -days 36500 -nodes -keyout server-key.pem -out server-cert.pem -subj "/O=$SERVER_NAME/CN=$SERVER_IP" -addext "subjectAltName=IP:$SERVER_IP"
