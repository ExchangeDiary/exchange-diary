# https://github.com/brunoshiroma/go-gin-poc/blob/main/Dockerfile
# Build Step
FROM golang:alpine AS build

RUN apk update
RUN apk add make
RUN apk add git

WORKDIR /go/github.com/ExchangeDiary/exchange-diary
COPY . .
RUN go mod tidy
RUN GO111MODULE=on go build -ldflags="-s -w" -o exchange-diary ./application/cmd/main.go

# Final Step
FROM alpine as runtime

# Base packages
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --update tzdata
RUN rm -rf /var/cache/apk/*

WORKDIR /home
# Copy binary from build step
COPY --from=build /go/github.com/ExchangeDiary/exchange-diary/exchange-diary exchange-diary
# Copy config files to runtime
COPY --from=build /go/github.com/ExchangeDiary/exchange-diary/infrastructure infrastructure
# Define timezone
ENV TZ=Asia/Seoul

# Define the ENTRYPOINT

CMD [ "/home/exchange-diary", "-phase=sandbox" ]
# ENTRYPOINT ./exchange-diary

# # Document that the service listens on port 8080.
# EXPOSE 8080