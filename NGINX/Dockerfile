FROM nginx

COPY ./default-http.conf /etc/nginx/conf.d/default.conf
COPY ./ssl.crt /etc/nginx/ssl/domain_pack.crt
COPY ./ssl_private.key /etc/nginx/ssl/private.key

USER $user
EXPOSE 80