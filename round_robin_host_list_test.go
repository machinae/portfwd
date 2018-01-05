package portfwd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobinHostList(t *testing.T) {
	testHosts := []string{
		"localhost:1234",
		"localhost:5678",
	}
	assert := assert.New(t)
	l := NewRoundRobinHostList()
	err := l.AddHost(testHosts[0])
	assert.NoError(err)
	err = l.AddHost(testHosts[1])
	assert.NoError(err)

	host, err := l.Host()
	assert.NoError(err)
	assert.Equal(testHosts[0], host)

	host, err = l.Host()
	assert.NoError(err)
	assert.Equal(testHosts[1], host)

	host, err = l.Host()
	assert.NoError(err)
	assert.Equal(testHosts[0], host)

	host, err = l.Host()
	assert.NoError(err)
	assert.Equal(testHosts[1], host)

}
