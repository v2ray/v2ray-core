package git

import (
	"testing"

  v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestRevParse(t *testing.T) {
	v2testing.Current(t)

	rev, err := RevParse("HEAD")
	assert.Error(err).IsNil()
	assert.Int(len(rev)).GreaterThan(0)
}

func TestRepoVersion(t *testing.T) {
	v2testing.Current(t)

	version, err := RepoVersionHead()
	assert.Error(err).IsNil()
	assert.Int(len(version)).GreaterThan(0)
}
