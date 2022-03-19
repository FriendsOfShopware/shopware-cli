---
 title: 'Schema of .shopware-project.yml' 
---

# Objects
* [`.shopware-project.yml`](#reference-config) (root object)
* [`Admin API credentials`](#reference-adminapi)
* [`Entity Sync Filter`](#reference-entitysyncfilterinner)
* [`MySQL dump configuration`](#reference-dump)
* [`Sync Settings`](#reference-sync)
    * [`Entity Sync`](#reference-entitysyncitem)
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
<a name="reference-entitysyncitem"></a>
## Entity Sync

**`Entity Sync` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**entity**|`string`|| &#10003; Yes|
|**exists**|`EntitySyncFilter` `[]`||No|
|**payload**|`object`|API payload| &#10003; Yes|

Additional properties are not allowed.

### EntitySyncItem.entity

* **Type**: `string`
* **Required**:  &#10003; Yes

### EntitySyncItem.exists

* **Type**: `EntitySyncFilter` `[]`
* **Required**: No

### EntitySyncItem.payload

API payload

* **Type**: `object`
* **Required**:  &#10003; Yes




---------------------------------------
<a name="reference-entitysyncfilterinner"></a>
## Entity Sync Filter

**`Entity Sync Filter` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**type**|`string`|filter type| &#10003; Yes|
|**field**|`string`|field| &#10003; Yes|
|**value**|`["string", "integer", "array", "boolean", "null"]`|value|No|
|**operator**|`string`||No|

Additional properties are not allowed.

### EntitySyncFilterInner.type

filter type

* **Type**: `string`
* **Required**:  &#10003; Yes
* **Allowed values**:
    * `"equals"`
    * `"multi"`
    * `"contains"`
    * `"prefix"`
    * `"suffix"`
    * `"not"`
    * `"range"`
    * `"until"`
    * `"equalsAll"`
    * `"equalsAny"`

### EntitySyncFilterInner.field

field

* **Type**: `string`
* **Required**:  &#10003; Yes

### EntitySyncFilterInner.value

value

* **Type**: `["string", "integer", "array", "boolean", "null"]`
* **Required**: No

### EntitySyncFilterInner.operator

* **Type**: `string`
* **Required**: No
* **Allowed values**:
    * `"AND"`
    * `"OR"`
    * `"XOR"`




---------------------------------------
<a name="reference-mailtemplateitemtranslation"></a>
## Mail Template Single Translation

**`Mail Template Single Translation` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**language**|`string`||No|
|**sender_name**|`string`||No|
|**subject**|`string`||No|
|**html**|`string`||No|
|**plain**|`string`||No|
|**custom_fields**|`["object", "null"]`||No|

Additional properties are not allowed.

### MailTemplateItemTranslation.language

* **Type**: `string`
* **Required**: No

### MailTemplateItemTranslation.sender_name

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

### MailTemplateItemTranslation.custom_fields

* **Type**: `["object", "null"]`
* **Required**: No




---------------------------------------
<a name="reference-mailtemplateitem"></a>
## Mail Template Sync

**`Mail Template Sync` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**id**|`string`||No|
|**translations**|`MailTemplateItemTranslation` `[]`||No|

Additional properties are not allowed.

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
|**entity**|`EntitySyncItem` `[]`||No|

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

### Sync.entity

* **Type**: `EntitySyncItem` `[]`
* **Required**: No




---------------------------------------
<a name="reference-syncconfigitem"></a>
## System Config Sync

**`System Config Sync` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**sales_channel**|`string`||No|
|**settings**|`object`|| &#10003; Yes|

Additional properties are not allowed.

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

Additional properties are not allowed.

### ThemeConfigItem.name

* **Type**: `string`
* **Required**: No

### ThemeConfigItem.settings

* **Type**: `object`
* **Required**: No
