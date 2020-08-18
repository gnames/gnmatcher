FROM alpine

LABEL maintainer="Dmitry Mozzherin"

# RUN apk add --no-cache bash

WORKDIR /bin

COPY ./gnmatcher/gnmatcher /bin

ENTRYPOINT [ "gnmatcher" ]

CMD ["grpc", "-p", "8778"]
