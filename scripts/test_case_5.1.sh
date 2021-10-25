#!/bin/bash
#
# test case 5.1
#
# fund lite account
# id and server IP:Port needed
#
# check for command line parameters
#
if [ -z $1 ]; then
	echo "You must pass an ID to be funded"
	exit 0
fi
if [ -z $2 ]; then
	echo "You must pass IP:Port for the server to use"
	exit 0
fi

# call our faucet script

TxID=`./cli_faucet.sh $1 $2`

# remove the "s

TxID=`echo $TxID | /usr/bin/sed 's/"//g'`

# check transaction status

Status=`./cli_get_tx_status.sh $TxID $2`

if [ ! -z $Status ]; then
	echo "Invalid status received from get tx"
fi

# get our balance

bal=`./cli_get_balance.sh $1 $2`

echo $bal
