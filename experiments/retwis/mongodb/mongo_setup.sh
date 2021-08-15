#!/bin/bash

mongo mongodb://mongodb-1:27017 init_mongo_rs.js

if [[ $? != 0 ]]; then
    echo "Failed to setup mongodb"
    exit 1
fi

mongo mongodb://mongodb-1:27017 --eval "rs.conf()"

sleep infinity
