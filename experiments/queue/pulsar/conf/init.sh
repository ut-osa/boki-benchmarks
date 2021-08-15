#!/bin/bash

ZK_HOST="zookeeper:2181"

bin/pulsar initialize-cluster-metadata \
  --cluster pulsar-cluster \
  --zookeeper ${ZK_HOST} --configuration-store ${ZK_HOST} \
  --web-service-url http://pulsar-broker:8080 \
  --broker-service-url pulsar://pulsar-broker:6650

if [[ $? != 0 ]]; then
    echo "Failed to setup pulsar"
    exit 1
fi

sleep 20

bin/pulsar-admin topics create-partitioned-topic \
  persistent://public/default/global-queue-0 \
  --partitions 64

sleep infinity
