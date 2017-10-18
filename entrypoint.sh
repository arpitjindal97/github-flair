#!/bin/bash

service tor restart
ls -alh /arpit
chmod +x /arpit/main
echo "Starting server..."
/arpit/main
