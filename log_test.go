package log

import "testing"

func Test_go_init(t *testing.T) {
	config := Config{
		LogSaveName: "test.log",
		LogSavePath: "./",
	}
	logClient := Init(config)
	logClient.Info("Test go init success")
}
