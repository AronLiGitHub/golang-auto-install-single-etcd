#!/bin/bash
basedir=$(cd `dirname $0`;pwd)

## function
function create_etcd_csr(){

host_ip=`python -c "import socket;print([(s.connect(('8.8.8.8', 53)), s.getsockname()[0], s.close()) for s in [socket.socket(socket.AF_INET, socket.SOCK_DGRAM)]][0][1])"`

cat <<EOF > $basedir/etcd-csr.json
{
  "CN": "etcd",
  "hosts": [
    "127.0.0.1",
    "$host_ip"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "shenzhen",
      "L": "shenzhen",
      "O": "etcd",
      "OU": "System"
    }
  ]
}
EOF
}

create_etcd_csr
