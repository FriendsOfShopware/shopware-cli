---
title: Bundle Support
---

Since 0.3.0 it's possible to use shopware-cli to build assets of normal Shopware bundles (base class does not extend Plugin and is registered in `config/bundles.php`). If you want to know more about Shopware bundles vs Plugins, checkout [this blog post](https://shyim.me/blog/you-dont-need-a-plugin-to-customize-shopware-6/).

## Bundle directly embeded into the Shopware project

If your Bundle is directly included in the Shopware project without a own composer.json, you have to adjust the root `composer.json` to make the bundle available for shopware-cli.

```json
{
    "extra": {
        "shopware-bundles": {
            // The key is the relative path from project root to the bundle
            "src/MyBundle": {}
        }
    }
}
```

If your bundle folder names does not match your Bundle name, you can use the `name` key to map the folder to the bundle name.

```json
{
    "extra": {
        "shopware-bundles": {
            "src/MyBundle": {
                "name": "MyFancyBundle"
            }
        }
    }
}
```

## Bundle as composer package

If your bundle is a own composer package, make sure your composer type is `shopware-bundle` and that you have set a `shopware-bundle-name` in the `extra` part of the config like this:

```json
{
    "name": "my-vendor/my-bundle",
    "type": "shopware-bundle",
    "extra": {
        "shopware-bundle-name": "MyBundle"
    }
}
```

Now you can use `shopware-cli extension build <path>` to build the assets and distribute them together with your bundle. Also `shopware-cli project ci` detects know automatically this bundle and builds the assets for it.
