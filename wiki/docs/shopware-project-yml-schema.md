---
 title: 'Schema of .shopware-project.yml' 
---

# Objects
* [`.shopware-project.yml`](#reference-config) (root object)
* [`Admin API credentials`](#reference-adminapi)
* [`MySQL dump configuration`](#reference-dump)


---------------------------------------
<a name="reference-config"></a>
## .shopware-project.yml

**`.shopware-project.yml` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**url**|`string`|URL to Shopware instance|No|
|**admin_api**|`AdminApi`||No|
|**dump**|`Dump`||No|

Additional properties are not allowed.

### Config.url

URL to Shopware instance

* **Type**: `string`
* **Required**: No

### Config.admin_api

* **Type**: `AdminApi`
* **Required**: No

### Config.dump

* **Type**: `Dump`
* **Required**: No




---------------------------------------
<a name="reference-adminapi"></a>
## Admin API credentials

**`Admin API credentials` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**client_id**|`string`|Client ID of integration|No|
|**client_secret**|`string`|Client Secret of integration|No|
|**username**|`string`|Username of admin user|No|
|**password**|`string`|Password of admin user|No|

Additional properties are not allowed.

### AdminApi.client_id

Client ID of integration

* **Type**: `string`
* **Required**: No

### AdminApi.client_secret

Client Secret of integration

* **Type**: `string`
* **Required**: No

### AdminApi.username

Username of admin user

* **Type**: `string`
* **Required**: No

### AdminApi.password

Password of admin user

* **Type**: `string`
* **Required**: No




---------------------------------------
<a name="reference-dump"></a>
## MySQL dump configuration

**`MySQL dump configuration` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**rewrite**|`object`||No|
|**nodata**|`string` `[]`||No|
|**ignore**|`string` `[]`||No|
|**where**|`object`||No|

Additional properties are not allowed.

### Dump.rewrite

* **Type**: `object`
* **Required**: No

### Dump.nodata

* **Type**: `string` `[]`
* **Required**: No

### Dump.ignore

* **Type**: `string` `[]`
* **Required**: No

### Dump.where

* **Type**: `object`
* **Required**: No




---------------------------------------
<a name="reference-shopware-cli"></a>
## shopware-cli

shopware cli project configuration definition file
