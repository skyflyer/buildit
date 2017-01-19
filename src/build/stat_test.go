package build

import (
	"conf"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveReadLastHead(t *testing.T) {
	SaveLastHead("1234")
	lh := GetLastHead()
	require.Equal(t, "1234", lh)

	SaveLastHead("5678")
	lh = GetLastHead()
	require.Equal(t, "5678", lh)
}

func TestCleanup(t *testing.T) {
	os.Remove(conf.LastHeadMarker)
}
