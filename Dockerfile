FROM alpine

LABEL maintainer="Dmitry Mozzherin"

ENV LAST_FULL_REBUILD 2020-08-10

WORKDIR /bin

COPY ./gnmatcher/gnmatcher /bin

ENTRYPOINT [ "gnmatcher" ]
