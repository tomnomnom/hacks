FROM alpine

RUN apk --update add \
    bash \
    supervisor \
    php7-apache2 \
    php7-mysqlnd \
    php7-mysqli \
    mysql \
    mysql-client \
    curl

# For whatever reason the apk package doesn't create
# this dir so apache fails to start
RUN mkdir /run/apache2

# Guess what this one's for...
RUN mkdir /run/mysqld
RUN chown mysql: /run/mysqld
RUN su mysql -s /bin/sh -c 'mysql_install_db --datadir=/var/lib/mysql'

# Wordpress files
RUN mkdir /app
RUN curl -s https://en-gb.wordpress.org/latest-en_GB.tar.gz -o /app/wordpress.tgz
RUN tar xvzf /app/wordpress.tgz -C /app/ 
RUN chmod -R 777 /app/wordpress/wp-content

# Wordpress config
ADD wp-config.php /app/wordpress/
ADD wordpress.sql /app/
ADD configure-mysql.sh /app/
RUN /app/configure-mysql.sh

# Apache config
ADD apache.conf /etc/apache2/conf.d/
ADD vhost.conf /etc/apache2/conf.d/ 

# Supervisor config
ADD run-apache.ini /etc/supervisor.d/
ADD run-mysql.ini /etc/supervisor.d/

# PHP info file for testing
ADD index.php /app/

# The -n makes supervisord run in the foreground
ENTRYPOINT ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisord.conf"]
