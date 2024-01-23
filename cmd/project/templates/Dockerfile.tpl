#syntax=docker/dockerfile:1.4

# pin versions
FROM shopware/docker-base:{{.PHP.PhpVersion}} as base-image
FROM ghcr.io/friendsofshopware/shopware-cli:latest-php-{{.PHP.PhpVersion}} as shopware-cli

FROM base-image as base-extended

USER root

{{- if .PHP.Extensions }}
RUN /usr/local/bin/install-php-extensions {{- range $key, $value := .PHP.Extensions }} {{$value}} {{- end }}
{{- end }}
{{- if .PHP.Settings }}

COPY <<EOF /usr/local/etc/php/conf.d/99-z-custom.ini
{{- range $key, $value := .PHP.Settings }}

{{$key}} = {{$value}}

{{- end }}
EOF

{{- end }}

USER www-data

FROM shopware-cli as build

ADD . /src
WORKDIR /src

{{- if .BuildEnv }}
ENV {{ .BuildEnv }}
{{- end }}

RUN --mount=type=secret,id=composer_auth,dst=/src/auth.json \
    --mount=type=cache,target=/root/.composer \
    --mount=type=cache,target=/root/.npm \
    /usr/local/bin/entrypoint.sh shopware-cli project ci /src

FROM base-extended

{{- if .RunEnv }}
ENV {{ .RunEnv }}
{{- end }}

COPY --from=build --chown=www-data --link /src /var/www/html