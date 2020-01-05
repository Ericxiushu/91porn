FROM alpine:latest
MAINTAINER linyuan

RUN mkdir -p /91porn

COPY ./91porn /91porn
COPY ./conf /91porn/conf

RUN apk add aria2 && apk add nginx

RUN mkdir -p /aria2/aria2ng && mkdir -p /run/nginx
COPY ./aria2.conf /aria2/aria2.conf
COPY ./nginx/nginx.conf /etc/nginx/nginx.conf
COPY ./nginx/nginx_default.conf /etc/nginx/conf.d/default.conf
COPY ./AriaNg-1.1.4.zip /aria2/AriaNg-1.1.4.zip
RUN cd /aria2 && touch aria2.session && unzip AriaNg-1.1.4.zip -d ./aria2ng
RUN cd /run/nginx && touch nginx.pid

CMD nginx && aria2c --conf-path=/aria2/aria2.conf -D && /91porn/91porn 
