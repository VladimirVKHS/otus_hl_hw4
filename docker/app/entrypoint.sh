#!/bin/bash

# Start the first process
/usr/sbin/zabbix_agentd -c /etc/zabbix/zabbix_agentd.conf &

# Start the second process
/usr/bin/main