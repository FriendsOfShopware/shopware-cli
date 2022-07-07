FROM php:7.4-cli-alpine

LABEL org.opencontainers.image.source https://github.com/FriendsOfShopware/shopware-cli
COPY --from=mlocati/php-extension-installer /usr/bin/install-php-extensions /usr/bin/
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

RUN apk add --no-cache git nodejs npm && install-php-extensions bcmath gd intl mysqli pdo_mysql sockets bz2 gmp soap zip gmp pcntl posix redis pcov imagick xsl calendar

COPY shopware-cli /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/shopware-cli"]
