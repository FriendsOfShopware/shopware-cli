---
title: Install
weight: 9
---

You can install the pre-compiled binary (in several different ways), use Docker or compile from source.

Below you can find the steps for each of them.


## Install the pre-compiled binary

### Homebrew

```bash
brew install FriendsOfShopware/tap/shopware-cli
```

### deb,rpm apt packages

Download the .deb, .rpm or .apk packages from the [releases](https://github.com/FriendsOfShopware/shopware-cli/releases/) page and install them with the appropriate tools.

### go install

```bash
go install github.com/FriendsOfShopware/shopware-cli@latest
```

## manually

Download the pre-compiled binaries from the [releases](https://github.com/FriendsOfShopware/shopware-cli/releases/) page and copy them to the desired location.

## Running with Docker

You can also use it within a Docker container. To do that, you'll need to execute something more-or-less like the examples below.

Registries:

- [ghcr.io/friendsofshopware/shopware-cli](https://github.com/FriendsOfShopware/shopware-cli/pkgs/container/shopware-cli)

Example usage:

Builds assets of an extension

```
docker run \
    --rm \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    -u $(id -u) \
    ghcr.io/friendsofshopware/shopware-cli \
    extension build FroshPlatformAdminer
```

## Compiling from source

If you just want to build from source for whatever reason, follow these steps:

### clone:

```
git clone https://github.com/FriendsOfShopware/shopware-cli
cd shopware-cli
```

### get the dependencies:

```
go mod tidy
```

### build:

```
go build -o shopware-cli .
```

### verify it works:

```
./shopware-cli --version
```