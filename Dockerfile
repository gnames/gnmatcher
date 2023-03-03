FROM alpine:3.17

LABEL maintainer="Dmitry Mozzherin"

# RUN apk add --no-cache bash

WORKDIR /bin

ENV LANG en_US.UTF-8
ENV LC_COLLATE C

COPY ./gnmatcher /bin

ENTRYPOINT [ "gnmatcher" ]

CMD ["rest", "-p", "8080"]
