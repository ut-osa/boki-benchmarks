#!/bin/bash

export ZOO_LOG4J_PROP="WARN,CONSOLE"

SERVER="${ZK_HOST:-zookeeper:2181}"

./bin/zkCli.sh -server $SERVER <<EOF
deleteall /faas
create /faas
create /faas/node
create /faas/view
create /faas/freeze
create /faas/cmd
quit
EOF
if [[ $? != 0 ]]; then
    echo "Failed to setup zookeeper"
    exit 1
fi

sleep 50

./bin/zkCli.sh -server $SERVER create /faas/cmd/start

sleep infinity
