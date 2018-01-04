package portfwd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomHostList(t *testing.T) {
	testHosts := []string{
		"localhost:1234",
		"localhost:5678",
		"localhost:1234",
	}
	assert := assert.New(t)
	l := NewRandomHostList()
	err := l.AddHost(testHosts[0])
	assert.NoError(err)
	err = l.AddHost(testHosts[1])
	assert.NoError(err)
	err = l.AddHost(testHosts[2])
	assert.Error(err)

	host, err := l.Host()
	assert.NoError(err)
	assert.NotEmpty(host)
}
