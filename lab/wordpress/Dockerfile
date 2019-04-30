FROM alpine

RUN apk --update add \
    bash \
    supervisor \
    php7-apache2 \
    mysql \
    mysql-client


# For whatever reason the apk package doesn't create
# this dir so apache fails to start
RUN mkdir /run/apache2

# Guess what this one's for...
RUN mkdir /run/mysqld
RUN chown mysql: /run/mysqld
RUN su mysql -s /bin/sh -c 'mysql_install_db --datadir=/var/lib/mysql'

# Apache config
ADD apache.conf /etc/apache2/conf.d/
ADD vhost.conf /etc/apache2/conf.d/ 

# Supervisor config
ADD run-apache.ini /etc/supervisor.d/
ADD run-mysql.ini /etc/supervisor.d/

# Temp for testing
ADD index.php /app/

ENTRYPOINT ["/bin/bash"]
