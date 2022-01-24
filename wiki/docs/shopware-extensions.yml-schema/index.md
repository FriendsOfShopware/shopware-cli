---
title: 'Schema of .shopware-extensions.yml'
weight: 500 
---

# Objects
* [`.shopware-extension.yml`](#reference-config) (root object)
* [`build`](#reference-build)
* [`store`](#reference-store)
* [`StoreInfoFaqQuestion`](#reference-storeinfofaqquestion)


---------------------------------------
<a name="reference-config"></a>
## .shopware-extension.yml

**`.shopware-extension.yml` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**build**|`Build`||No|
|**store**|`Store`||No|

Additional properties are not allowed.

### Config.build

* **Type**: `Build`
* **Required**: No

### Config.store

* **Type**: `Store`
* **Required**: No




---------------------------------------
<a name="reference-build"></a>
## build

**`build` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**zip**|`object`||No|

Additional properties are not allowed.

### Build.zip

* **Type**: `object`
* **Required**: No




---------------------------------------
<a name="reference-shopware-cli"></a>
## shopware-cli

shopware cli extension configuration definition file



---------------------------------------
<a name="reference-store"></a>
## store

**`store` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**icon**|`string`|Specifies the Path to the icon (128x128 px) for store.|No|
|**availabilities**|`string` `[]`|Specifies the visibility in stores.|No|
|**localizations**|`string` `[]`|Specifies the languages the extension is translated.|No|
|**categories**|`string` `[*-2]`|Specifies the categories in which the extension can be found.|No|
|**default_locale**|`string`||No|
|**type**|`string`|Specifies the type of this extension.|No|
|**automatic_bugfix_version_compatibility**|`boolean`|Specifies whether the extension should automatically be set compatible with Shopware bugfix versions.|No|
|**videos**|`object`|Specifies the links of YouTube-Videos to show or describe the extension.|No|
|**tags**|`object`|Specifies the tags of the extension.|No|
|**highlights**|`object`|Specifies the highlights of the extension.|No|
|**features**|`object`|Specifies the features of the extension.|No|
|**faq**|`object`|Specifies Frequently Asked Questions for the extension.|No|
|**description**|`object`|Specifies the description of the extension in store.|No|
|**installation_manual**|`object`|Installation manual of the extension in store.|No|
|**images**|`object` `[1-*]`|Specifies images for the extension in the store.|No|

Additional properties are not allowed.

### Store.icon

Specifies the Path to the icon (128x128 px) for store.

* **Type**: `string`
* **Required**: No

### Store.availabilities

Specifies the visibility in stores.

* **Type**: `string` `[]`
    * Each element in the array must be unique.
    * Each element in the array must be one of the following values:
        * `German`
        * `International`
* **Required**: No

### Store.localizations

Specifies the languages the extension is translated.

* **Type**: `string` `[]`
    * Each element in the array must be unique.
    * Each element in the array must be one of the following values:
        * `de_DE`
        * `en_GB`
        * `es_ES`
        * `fi_FI`
        * `fr_FR`
        * `it_IT`
        * `nb_NO`
        * `nl_NL`
        * `pl_PL`
        * `sv_SE`
        * `bg_BG`
        * `cs_CZ`
        * `pt_PT`
        * `hy`
        * `de_CH`
        * `tr`
        * `da_DK`
        * `ru_RU`
* **Required**: No

### Store.categories

Specifies the categories in which the extension can be found.

* **Type**: `string` `[*-2]`
    * Each element in the array must be unique.
    * Each element in the array must be one of the following values:
        * `Administration`
        * `BackendBearbeitung`
        * `System`
        * `SEOOptimierung`
        * `Bonitaetsprüfung`
        * `Rechtssicherheit`
        * `MobileAdministration`
        * `Auswertung`
        * `KommentarFeedback`
        * `Tracking`
        * `MobileAuswertung`
        * `Integration`
        * `Shopsystem`
        * `PreissuchmaschinenPortale`
        * `Warenwirtschaft`
        * `Versand`
        * `Bezahlung`
        * `StorefrontDetailanpassungen`
        * `Sprache`
        * `Suche`
        * `HeaderFooter`
        * `Produktdarstellung`
        * `UebergreifendeDarstellung`
        * `Detailseite`
        * `MenueKategorien`
        * `Bestellprozess`
        * `KundenkontoPersonalisierung`
        * `IconsButons`
        * `Schriftarten`
        * `WidgetsSnippets`
        * `Sonderfunktionen`
        * `Themes`
        * `Branche`
        * `Home+Furnishings`
        * `FashionBekleidung`
        * `GartenNatur`
        * `KosmetikGesundheit`
        * `EssenTrinken`
        * `KinderPartyGeschenke`
        * `SportLifestyleReisen`
        * `TechnikIT`
        * `IndustrieGroßhandel`
        * `MigrationTools`
        * `Einkaufswelten`
        * `ConversionOptimierung`
        * `Extensions`
        * `MarketingTools`
        * `B2BExtensions`
        * `Blog`
* **Required**: No

### Store.default_locale

* **Type**: `string`
* **Required**: No
* **Allowed values**:
    * `"en_GB"`
    * `"de_DE"`

### Store.type

Specifies the type of this extension.

* **Type**: `string`
* **Required**: No
* **Allowed values**:
    * `"extension"`
    * `"theme"`

### Store.automatic_bugfix_version_compatibility

Specifies whether the extension should automatically be set compatible with Shopware bugfix versions.

* **Type**: `boolean`
* **Required**: No

### Store.videos

Specifies the links of YouTube-Videos to show or describe the extension.

* **Type**: `object`
* **Required**: No

### Store.tags

Specifies the tags of the extension.

* **Type**: `object`
* **Required**: No

### Store.highlights

Specifies the highlights of the extension.

* **Type**: `object`
* **Required**: No

### Store.features

Specifies the features of the extension.

* **Type**: `object`
* **Required**: No

### Store.faq

Specifies Frequently Asked Questions for the extension.

* **Type**: `object`
* **Required**: No

### Store.description

Specifies the description of the extension in store.

* **Type**: `object`
* **Required**: No

### Store.installation_manual

Installation manual of the extension in store.

* **Type**: `object`
* **Required**: No

### Store.images

Specifies images for the extension in the store.

* **Type**: `object` `[1-*]`
* **Required**: No




---------------------------------------
<a name="reference-storeinfofaqquestion"></a>
## StoreInfoFaqQuestion

**`StoreInfoFaqQuestion` Properties**

|   |Type|Description|Required|
|---|---|---|---|
|**question**|`string`|| &#10003; Yes|
|**answer**|`string`|| &#10003; Yes|

Additional properties are not allowed.

### StoreInfoFaqQuestion.question

* **Type**: `string`
* **Required**:  &#10003; Yes

### StoreInfoFaqQuestion.answer

* **Type**: `string`
* **Required**:  &#10003; Yes
