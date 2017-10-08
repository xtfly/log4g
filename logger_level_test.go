package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {
	assert.Equal(t, Error, LevelFrom("error"))
	assert.Equal(t, Error, LevelFrom("Error"))
	assert.Equal(t, Error, LevelFrom("eRror"))
	assert.Equal(t, Uninitialized, LevelFrom("xx"))

	assert.Equal(t, "ERROR", Error.String())
	assert.Equal(t, "ERR", Error.ShortStr())
}
