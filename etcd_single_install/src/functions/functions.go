package functions

import "../tools"

//定义全局变量
var (
	externalDir = "../../external/"
	etcdSourceDir  = externalDir + "software/etcd/"        //设置etcdTLS创建证书的目录
	cfsslDir    = externalDir + "software/cfssl/" // 设置cfssl的软件目录
	optDir  = "/opt/"    // 设置定义生成shell脚本的根目录
	optInstallDir = optDir + "install_etcd_single/" //定义shell脚本的执行目录
	basedir = "basedir=$(cd `dirname $0`;pwd)" //定义shell脚本的当前目录代码
	//定义shell调用python获取本地IP的语句
	hostip = "host_ip=`python -c \"import socket;print([(s.connect(('8.8.8.8', 53)), s.getsockname()[0], s.close()) for s in [socket.socket(socket.AF_INET, socket.SOCK_DGRAM)]][0][1])\"`"
)

// 执行所有步骤的函数
func Implements()  {
	init_dir()
	init_software()
	create_ca()
	create_setup_etcd()
	create_check_etcd()
	create_export_etcd_endpoint()
	create_implements_shell()
	start_install_etcd()
	end_deal()
}

// 初始化相关的文件目录
func init_dir(){
	// 创建目录 install_etcd_single
	sh1 := "mkdir -p "+ optDir +"install_etcd_single"
	tools.ExecShell(sh1)
	
	// 创建cfssl软件包路径
	sh2 := "mkdir -p "+ optDir +"install_etcd_single/Step1_file/cfssl"
	tools.ExecShell(sh2)
	
	// 创建etcdSSL证书路径
	sh3 := "mkdir -p "+ optDir +"install_etcd_single/Step1_file/etcdSSL"
	tools.ExecShell(sh3)
	
	// 创建etcd的rpm包文件夹路径
	sh4 := "mkdir -p "+ optDir +"install_etcd_single/Step2_file/etcd_rpm"
	tools.ExecShell(sh4)
}

// 拷贝文件夹中所需的软件至预备安装的路径
func init_software(){
	
	// 拷贝cfssl的二进制可执行文件至安装路径
	cfsslDir := externalDir + "software/cfssl/" // 设置cfssl的软件目录
	cfsslSoftDir := optDir + "install_etcd_single/Step1_file/cfssl/"                 //设置拷贝的目标目录
	sh1 := "cp " + cfsslDir + "cfssl* " + cfsslSoftDir
	sh2 := "cd " + cfsslSoftDir + " && ls cfssl*"
	sh3 := "cd " + cfsslSoftDir + " && chmod +x cfssl*"
	tools.ExecShell(sh1)
	tools.ExecShell(sh2)
	tools.ExecShell(sh3)
	
	// 拷贝ca证书生成文件
	sh4 := "cp " + cfsslDir + "ca* " + cfsslSoftDir
	sh5 := "cp " + cfsslDir + "create_etcd_csr.sh " + cfsslSoftDir
	sh6 := "cd " + cfsslSoftDir + " && chmod +x create_etcd_csr.sh"
	tools.ExecShell(sh4)
	tools.ExecShell(sh5)
	tools.ExecShell(sh6)
	
	// 拷贝etcd安装所需的文件
	etcdSoftDir := optInstallDir + "Step2_file"
	etcdSoftRpmDir := optInstallDir + "Step2_file/etcd_rpm"
	sh7 := "cp -v "+ etcdSourceDir +"create* " + etcdSoftDir
	sh8 := "cd " + etcdSoftDir + " && chmod +x create*"
	sh9 := "cp -v "+ etcdSourceDir +"etcd* " + etcdSoftRpmDir
	tools.ExecShell(sh7)
	tools.ExecShell(sh8)
	tools.ExecShell(sh9)
}

// 创建执行创建CA证书的步骤文件
func create_ca(){
	
	// 1、定义shell脚本内容
	shell := `
configdir=$basedir/Step1_file
cfssldir=$configdir/cfssl
ssldir=$configdir/etcdSSL
ipfile=$basedir/set_ip.txt

set -e

## function and implments
function set_ip(){
`+ hostip + `
cat <<EOF > $ipfile
################ auto Set config-base file ################
ETCD_SERVER1=infra1=$host_ip
EOF
}

set_ip

function check_firewalld_selinux(){
  setenforce 0 || true
  systemctl stop iptables.service || true
  systemctl stop firewalld.service || true
  systemctl status firewalld
  /usr/sbin/sestatus -v
}

#check_firewalld_selinux

function install_cfssl(){
  cp $cfssldir/cfssl_linux-amd64 /usr/local/bin/cfssl
  cp $cfssldir/cfssljson_linux-amd64 /usr/local/bin/cfssljson
  cp $cfssldir/cfssl-certinfo_linux-amd64 /usr/local/bin/cfssl-certinfo
  cd /usr/local/bin && ls cfssl*
}

install_cfssl

function create_etcd_csr(){
 sh $cfssldir/create_etcd_csr.sh
}

create_etcd_csr

function create_ssl(){
  cd $configdir && rm -rf ssl
  mkdir -p $ssldir
  cp $cfssldir/ca-config.json $ssldir/ca-config.json
  cp $cfssldir/ca-csr.json $ssldir/ca-csr.json
  cat $ssldir/ca-config.json
  cat $ssldir/ca-csr.json
  cd $ssldir && cfssl gencert -initca ca-csr.json | cfssljson -bare ca
  cd $ssldir && ls ca*
}

create_ssl

function create_etcd_certificates(){
  cp $cfssldir/etcd-csr.json  $ssldir/etcd-csr.json
  cat $ssldir/etcd-csr.json
  cd $ssldir && cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=etcd etcd-csr.json | cfssljson -bare etcd
  cd $ssldir && ls etcd*
}

create_etcd_certificates

function set_etcd_ssl(){
 mkdir -p /etc/etcd && rm -rf /etc/etcd/etcdSSL
 cp -r $ssldir /etc/etcd
}

set_etcd_ssl
	`
	generate_shell("Step1_create_CA.sh",shell)
}

