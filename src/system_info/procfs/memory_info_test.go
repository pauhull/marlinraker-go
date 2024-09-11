package procfs

import (
	"gotest.tools/assert"
	"testing"
)

func TestMemAvail(t *testing.T) {
	avail, err := getMemAvailImpl("testdata/meminfo1")
	assert.NilError(t, err)
	assert.Equal(t, avail, int64(24552156))

	avail, err = getMemAvailImpl("testdata/meminfo2")
	assert.NilError(t, err)
	assert.Equal(t, avail, int64(302256))
}
