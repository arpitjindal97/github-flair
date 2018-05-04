#!/bin/bash

ls -alh /arpit
chmod +x /arpit/main
echo "Starting server..."
cd /arpit
/arpit/main &
