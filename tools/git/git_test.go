package git

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestRevParse(t *testing.T) {
	assert := unit.Assert(t)

	rev, err := RevParse("HEAD")
	assert.Error(err).IsNil()
	assert.Int(len(rev)).GreaterThan(0)

	rev, err = RevParse("v0.8")
	assert.Error(err).IsNil()
	assert.String(rev).Equals("de7a1d30c3e6bda6a1297b5815369fcfa0e74f0e")
}

func TestRepoVersion(t *testing.T) {
	assert := unit.Assert(t)

	version, err := RepoVersionHead()
	assert.Error(err).IsNil()
	assert.Int(len(version)).GreaterThan(0)

	version, err = RepoVersion("de7a1d30c3e6bda6a1297b5815369fcfa0e74f0e")
	assert.Error(err).IsNil()
	assert.String(version).Equals("v0.8")
}
