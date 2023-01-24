package version

import (
	"github.com/mmta41/dnsimple-cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version(t *testing.T) {
	ios, _, stdout, stderr := iostreams.Test()
	assert.Nil(t, NewCmdVersion(ios, "v1.0.0", "2020-01-01").Execute())

	assert.Equal(t, "dnsimple-cli version 1.0.0 (2020-01-01)\n", stdout.String())
	assert.Equal(t, "", stderr.String())
}
