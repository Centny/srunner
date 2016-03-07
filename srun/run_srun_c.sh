#!/bin/bash
cd `dirname ${0}`
export PATH=`pwd`:`dirname ${0}`:$PATH
srun conf/srun.properties