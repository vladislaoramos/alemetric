package repo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreFilePath(t *testing.T) {
	repo := &MetricsRepo{}
	path := "./path"
	op := StoreFilePath(path)
	op(repo)
	require.Equal(t, path, repo.StoreFilePath)
}

func TestRestore(t *testing.T) {
	repo := &MetricsRepo{}
	op := Restore()
	op(repo)
	require.True(t, repo.Restore)
}
