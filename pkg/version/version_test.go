package version

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, Version, "1.2.0")
}
