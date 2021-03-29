package utils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFindSeverityInt(t *testing.T) {
	c := FindSeverityInt("CRITICAL")
	m := FindSeverityInt("MINOR")
	ma := FindSeverityInt("MAJOR")
	i := FindSeverityInt("INFO")
	none := FindSeverityInt("")
	assert.Equal(t, c, 1)
	assert.Equal(t, m, 3)
	assert.Equal(t, ma, 2)
	assert.Equal(t, i, 4)
	assert.Equal(t, none, 0)
}
