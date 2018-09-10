FROM alpine

RUN apk add --no-cache ca-certificates

ARG BUILD_PORT

COPY bin/burgundy ./
COPY bin/swagger.json ./
COPY conf ./
ENV PORT $BUILD_PORT
EXPOSE $BUILD_PORT
ENTRYPOINT ["/burgundy"]
