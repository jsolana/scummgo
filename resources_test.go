package scummgo_test

import (
	"testing"

	"github.com/jsolana/scummgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceRefParse(t *testing.T) {
	ref, err := scummgo.ParseResourceRef("pkg:foo/bar")
	require.NoError(t, err)

	assert.Equal(t, scummgo.ResourcePackage("pkg"), ref.Package())
	assert.Equal(t, scummgo.ResourceID("foo/bar"), ref.ID())
}
