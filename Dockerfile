# BUILD
FROM golang:latest as builder

# copy
WORKDIR /src/
COPY . /src/

# build
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags netgo \
    -ldflags '-w -extldflags "-static"' \
    -mod vendor \
    -o app

# RUN
FROM gcr.io/distroless/static
COPY --from=builder /src/app .
ENTRYPOINT ["/app"]