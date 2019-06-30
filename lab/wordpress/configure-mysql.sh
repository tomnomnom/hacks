#!/bin/bash
su mysql -s /bin/sh -c '/usr/bin/mysqld --datadir=/var/lib/mysql' &
sleep 3
mysql -uroot -e "create database wordpress;"
mysql -uroot wordpress < /app/wordpress.sql
