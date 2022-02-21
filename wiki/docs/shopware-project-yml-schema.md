---
 title: 'Schema of .shopware-project.yml' 
---

# Objects
* [`.shopware-project.yml`](#reference-config) (root object)
* [`Admin API credentials`](#reference-adminapi)
* [`MySQL dump configuration`](#reference-dump)
* [`Sync Settings`](#reference-sync)
    * [`Mail Template Sync`](#reference-mailtemplateitem)
        * [`Mail Template Single Translation`](#reference-mailtemplateitemtranslation)
    * [`System Config Sync`](#reference-syncconfigitem)
    * [`Theme Config Sync`](#reference-themeconfigitem)


---------------------------------------
<a name="reference-config"></a>
## .shopware-project.yml

**`.shopware-project.yml` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**url**|`string`|URL to Shopware instance|No|
|**admin_api**|`AdminApi`||No|
|**dump**|`Dump`||No|
|**sync**|`Sync`||No|

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

### Config.sync

* **Type**: `Sync`
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
<a name="reference-mailtemplateitemtranslation"></a>
## Mail Template Single Translation

**`Mail Template Single Translation` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**language**|`string`||No|
|**senderName**|`string`||No|
|**subject**|`string`||No|
|**html**|`string`||No|
|**plain**|`string`||No|
|**customFields**|`object`||No|

Additional properties are allowed.

### MailTemplateItemTranslation.language

* **Type**: `string`
* **Required**: No

### MailTemplateItemTranslation.senderName

* **Type**: `string`
* **Required**: No

### MailTemplateItemTranslation.subject

* **Type**: `string`
* **Required**: No

### MailTemplateItemTranslation.html

* **Type**: `string`
* **Required**: No

### MailTemplateItemTranslation.plain

* **Type**: `string`
* **Required**: No

### MailTemplateItemTranslation.customFields

* **Type**: `object`
* **Required**: No




---------------------------------------
<a name="reference-mailtemplateitem"></a>
## Mail Template Sync

**`Mail Template Sync` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**id**|`string`||No|
|**translations**|`MailTemplateItemTranslation` `[]`||No|

Additional properties are allowed.

### MailTemplateItem.id

* **Type**: `string`
* **Required**: No

### MailTemplateItem.translations

* **Type**: `MailTemplateItemTranslation` `[]`
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



---------------------------------------
<a name="reference-sync"></a>
## Sync Settings

**`Sync Settings` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**config**|`SyncConfigItem` `[]`||No|
|**theme**|`ThemeConfigItem` `[]`||No|
|**mail_template**|`MailTemplateItem` `[]`||No|

Additional properties are not allowed.

### Sync.config

* **Type**: `SyncConfigItem` `[]`
* **Required**: No

### Sync.theme

* **Type**: `ThemeConfigItem` `[]`
* **Required**: No

### Sync.mail_template

* **Type**: `MailTemplateItem` `[]`
* **Required**: No




---------------------------------------
<a name="reference-syncconfigitem"></a>
## System Config Sync

**`System Config Sync` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**sales_channel**|`string`||No|
|**settings**|`object`|| &#10003; Yes|

Additional properties are allowed.

### SyncConfigItem.sales_channel

* **Type**: `string`
* **Required**: No

### SyncConfigItem.settings

* **Type**: `object`
* **Required**:  &#10003; Yes




---------------------------------------
<a name="reference-themeconfigitem"></a>
## Theme Config Sync

**`Theme Config Sync` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**name**|`string`||No|
|**settings**|`object`||No|

Additional properties are allowed.

### ThemeConfigItem.name

* **Type**: `string`
* **Required**: No

### ThemeConfigItem.settings

* **Type**: `object`
* **Required**: No
