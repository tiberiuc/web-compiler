#!/bin/bash

SCRIPT=`realpath $0`
SCRIPTPATH=`dirname $SCRIPT`

echo $SCRIPTPATH

export PATH=$SCRIPTPATH/vendors/bin:$PATH

export GEM_HOME="$SCRIPTPATH/vendors/"
export GEM_PATH="${GEM_HOME}"

