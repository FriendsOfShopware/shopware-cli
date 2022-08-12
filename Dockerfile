FROM php:7.4-cli-debian

LABEL org.opencontainers.image.source https://github.com/FriendsOfShopware/shopware-cli
COPY --from=mlocati/php-extension-installer /usr/bin/install-php-extensions /usr/bin/
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash - && \
    apt-get install -y git nodejs && \
    install-php-extensions bcmath gd intl mysqli pdo_mysql sockets bz2 soap zip gmp pcntl redis pcov imagick xsl calendar

COPY shopware-cli /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/shopware-cli"]
