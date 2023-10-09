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

Creates a zip file from extension folder

Parameters:

* path - Path to extension folder. F.e: `shopware-cli extension zip MyPlugin`

Environment-Variables:

* SHOPWARE_PROJECT_ROOT (optional) - Path to a installed shopware to speed up building. F.e: `SHOPWARE_PROJECT_ROOT=/var/www/myshop/ shopware-cli extension zip MyPlugin`


## shopware-cli extension build

Builds the JS and CSS assets into the extension folder

Parameters:

* path - Path to extension folder. This can be also multiple directories. F.e: `SHOPWARE_PROJECT_ROOT=/var/www/myshop/ shopware-cli extension build MyPlugin MySecondPlugin`

Environment-Variables:

* SHOPWARE_PROJECT_ROOT (optional) - Path to a installed shopware to speed up building. F.e: `SHOPWARE_PROJECT_ROOT=/var/www/myshop/ shopware-cli extension build MyPlugin`


## shopware-cli extension admin-watch

Starts a admin-watcher using ESBuild to build the JS and CSS assets into the extension folder

Parameters:

* path - Path to extension folder.
* url - A URL of a running Shopware instance. This is used to proxy the requests to the admin API like http://localhost

Options:

* `--listen` - Listen Address for Server
* `--external-url` - Use this URL in the browser. Needed for reverse proxy setups
