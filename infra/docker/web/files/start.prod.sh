#!/bin/bash

echo "Replacing \"{{API_ADDR}}\" in /etc/nginx/nginx.conf with $API_SVC_SERVICE_HOST..."
sed -i "s/{{API_ADDR}}/${API_SVC_SERVICE_HOST}/g;" /etc/nginx/nginx.conf
echo "Replacing \"{{ROUTER_ADDR}}\" in /etc/nginx/nginx.conf with $ROUTER_SVC_SERVICE_HOST..."
sed -i "s/{{ROUTER_ADDR}}/${ROUTER_SVC_SERVICE_HOST}/g;" /etc/nginx/nginx.conf
echo "Replacing \"{{GOPHR_CERT_SECRET}}\" in /etc/nginx/nginx.conf with cert.prod.crt..."
sed -i "s/{{GOPHR_CERT_SECRET}}/cert.prod.crt/g;" /etc/nginx/nginx.conf
echo "Replacing \"{{GOPHR_KEY_SECRET}}\" in /etc/nginx/nginx.conf with cert.prod.key..."
sed -i "s/{{GOPHR_KEY_SECRET}}/cert.prod.key/g;" /etc/nginx/nginx.conf
echo "Starting nginx..."
nginx -g 'daemon off;'
