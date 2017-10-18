#!/bin/bash

service tor restart
ls -alh /arpit
chmod +x /arpit/main
cat /arpit/private.key
echo "Starting server..."
bash /arpit/main
