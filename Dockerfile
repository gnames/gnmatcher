FROM alpine

LABEL maintainer="Dmitry Mozzherin"

WORKDIR /bin

COPY ./gnmatcher/gnmatcher /bin

ENTRYPOINT [ "gnmatcher" ]

CMD ["grpc", "-p", "8778"]
