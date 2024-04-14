FROM cgr.dev/chainguard/wolfi-base

RUN <<EOF
    set -eo pipefail

apk add --no-cache \
    libstdc++ \
    sudo-rs \
    nodejs-21 \
    curl \
    npm \
    go \
    golangci-lint \
    openssh \
    openssh-keygen \
    openssh-client \
    git \
    git-lfs \
    composer \
    php-8.2 \
    php-8.2-fileinfo \
    php-8.2-openssl \
    php-8.2-ftp \
    php-8.2-ctype \
    php-8.2-xml \
    php-8.2-dom \
    php-8.2-simplexml \
    php-8.2-xmlreader \
    php-8.2-xmlwriter \
    php-8.2-bcmath \
    php-8.2-iconv \
    php-8.2-gd \
    php-8.2-intl \
    php-8.2-pdo \
    php-8.2-pdo_mysql \
    php-8.2-mysqlnd \
    php-8.2-pcntl \
    php-8.2-sockets \
    php-8.2-bz2 \
    php-8.2-gmp \
    php-8.2-soap \
    php-8.2-zip \
    php-8.2-redis \
    php-8.2-opcache \
    php-8.2-amqp

# Refresh ldconfig cache
ldconfig
echo '%gitpod ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
addgroup -g 33333 gitpod
adduser -u 33333 -G gitpod -h /home/gitpod -s /bin/bash -D gitpod
EOF

USER gitpod

