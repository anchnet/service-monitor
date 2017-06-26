#!/bin/bash
echo "$(pwd)/lib" > /etc/ld.so.conf.d/oracle-client.conf
ldconfig
echo 127.0.0.1 ${HOSTNAME} >> /etc/hosts

