---
title: Extension Commands
weight: 20
---

## shopware-cli extension validate

Validate extension for store compliance. Supported PHP Versions are: 7.3, 7.4, 8.1, 8.2

Parameters:

* path - Path to zip or extension folder


## shopware-cli extension prepare

Installs composer dependencies of the extension

Parameters:

* path - Path to extension folder


## shopware-cli extension zip

Creates a zip file from the extension folder.

Parameters:

* path - Path to extension folder. For example: `shopware-cli extension zip MyPlugin`.

Environment-Variables:

* SHOPWARE_PROJECT_ROOT (optional) - Path to an installed Shopware to speed up building. For example: `SHOPWARE_PROJECT_ROOT=/var/www/myshop/ shopware-cli extension zip MyPlugin`.


## shopware-cli extension build

Builds the JS and CSS assets into the extension folder

Parameters:

* path - Path to extension folder. This can be also multiple directories. For example: `SHOPWARE_PROJECT_ROOT=/var/www/myshop/ shopware-cli extension build MyPlugin MySecondPlugin`.

Environment-Variables:

* SHOPWARE_PROJECT_ROOT (optional) - Path to an installed Shopware to speed up building. For example: `SHOPWARE_PROJECT_ROOT=/var/www/myshop/ shopware-cli extension build MyPlugin`.


## shopware-cli extension admin-watch

Starts an admin-watcher using ESBuild to build the JS and CSS assets into the extension folder.

Parameters:

* path - Path to extension folder.
* url - A URL of a running Shopware instance. This is used to proxy the requests to the admin API like http://localhost

Options:

* `--listen` - Listen Address for Server
* `--external-url` - Use this URL in the browser. Needed for reverse proxy setups

## shopware-cli extension get-changelog

Get the changelog of an extension

Arguments:

* `path` - Path to extension folder/zip

Parameters:

* `--german` - Get the German changelog.

## shopware-cli extension get-version

Get the version of an extension

Arguments:

* `path` - Path to extension folder/zip

Parameters:

* `--german` - Get the German changelog.