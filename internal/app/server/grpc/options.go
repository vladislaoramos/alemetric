package grpc

import "time"

type OptionFunc func(tool *GRPCServer)

// WriteFileWithDuration sets the writing to the file with duration.
func WriteFileWithDuration(duration time.Duration) OptionFunc {
	return func(gs *GRPCServer) {
		gs.writeToFileWithDuration = true
		gs.writeFileDuration = duration
	}
}

// SyncWriteFile sets the synchronous writing to the file.
func SyncWriteFile() OptionFunc {
	return func(gs *GRPCServer) {
		gs.syncWriteFile = true
	}
}

// CheckDataSign sets the data signing of all metrics for the tool.
func CheckDataSign(key string) OptionFunc {
	return func(gs *GRPCServer) {
		gs.encryptionKey = key
		gs.checkDataSign = true
	}
}
