---
title: App Commands
weight: 20
--

### app push

This command uploads app files to a shopware instance and installs and activates the app over the api.

* url - Url pointing to a shop, for example `http://localhost:8000`. **Required**
* dir - Path to the app folder. Default: `./`
* credentials - Path to a json file containing user crendentials
  * For admin users: `{ "Username": "admin", "Password": "shopware" }`
  * For integration: `{ "Id": "<client-id>", "Secret": "<client-secret>" }`
