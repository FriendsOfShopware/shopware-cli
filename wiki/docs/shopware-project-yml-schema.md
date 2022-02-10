---
 title: 'Schema of .shopware-project.yml' 
---

# Objects
* [`.shopware-project.yml`](#reference-config) (root object)
* [`admin api`](#reference-adminapi)


---------------------------------------
<a name="reference-config"></a>
## .shopware-project.yml

**`.shopware-project.yml` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**url**|`string`|URL to Shopware instance|No|
|**admin_api**|`AdminApi`||No|

Additional properties are not allowed.

### Config.url

URL to Shopware instance

* **Type**: `string`
* **Required**: No

### Config.admin_api

* **Type**: `AdminApi`
* **Required**: No




---------------------------------------
<a name="reference-adminapi"></a>
## admin api

**`admin api` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**client_id**|`string`|Client ID of integreation|No|
|**client_secret**|`string`|Client Secret of integreation|No|
|**username**|`string`|Username of admin user|No|
|**password**|`string`|Password of admin user|No|

Additional properties are not allowed.

### AdminApi.client_id

Client ID of integreation

* **Type**: `string`
* **Required**: No

### AdminApi.client_secret

Client Secret of integreation

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
<a name="reference-shopware-cli"></a>
## shopware-cli

shopware cli project configuration definition file
