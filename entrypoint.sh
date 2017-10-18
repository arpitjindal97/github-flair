#!/bin/bash

service tor restart
ls -alh /arpit
chmod +x /arpit/main
echo "Starting server..."
cd /arpit
mongod &
/arpit/main
