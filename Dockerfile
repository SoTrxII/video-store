# Dockerfile for an appservice go container (expected size 18Mo)
####################################################################################################
## Builder
####################################################################################################
FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
# Build the app, strip it (LDFLAGS) and optimize it with UPX
RUN apk add upx && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/server -ldflags="-w -s" && \
    upx --best --lzma ./bin/server

####################################################################################################
## Final image
####################################################################################################
FROM alpine as release
# The runtime user, having no home dir nor password
RUN adduser -HD -s /bin/ash appuser

WORKDIR /app
# Copy the built app, only allowing our app user to execute it
COPY --from=builder --chmod=0500 --chown=appuser:appuser  /app/bin/server ./
EXPOSE 8080
ENTRYPOINT [ "/app/server" ]