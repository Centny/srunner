#!/bin/bash
cd `dirname ${0}`
export PATH=`pwd`:`dirname ${0}`:$PATH
conf=conf/srun.properties
if [ "$SR_C" != "" ];then
	conf=$SR_C
fi
srun $conf
