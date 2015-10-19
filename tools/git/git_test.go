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
}

func TestRepoVersion(t *testing.T) {
	assert := unit.Assert(t)

	version, err := RepoVersionHead()
	assert.Error(err).IsNil()
	assert.Int(len(version)).GreaterThan(0)
}
