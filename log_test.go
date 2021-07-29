package log

import "testing"

func Test_go_init(t *testing.T) {
	config := &Config{
		LogSaveName: "test.log",
		LogSavePath: "./logs",
	}
	log := Init(config)
	log.Info("Test log init success")
	log.Infof("%s 载人登月成功", "中国")
	log.Warn("Test log init success")
	log.Warnf("%s 载人登月成功", "中国")
	log.Debug("Test log init success")
	log.Debugf("%s 载人登月成功", "中国")
}
