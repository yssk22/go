package xlog

import "testing"

func benchLogText(b *testing.B, logger *Logger) {
	for i := 0; i < b.N; i++ {
		logger.Infof("[%d] This is bench", i)
	}
}

func BenchmarkLoggerStackCapture(b *testing.B) {
	logger := New(NullSink)
	logger.MinStackCaptureOn = LevelTrace
	benchLogText(b, logger)
}

func BenchmarkLoggerNoStackCapture(b *testing.B) {
	logger := New(NullSink)
	logger.MinStackCaptureOn = LevelFatal
	benchLogText(b, logger)
}

func BenchmarkLoggerFewStackCapture(b *testing.B) {
	logger := New(NullSink)
	logger.MinStackCaptureOn = LevelTrace
	logger.StackCaptureDepth = 1
	benchLogText(b, logger)
}
