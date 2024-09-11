package procfs

import (
	"gotest.tools/assert"
	"testing"
)

func TestProcessTime(t *testing.T) {
	time, err := getProcessTimeImpl("testdata/proc_pid_stat", 100)
	assert.NilError(t, err)
	assert.Equal(t, time, float32(0.02))
}
