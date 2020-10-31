FROM golang:1.15.3 as golangBuild

RUN mkdir -p /output
WORKDIR /output
COPY ./  ./
RUN GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor -o 91porn

FROM senlixiushu/91porn-runtime:0.0.1
MAINTAINER linyuan

# RUN mkdir -p /91porn && mkdir -p /aria2 && mkdir -p /run/nginx

COPY ./conf /91porn/conf
COPY --from=golangBuild /output/91porn /91porn

CMD /91porn/91porn 
