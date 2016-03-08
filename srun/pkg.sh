#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export cpwd=`pwd`
export LD_LIBRARY_PATH=/usr/local/lib:/usr/lib
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
o_dir=build
if [ "$2" != "" ];then
	o_dir=$2/srun
fi
rm -rf $o_dir
mkdir -p $o_dir

#### Package ####
sys_n=`uname`
n_srv=srun
v_srv=0.0.1
n_srv_=$n_srv

if [ ${sys_n:0:7} = "MSYS_NT" ];then
	n_srv_=$n_srv".exe"
fi
##
d_srv="$n_srv"d
o_srv=$o_dir/$n_srv
mkdir $o_srv
mkdir $o_srv/conf
mkdir $o_srv/www
go build -o $o_srv/$n_srv_ github.com/Centny/srunner/srun
cp "$n_srv".properties $o_srv/conf
cp *.bat $o_srv
cp *.sh $o_srv

###
if [ "$1" != "" ];then
	curl -o $o_srv/srvd $1/srvd
	curl -o $o_srv/srvd_i $1/srvd_i
	chmod +x $o_srv/srvd
	chmod +x $o_srv/srvd_i
	echo "./srvd_i \$1 $n_srv \$2 \$3" >$o_srv/install.sh
	chmod +x $o_srv/install.sh
fi 
cd $o_dir
zip -r -q $n_srv.zip $n_srv
cd ../
echo "Package $n_srv..."
