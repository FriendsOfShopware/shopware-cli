FROM alpine
RUN apk add --no-cache git nodejs npm
ENTRYPOINT ["/shopware-cli"]
COPY shopware-cli /
