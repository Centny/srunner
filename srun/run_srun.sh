#!/bin/bash
cd `dirname ${0}`
export PATH=`pwd`:`dirname ${0}`:$PATH
conf=conf/srun.properties
if [ "$SR_C" != "" ];then
	conf=$SR_C
fi
for (( i = 1; i>0; i++ )); do
	if [ -f ".sr_done" ]; then
		break
	fi
	srun $conf
done
rm -f .sr_done