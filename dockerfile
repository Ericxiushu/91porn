FROM alpine:latest
MAINTAINER linyuan

RUN mkdir -p /91porn

COPY ./91porn /91porn
COPY ./conf /91porn/conf

CMD /91porn/91porn
