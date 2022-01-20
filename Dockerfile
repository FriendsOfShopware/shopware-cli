FROM alpine
RUN apk add --no-cache git nodejs npm
ENTRYPOINT ["/usr/local/bin/shopware-cli"]
COPY shopware-cli /usr/local/bin/
