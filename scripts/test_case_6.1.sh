#!/bin/bash
#
# test case 6.1
#
# xfer funds between lite accounts
# ids, amount and server IP:Port needed
#
# check for command line parameters
#

# if IDs not entered on the command line, prompt for one and exit

if [ -z $1 ]; then
        echo "Usage: test_case_6.1.sh fromID toID numTokens IPAddress:Port"
        exit 0
fi

# see if $1 is really an ID

id1=$1
size=${#id1}
if [ $size -lt 59 ]; then
        echo "Expected acc://<48 byte string>/ACME"
        exit 0
fi

if [ -z $2 ]; then
        echo "Usage: test_case_6.1.sh fromID toID numTokens IPAddress:Port"
        exit 0
fi

# see if $2 is really an ID

id2=$2
size=${#id2}
if [ $size -lt 59 ]; then
        echo "Expected acc://<48 byte string>/ACME"
        exit 0
fi

if [ -z $3 ]; then
        echo "Usage: test_case_6.1.sh fromID toID numTokens IPAddress:Port"
        exit 0
fi

# see if the IP address and port were entered on the command line

if [ -z $4 ]; then
        echo "You must enter an IPAddress:Port for a server to transfer tokens"
        echo "Usage: test_case_6.1.sh fromID toID numTokens IPAddress:Port"
        exit 0
fi

# call our xfer script

./cli_xfer_tokens.sh $id1 $id2 $3 $4

