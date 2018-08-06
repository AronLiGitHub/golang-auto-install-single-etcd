#!/bin/bash
basedir=$(cd `dirname $0`;pwd)

set -e

## function and implments
function set_ip(){
host_ip=`python -c "import socket;print([(s.connect(('8.8.8.8', 53)), s.getsockname()[0], s.close()) for s in [socket.socket(socket.AF_INET, socket.SOCK_DGRAM)]][0][1])"`
echo $host_ip
}

set_ip

