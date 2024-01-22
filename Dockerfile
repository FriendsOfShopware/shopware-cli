#syntax=docker/dockerfile:1.4

# pin versions
FROM shopware/docker-base:8.2 as base-image
FROM ghcr.io/friendsofshopware/shopware-cli:latest-php-8.2 as shopware-cli

FROM base-image as base-extended
RUN install-php-extensions opentelemetry

RUN echo "display_errors = 1" >> /usr/local/etc/php/conf.d/99-z-custom.ini


FROM shopware-cli as build

ADD . /src
WORKDIR /src

RUN --mount=type=secret,id=composer_auth,dst=/src/auth.json \
    --mount=type=cache,target=/root/.composer \
    --mount=type=cache,target=/root/.npm \
    /usr/local/bin/entrypoint.sh shopware-cli project ci /src

FROM base-extended

COPY --from=build --chown=www-data --link /src /var/www/html