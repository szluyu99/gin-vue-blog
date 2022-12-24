#!/bin/bash

mysql -uroot -p$MYSQL_ROOT_PASSWORD << EOF
system echo '================Start create database gvb====================';
source $WORK_PATH/gvb.sql
system echo '================OK!====================';
EOF