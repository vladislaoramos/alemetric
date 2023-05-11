package usecase

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWriteFileWithDuration(t *testing.T) {
	mt := &ToolUseCase{}
	dur := time.Second
	op := WriteFileWithDuration(dur)
	op(mt)
	require.Equal(t, true, mt.writeToFileWithDuration)
	require.Equal(t, time.Second, mt.writeFileDuration)
}

func TestSyncWriteFile(t *testing.T) {
	mt := &ToolUseCase{}
	op := SyncWriteFile()
	op(mt)
	require.True(t, mt.syncWriteFile)
}

func TestCheckDataSign(t *testing.T) {
	mt := &ToolUseCase{}
	key := "string"
	op := CheckDataSign(key)
	op(mt)
	require.Equal(t, "string", mt.encryptionKey)
	require.True(t, mt.checkDataSign)
}
