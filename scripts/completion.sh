#!/usr/bin/env bash

rm -rf completions
mkdir completions
go run . completion bash > completions/shopware-cli.bash
go run . completion zsh > completions/shopware-cli.zsh
go run . completion fish > completions/shopware-cli.fish