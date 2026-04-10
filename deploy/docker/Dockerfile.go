ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
ARG APP=api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /out/flowscope ./cmd/${APP}

FROM alpine:3.20
RUN adduser -D -u 10001 flowscope
USER flowscope
COPY --from=build /out/flowscope /usr/local/bin/flowscope
ENTRYPOINT ["/usr/local/bin/flowscope"]
