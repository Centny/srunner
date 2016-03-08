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
	srun $conf 1>logs/srun_out.log 2>logs/srun_err.log
done
rm -f .sr_done