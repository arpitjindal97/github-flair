#!/bin/bash

ls -alh /arpit
chmod +x /arpit/main
echo "Starting mongod..."
mongod &
echo "Starting server..."
cd /arpit
/arpit/main
