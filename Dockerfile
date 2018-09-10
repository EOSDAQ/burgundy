FROM alpine

ENV TZ 'Asia/Seoul'

RUN apk add --no-cache ca-certificates && \
        echo $TZ > /etc/timezone && \
        apt-get update && \
        apt-get install -y tzdata && \
        rm /etc/localtime && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
        dpkg-reconfigure -f noninteractive tzdata && \
        apt-get clean

ARG BUILD_PORT

COPY bin/burgundy ./
COPY bin/swagger.json ./
COPY conf ./
ENV PORT $BUILD_PORT
EXPOSE $BUILD_PORT
ENTRYPOINT ["/burgundy"]
