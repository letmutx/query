package query

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const base = `SELECT image_id FROM images i`

func Test_Query(t *testing.T) {
	require := require.New(t)
	q, params, err := New(base).
		Where(`i.name LIKE {{ .Param "name" }}`).
		And(`i.size != {{ .Param "size" }}`).
		Build(Params{"name": "IMG_%", "size": 22})
	require.Nil(err)
	require.Equal(len(params), 2, "they should be equal")
	require.True(strings.HasPrefix(q, base))
}
