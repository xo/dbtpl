#!/bin/bash

set -ex

yaegi extract github.com/xo/dbtpl/loader
yaegi extract github.com/xo/dbtpl/types
yaegi extract os/exec
yaegi extract github.com/kenshaw/glob
yaegi extract github.com/goccy/go-yaml
yaegi extract github.com/kenshaw/inflector
yaegi extract github.com/kenshaw/snaker
yaegi extract github.com/Masterminds/sprig/v3
yaegi extract github.com/yookoala/realpath
yaegi extract golang.org/x/tools/imports
yaegi extract mvdan.cc/gofumpt/format
perl -pi -e 's/.*\n// if /Custom/' github_com-goccy-go-yaml.go
