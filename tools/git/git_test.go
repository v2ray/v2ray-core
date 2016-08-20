package git

import (
	"testing"

	"v2ray.com/core/testing/assert"
)

func TestRevParse(t *testing.T) {
	assert := assert.On(t)

	rev, err := RevParse("HEAD")
	assert.Error(err).IsNil()
	assert.Int(len(rev)).GreaterThan(0)
}

func TestRepoVersion(t *testing.T) {
	assert := assert.On(t)

	version, err := RepoVersionHead()
	assert.Error(err).IsNil()
	assert.Int(len(version)).GreaterThan(0)
}
