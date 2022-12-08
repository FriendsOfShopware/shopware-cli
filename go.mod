module github.com/FriendsOfShopware/shopware-cli

go 1.19

require (
	github.com/bep/godartsass v0.16.0
	github.com/caarlos0/env/v6 v6.10.1
	github.com/doutorfinancas/go-mad v0.0.0-20221115152854-f38f7c284800
	github.com/evanw/esbuild v0.16.2
	github.com/friendsofshopware/go-shopware-admin-api-sdk v0.0.0-20220325180335-81b5b971debc
	github.com/google/uuid v1.3.0
	github.com/gorilla/schema v1.2.0
	github.com/manifoldco/promptui v0.9.0
	github.com/mholt/archiver/v3 v3.5.1
	github.com/microcosm-cc/bluemonday v1.0.21
	github.com/olekukonko/tablewriter v0.0.5
	github.com/otiai10/copy v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/schollz/progressbar/v3 v3.12.2
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.6.1
	github.com/vulcand/oxy v1.4.2
	github.com/yuin/goldmark v1.5.3
	go.uber.org/zap v1.24.0
	gopkg.in/antage/eventsource.v1 v1.0.0-20150318155416-803f4c5af225
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/cli/safeexec v1.0.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jaswdr/faker v1.15.0 // indirect
	github.com/klauspost/compress v1.15.12 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/oauth2 v0.1.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

// remove when https://github.com/doutorfinancas/go-mad/pull/42 is merged
replace github.com/doutorfinancas/go-mad v0.0.0-20221031120329-288cd003774876db34e5798f39b382ffa53d204c => github.com/shyim/go-mad v0.0.0-20221031184014-288cd0037748
