FROM alpine:3.13.4

LABEL Author="rangh"
LABEL Application="magicLite"

RUN mkdir -p /var/app
RUN mkdir -p /var/app/config

WORKDIR /var/app

ADD magicLite /var/app/magicLite
ADD bootstrap.sh /var/app/bootstrap.sh

RUN chmod +x /var/app/magicLite
RUN chmod +x /var/app/bootstrap.sh

ENTRYPOINT [ "/var/app/bootstrap.sh" ]

EXPOSE 8880
