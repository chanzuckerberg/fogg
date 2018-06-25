package apply

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveExtension(t *testing.T) {
	var x string

	x = removeExtension("foo")
	assert.Equal(t, "foo", x)

	x = removeExtension("foo.txt")
	assert.Equal(t, "foo", x)

	x = removeExtension("foo.txt.asdf")
	assert.Equal(t, "foo.txt", x)
}
