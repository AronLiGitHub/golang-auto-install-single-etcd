#!/bin/bash
basedir=$(cd `dirname $0`;pwd)
setip=$basedir/../set_ip.txt

################## Set ETCD SERVER IP ######################
function createEtcdConf(){
ETCD_SERVER1_NAME="`cat $setip | grep ETCD_SERVER1 | head -n 1  | cut -f 2 -d "="`"
ETCD_SERVER1_IP="`cat $setip | grep ETCD_SERVER1 | head -n 1  | cut -f 3 -d "="`"
ETCD_SERVER1_2380="https://$ETCD_SERVER1_IP:2380"
echo "ETCD_SERVER1_2380:$ETCD_SERVER1_2380"
ETCD_SERVER1_2379="https://$ETCD_SERVER1_IP:2379"
echo "ETCD_SERVER1_2380:$ETCD_SERVER1_2379"

cat <<EOF > $basedir/etcd.conf
#[member]
ETCD_NAME=$ETCD_SERVER1_NAME
ETCD_DATA_DIR="/var/lib/etcd"
ETCD_LISTEN_PEER_URLS="$ETCD_SERVER1_2380"
ETCD_LISTEN_CLIENT_URLS="$ETCD_SERVER1_2379"

#[cluster]
ETCD_INITIAL_ADVERTISE_PEER_URLS="$ETCD_SERVER1_2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_ADVERTISE_CLIENT_URLS="$ETCD_SERVER1_2379"

EOF
}

createEtcdConf
