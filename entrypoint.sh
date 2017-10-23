#!/bin/bash

mongod &
service tor restart
ls -alh /arpit
chmod +x /arpit/main
echo "Starting server..."
cd /arpit
varnishd -f /arpit/default.vcl -s malloc,100M -T 127.0.0.1:2000 -a 0.0.0.0:8081
/arpit/https_server &
/arpit/main
