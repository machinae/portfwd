package main

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testConfigFile = `
# Keys can be named anything
[forward1]
from = ":8000"
to = ["rhost1:8000", "rhost2:8000"]
strategy = "random"

[forward2]
from = ":9000"
to = "rhost1:9000"
`

func TestConfig(t *testing.T) {
	assert := assert.New(t)
	r := strings.NewReader(testConfigFile)
	viper.SetConfigType("toml")
	err := viper.ReadConfig(r)
	assert.NoError(err)
	err = parseConfig()
	assert.NoError(err)
}
