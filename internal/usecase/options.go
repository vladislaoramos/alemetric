package usecase

import "time"

type OptionFunc func(tool *ToolUseCase)

// WriteFileWithDuration sets the writing to the file with duration.
func WriteFileWithDuration(duration time.Duration) OptionFunc {
	return func(mt *ToolUseCase) {
		mt.writeToFileWithDuration = true
		mt.writeFileDuration = duration
	}
}

// SyncWriteFile sets the synchronous writing to the file.
func SyncWriteFile() OptionFunc {
	return func(mt *ToolUseCase) {
		mt.syncWriteFile = true
	}
}

// CheckDataSign sets the data signing of all metrics for the tool.
func CheckDataSign(key string) OptionFunc {
	return func(mt *ToolUseCase) {
		mt.encryptionKey = key
		mt.checkDataSign = true
	}
}
