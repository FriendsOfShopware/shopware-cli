---
title: Extension Commands
weight: 20
---

## shopware-cli extension validate

Validate extension for store compliance. The PHP code will be sent to an [external service](https://github.com/FriendsOfShopware/aws-php-syntax-checker-lambda) to verify the PHP syntax.

Parameters:

* path - Path to zip or extension folder


## shopware-cli extension prepare

Installs composer dependencies of the extension

Parameters:

* path - Path to extension folder


## shopware-cli extension zip

Creates a zip file from extension folder

Parameters:

* path - Path to extension folder

## shopware-cli extension build

Builds the JS and CSS assets into the extension folder

Parameters:

* path - Path to extension folder. This can be also multiple directories


## shopware-cli extension admin-watch

Starts a admin-watcher using ESBuild to build the JS and CSS assets into the extension folder

Parameters:

* path - Path to extension folder.
* url - A URL of a running Shopware instance. This is used to proxy the requests to the admin API like http://localhost