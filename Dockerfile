ARG PHP_VERSION

FROM ghcr.io/friendsofshopware/shopware-cli-base:phpv${PHP_VERSION}

COPY shopware-cli /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/entrypoint.sh", "shopware-cli"]
CMD ["--help"]
