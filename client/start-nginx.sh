#!/bin/sh

envsubst '${UPSTREAM_SERVER} ${SERVER_NAME}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

exec nginx -g 'daemon off;'