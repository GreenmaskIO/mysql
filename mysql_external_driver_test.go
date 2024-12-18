package mysql

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternalDriver_EncodeValueByTypeName(t *testing.T) {
	d := NewExternalDriver()

	t.Run("text", func(t *testing.T) {
		res, err := d.EncodeValueByTypeName("name", "src", []byte{})
		require.NoError(t, err)
		require.Equal(t, []byte("src"), res)
	})

	t.Run("blob", func(t *testing.T) {
		res, err := d.EncodeValueByTypeName("name", []byte("src"), []byte{})
		require.NoError(t, err)
		require.Equal(t, []byte("src"), res)
	})

	t.Run("raw data", func(t *testing.T) {
		res, err := d.EncodeValueByTypeName("name", sql.RawBytes("test"), []byte{})
		require.NoError(t, err)
		require.Equal(t, []byte("test"), res)
	})

	t.Run("int", func(t *testing.T) {
		res, err := d.EncodeValueByTypeName("name", int32(1), []byte{})
		require.NoError(t, err)
		require.Equal(t, []byte("1"), res)
	})

}
