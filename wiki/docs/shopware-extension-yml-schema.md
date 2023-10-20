---
 title: 'Schema of .shopware-extension.yml' 
---

Any configuration field is optional. When you create a `.shopware-extension.yml`, you get also IDE autocompletion for all fields.


```yaml
# .shopware-extension.yml
build:
    # override the auto detection of the shopware constraint
    shopwareVersionConstraint: `~6.5.0`

    # build additional bundles for assets
    extraBundles:
        # when only path passed the folder name will be used as bundle name
        - path: MySecondBundle
        # or specify it explicitly
        - name: DifferentName
          path: Bundle
    zip:
        composer:
            # disable composer install, enabled by default
            enabled: true

            # run commands before composer install or after
            before_hooks:
                - echo "Before"
            after_hooks:
                - echo "After"

            # exclude packages from installing into the vendor folder of the plugin
            excluded_packages:
                - symfony/filesystem
        assets:
            # disable building assets, builds by default
            enabled: true

            # run commands before asset building or after
            before_hooks:
                - echo "Before"
            after_hooks:
                - echo "After"

            # use bundled esbuild instead of default compile
            enable_es_build_for_admin: false

            # use bundled esbuild instead of default compile
            enable_es_build_for_storefront: false

        pack:
            # run commands before packing the zip
            before_hooks:
                - echo "Before"

            # paths to skip for zipping
            excludes:
                paths:
                    - .idea

store:
    # override default icon path
    icon: icon.png

    # store availabilities
    availabilities:
        - German
        - International

    # store localizations, see json schema for all names
    localizations:
        - de_DE
        - en_GB

    # default locale
    default_locale: en_GB

    # store category, see json schema for all names
    categories:
        - Administration

    # type: extension or theme
    type: extension

    # auto mark patch versions as compatible
    automatic_bugfix_version_compatibility: true

    videos:
        en:
            - https://yotuube.com/...
        de:
            - ....

    tags:
        en:
            - tag
        de:
            - ...

    highlights:
        en:
            - tag
        de:
            - ...

    features:
        en:
            - tag
        de:
            - ...

    faq:
        - question: Can do the extension this ?
          answer: Yes, we can ....


    description:
        # inline
        en: |
            Write inline
        # embedd an html or markdown file
        de: file:src/Resources/store/description.md

    installation_manual:
        # inline
        en: |
            Write inline
        # embedd an html or markdown file
        de: file:src/Resources/store/manual.md

    images:
        - file: src/Resources/store/images/1.png
          # toggle visibility in that language
          activate:
            en: true
            de: true
          # preview image in that language
          preview:
            en: true
            de: true
          # sorting of the images
          priority: 1

changelog:
    # enable automatic changelog generation
    enabled: false

    # use openai to generate a better changelog based on commit messages. Requires OPENAI_TOKEN set
    ai_enabled: false

    # limit with regex which commits should be considered
    pattern: ''

    # allows to override the changelog generation template
    template: |
        {{range .Commits}}- [{{ .Message }}]({{ $.Config.VCSURL }}/{{ .Hash }})
        {{end}}

    # Allows to write RegEx groups into variables which can be used in the template
    variables:
        # extract the ticket number into variable.
        # can be then used in the template with {{ .Variables.ticket }}
        ticket: ^(NEXT-[0-9]+)
```