// 创建执行安装etcd的步骤文件
func create_setup_etcd(){
	
	// 1、定义shell脚本内容
	shell := `
configdir=$basedir/Step2_file
etcdSSLDir=/etc/etcd/etcdSSL
etcdRpm=$configdir/etcd_rpm/etcd-3.2.22-1.el7.x86_64.rpm

## function and implments
function install_etcd(){
  #rpm -qa|grep etcd
  rpm -e etcd-3.2.22-1.el7.x86_64
  service etcd stop
  yum erase etcd -y
  cd /var/lib/ && rm -rf etcd
  rpm -ivh $etcdRpm
}

install_etcd

function set_service(){
sh $configdir/create_etcd_conf.sh
sh $configdir/create_etcd_service.sh
cd /usr/lib/systemd/system && mv etcd.service etcd.service.bak
cd $configdir && cp etcd.service /usr/lib/systemd/system
cd /etc/etcd/ && mv etcd.conf etcd.conf.bak
cd $configdir && cp etcd.conf /etc/etcd
}

set_service

function start_etcd(){
  systemctl daemon-reload
  systemctl start etcd
  systemctl status etcd
}

start_etcd

function check_health(){
etcdctl \
  --ca-file=$etcdSSLDir/ca.pem \
  --cert-file=$etcdSSLDir/etcd.pem \
  --key-file=$etcdSSLDir/etcd-key.pem \
  cluster-health
}

check_health
`
	generate_shell("Step2_setup_etcd1.sh",shell)
}

// 创建检查etcd健康的文件
func create_check_etcd()  {
	// 1、定义shell脚本内容
	shell := `
etcdSSLDir=/etc/etcd/etcdSSL

function check_health(){
etcdctl \
  --ca-file=$etcdSSLDir/ca.pem \
  --cert-file=$etcdSSLDir/etcd.pem \
  --key-file=$etcdSSLDir/etcd-key.pem \
  cluster-health
}

check_health
`
	generate_shell("Step3_check_etcd1.sh",shell)
}

// 创建导出etcd访问信息的文件
func create_export_etcd_endpoint(){
	
	// 0、单独处理带 ` 的字符串
	ETCD_SERVER1_NAME := "`cat $setip | grep ETCD_SERVER1 | head -n 1  | cut -f 2 -d \"=\"`"
	ETCD_SERVER1_IP := "`cat $setip | grep ETCD_SERVER1 | head -n 1  | cut -f 3 -d \"=\"`"
	
	// 1、定义shell脚本内容
	shell := `
setip=$basedir/set_ip.txt
etcdSSLDir=/etc/etcd/etcdSSL

###################### set ip ###############################

ETCD_SERVER1_NAME=` + ETCD_SERVER1_NAME + `
ETCD_SERVER1_IP=` + ETCD_SERVER1_IP + `
ETCD_SERVER1_2379="$ETCD_SERVER1_IP:2379"

## function and implments
function export_endpoints(){
cat <<EOF > /opt/ETCD_CLUSER_INFO
ETCD_ENDPOINT_2379=https://$ETCD_SERVER1_2379
CA_FILE=$etcdSSLDir/ca.pem
CERT_FILE=$etcdSSLDir/etcd.pem
KEY_FILE=$etcdSSLDir/etcd-key.pem
EOF
}

export_endpoints
`
	generate_shell("Step4_export_etcd_endpoint.sh",shell)
}

// 创建shell脚本的总执行文件
func create_implements_shell(){
	// 1、定义shell脚本内容
	shell := `
./Step1_create_CA.sh
./Step2_setup_etcd1.sh
./Step3_check_etcd1.sh
./Step4_export_etcd_endpoint.sh
`

	generate_shell("Implement.sh",shell)
}

// 开始安装etcd服务
func start_install_etcd(){
	sh := "cd "+ optInstallDir + " && ./Implement.sh"
	tools.ExecShell(sh)
}

// 安装后的处理
func end_deal(){
	//备份一个检查etcd服务健康的脚本
	sh1 := "cd "+ optInstallDir + " && cp Step3_check_etcd1.sh " + optDir
	tools.ExecShell(sh1)
	
	//删除安装的执行脚本
	sh2 := "cd " + optDir + " && rm -rf install_etcd_single"
	tools.ExecShell(sh2)
}

// 抽象生成shell脚本的方法
func generate_shell(filename string,shell string){
	// 1、生成shell脚本
	//filename := "./" + filename
	content := `#!/bin/bash
`+ basedir + shell
	
	tools.WriteWithIoutil(filename,content)
	
	// 2、拷贝shell脚本至指定目录
	sh1 := "cp -v " + filename + " " + optInstallDir
	sh2 := "cd " + optInstallDir + " && chmod +x " + filename
	sh3 := "rm -rf " + filename
	tools.ExecShell(sh1)
	tools.ExecShell(sh2)
	tools.ExecShell(sh3)
}