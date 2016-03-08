#!/bin/bash
cd `dirname ${0}`
export PATH=`pwd`:`dirname ${0}`:$PATH
touch .sr_done
srun -c http://127.0.0.1:3010/_exit_