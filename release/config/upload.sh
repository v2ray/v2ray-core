#!/bin/sh

journalctl -u v2ray --since today | curl -X PUT -s --upload-file "-" https://transfer.sh/v2ray.log | awk '{print "{\"value1\":\""$1"\"}"}' | curl -s --header "Content-Type: application/json" --request POST --data @- https://www.v2ray.com/logupload/
