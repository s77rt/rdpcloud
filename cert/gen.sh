#!/usr/bin/env bash

if [ "$IS_FREE_TRIAL" == "FALSE" ]; then
	DAYS=36500
elif [ "$IS_FREE_TRIAL" == "TRUE" ]; then
	DAYS=$FREE_TRIAL_DURATION
else
	echo IS_FREE_TRIAL can either be TRUE or FALSE
	exit 1
fi

rm -f *.pem

openssl req -x509 -newkey rsa:4096 -days $DAYS -nodes -keyout server-key.pem -out server-cert.pem -subj "/O=$SERVER_NAME/CN=$SERVER_IP" -addext "subjectAltName=IP:$SERVER_IP"
