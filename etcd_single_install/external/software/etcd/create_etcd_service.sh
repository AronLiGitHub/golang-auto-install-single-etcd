#!/bin/bash
basedir=$(cd `dirname $0`;pwd)
setip=$basedir/../set_ip.txt
etcdSSLDir=/etc/etcd/etcdSSL
etcdpem=$etcdSSLDir/etcd.pem
etcdkey=$etcdSSLDir/etcd-key.pem
capem=$etcdSSLDir/ca.pem

################## Set ETCD SERVER IP ######################

ETCD_SERVER1_NAME="`cat $setip | grep ETCD_SERVER1 | head -n 1  | cut -f 2 -d "="`"
ETCD_SERVER1_IP="`cat $setip | grep ETCD_SERVER1 | head -n 1  | cut -f 3 -d "="`"
ETCD_SERVER1_2380="$ETCD_SERVER1_NAME=https://$ETCD_SERVER1_IP:2380"
echo "ETCD_SERVER1_2380:$ETCD_SERVER1_2380"

ETCD_SERVER_2380="$ETCD_SERVER1_2380"
echo $ETCD_SERVER_2380

################# create etcd.service #####################
function createEtcdService(){
cat <<EOF > $basedir/etcd.service
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
WorkingDirectory=/var/lib/etcd/
EnvironmentFile=-/etc/etcd/etcd.conf
# set GOMAXPROCS to number of processors
ExecStart=/usr/bin/etcd `echo '\'`
  --name `echo '${ETCD_NAME} \'`
  --cert-file=$etcdpem `echo '\'`
  --key-file=$etcdkey `echo '\'`
  --peer-cert-file=$etcdpem `echo '\'`
  --peer-key-file=$etcdkey `echo '\'`
  --trusted-ca-file=$capem `echo '\'`
  --peer-trusted-ca-file=$capem `echo '\'`
  --initial-advertise-peer-urls `echo '${ETCD_INITIAL_ADVERTISE_PEER_URLS}'` `echo '\'`
  --listen-peer-urls `echo '${ETCD_LISTEN_PEER_URLS}'` `echo '\'`
  --listen-client-urls `echo '${ETCD_LISTEN_CLIENT_URLS}'`,http://127.0.0.1:2379 `echo '\'`
  --advertise-client-urls `echo '${ETCD_ADVERTISE_CLIENT_URLS}'` `echo '\'`
  --initial-cluster-token `echo '${ETCD_INITIAL_CLUSTER_TOKEN}'` `echo '\'`
  --initial-cluster $ETCD_SERVER_2380 `echo '\'`
  --initial-cluster-state new `echo '\'`
  --data-dir=`echo '${ETCD_DATA_DIR}'`

Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target

EOF
}

createEtcdService
