module shopware-cli

go 1.17

require (
	github.com/caarlos0/env/v6 v6.9.1
	github.com/gorilla/schema v1.2.0
	github.com/hashicorp/go-version v1.4.0
	github.com/manifoldco/promptui v0.9.0
	github.com/microcosm-cc/bluemonday v1.0.17
	github.com/olekukonko/tablewriter v0.0.5
	github.com/otiai10/copy v1.7.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.3.0
	github.com/yuin/goldmark v1.3.5
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20211210111614-af8b64212486 // indirect
)

replace github.com/hashicorp/go-version => ./version
