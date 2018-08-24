FROM golang:1.10 AS builder

RUN apt-get update -y
RUN apt-get install -y ca-certificates
RUN apt-get upgrade -y ca-certificates
RUN update-ca-certificates

# Download and install the latest release of dep
RUN go get -u github.com/golang/dep/cmd/dep
#RUN go get -u github.com/golang/lint/golint
#RUN go get -u github.com/sqs/goreturns
RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger
#RUN go get -u honnef.co/go/tools/cmd/megacheck

ARG VERSION
ARG BUILD_DATE
ARG BUILD_PKG
ARG BUILD_PORT

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/$BUILD_PKG
COPY . ./
RUN make vendor
RUN cd ./api && CGO_ENABLED=0 GOOS=linux go build -i -tags 'release' -a -installsuffix nocgo -ldflags "-X main.Version="$VERSION" -X main.BuildDate="$BUILD_DATE -o /$BUILD_PKG .
RUN cd ./api && swagger generate spec -o /swagger.json
RUN cp -fp conf/.env.json /.env.json

FROM alpine
ARG BUILD_PKG
ARG BUILD_PORT
COPY --from=builder /$BUILD_PKG ./
COPY --from=builder /swagger.json ./
COPY --from=builder /.env.json ./
ENV PORT $BUILD_PORT
EXPOSE $BUILD_PORT
ENTRYPOINT ["/burgundy"]
