// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

module github.com/staringfun/millsmess/services/socket

replace (
	github.com/staringfun/millsmess/libs/base => ../../libs/base
	github.com/staringfun/millsmess/libs/internal-core-api => ../../libs/internal-core-api
	github.com/staringfun/millsmess/libs/test-utils => ../../libs/test-utils
)

go 1.24.4

require github.com/stretchr/testify v1.10.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/sync v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
