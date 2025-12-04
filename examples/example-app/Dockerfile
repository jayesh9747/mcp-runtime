FROM golang:1.22-alpine AS build
ENV CGO_ENABLED=0 GO111MODULE=off
WORKDIR /src
COPY main.go .
RUN go build -o server main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=build /src/server /app/server
EXPOSE 8088
ENTRYPOINT ["/app/server"]
