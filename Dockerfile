FROM golang:1.16-alpine AS build
WORKDIR /go/src/app
COPY . .
RUN go mod tidy && go run build/mage.go install

FROM alpine:latest
RUN apk add --update curl
COPY --from=build /go/bin/scheduler /usr/bin/scheduler
LABEL org.opencontainers.image.authors="alfreddobradi@gmail.com"
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/health || exit 1
EXPOSE 80
ENTRYPOINT ["/usr/bin/scheduler"]