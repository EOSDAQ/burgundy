FROM golang:1.10 AS builder

# Download and install the latest release of dep
RUN go get -u github.com/golang/dep/cmd/dep
#RUN go get -u github.com/golang/lint/golint
#RUN go get -u github.com/sqs/goreturns
RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger
#RUN go get -u honnef.co/go/tools/cmd/megacheck

RUN apt-get install -y ca-certificates

ARG VERSION
ARG BUILD_DATE

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/burgundy
COPY . ./
RUN make vendor-update
RUN cd ./api && CGO_ENABLED=0 GOOS=linux go build -i -tags 'release' -a -installsuffix nocgo -ldflags "-X main.Version="$VERSION" -X main.BuildDate="$BUILD_DATE -o /burgundy .
RUN cd ./api && swagger generate spec -o /swagger.json
RUN cp -fp conf/.env.json /.env.json

FROM alpine
COPY --from=builder /burgundy ./
COPY --from=builder /swagger.json ./
COPY --from=builder /.env.json ./
ENV PORT 18889
EXPOSE 18889
ENTRYPOINT ["/burgundy"]
