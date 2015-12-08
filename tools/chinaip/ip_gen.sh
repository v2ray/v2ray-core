#!/bin/bash

APNIC_FILE="http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"
curl "${APNIC_FILE}" | grep ipv4 | grep CN | awk -F\| '{ printf("%s/%d\n", $4, 32-log($5)/log(2)) }' > ipv4.txt
