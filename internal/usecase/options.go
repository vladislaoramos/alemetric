package usecase

import "time"

type OptionFunc func(tool *ToolUseCase)

func WriteFileWithDuration(duration time.Duration) OptionFunc {
	return func(mt *ToolUseCase) {
		mt.writeToFileWithDuration = true
		mt.writeFileDuration = duration
	}
}

func SyncWriteFile() OptionFunc {
	return func(mt *ToolUseCase) {
		mt.syncWriteFile = true
	}
}
