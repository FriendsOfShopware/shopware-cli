ARG PHP_VERSION

FROM ghcr.io/friendsofshopware/shopware-cli-base:${PHP_VERSION}

COPY shopware-cli /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/shopware-cli"]
CMD ["--help"]
