package scummgo_test

import (
	"bytes"
	"testing"

	"github.com/jsolana/scummgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDumpPackage_ScriptResource(t *testing.T) {
	var idxBuf, datBuf bytes.Buffer
	enc, err := scummgo.NewResourceEncoder(&idxBuf, &datBuf)
	require.NoError(t, err)

	err = enc.EncodeScript(
		scummgo.ResourceID("hello"),
		&scummgo.Script{Language: scummgo.ScriptLua, Code: []byte("print('hi')")},
		scummgo.CompressionNone,
	)
	require.NoError(t, err)

	dump, err := scummgo.DumpPackageBytes(idxBuf.Bytes(), datBuf.Bytes())
	require.NoError(t, err)
	require.Len(t, dump.Resources, 1)

	r := dump.Resources[0]
	assert.Equal(t, "hello", r.ID)
	assert.Equal(t, "script", r.Type)
	assert.Equal(t, uint8(4), r.TypeCode)
	assert.Equal(t, "none", r.Compression)
	assert.NotEmpty(t, r.PayloadSHA256)
	assert.NotEmpty(t, r.DecodedSHA256)

	sd, ok := r.Detail.(scummgo.ScriptDumpDetail)
	require.True(t, ok)
	assert.Equal(t, "lua", sd.Language)
	assert.Equal(t, 11, sd.CodeSize)
	assert.Equal(t, "print('hi')", sd.CodePreview)
}

func TestDumpPackageJSON_Stable(t *testing.T) {
	var idxBuf, datBuf bytes.Buffer
	enc, err := scummgo.NewResourceEncoder(&idxBuf, &datBuf)
	require.NoError(t, err)
	require.NoError(t, enc.EncodeScript(
		scummgo.ResourceID("a"),
		&scummgo.Script{Language: scummgo.ScriptLua, Code: []byte("x")},
		scummgo.CompressionNone,
	))

	dump, err := scummgo.DumpPackageBytes(idxBuf.Bytes(), datBuf.Bytes())
	require.NoError(t, err)
	j1, err := scummgo.DumpPackageJSON(dump)
	require.NoError(t, err)
	j2, err := scummgo.DumpPackageJSON(dump)
	require.NoError(t, err)
	assert.Equal(t, j1, j2)
}
