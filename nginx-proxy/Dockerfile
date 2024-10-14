FROM nginx

ARG PROXY_HOST go-user
ARG PORT 8888

COPY default.conf /etc/nginx/conf.d

RUN sed -i "s|#PROXY_HOST#|${PROXY_HOST}|;s|#PORT#|${PORT}|" /etc/nginx/conf.d/default.conf
