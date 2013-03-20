#!/bin/bash -ex

if [ $USER != "root" ] ; then
	        echo "need to run as root"
	fi

CA_CERT=`puppet agent --genconfig |grep -v "#" |grep localcacert |awk '{print $NF}'`  
HOST_CERT=`puppet agent --genconfig |grep -v "#" |grep hostcert |awk '{print $NF}'`
PRI_KEY=`puppet agent --genconfig |grep -v "#" |grep hostprivkey |awk '{print $NF}'`

echo | openssl s_client -connect puppet:8140 -cert $HOST_CERT -key $PRI_KEY -CAfile $CA_CERT 2>&1 |grep Verify |tail -1

