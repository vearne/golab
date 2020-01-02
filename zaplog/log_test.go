package zaplog

import (
	//"fmt"
	"testing"
)

func BenchmarkZapFile(b *testing.B) {
	logger := InitLogger()
	defer logger.Sync()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("hello world!")
	}
}

func BenchmarkZapStdout(b *testing.B) {
	logger := InitLogger2()
	defer logger.Sync()
	fmt.Println()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("hello world!")
	}
}

func BenchmarkAsync(b *testing.B) {
	logger := InitLogger3()
	defer logger.Sync()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("hello world!")
	}
}
