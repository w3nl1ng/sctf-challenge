#!/bin/bash
set -e

service mysql start

mysql < /mysql/create_db.sql
sleep 3

service mysql restart

tail -f /dev/null