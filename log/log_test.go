package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Logger()
	}, "")
}

func TestLogInfo(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Info("test111111")
	}, "")
}

func TestLogInfof(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Infof("test111111: %s", "pppp")
	}, "")
}

func TestLogDebug(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Debug("test111111")
	}, "")
}

func TestLogDebugf(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Debugf("test111111: %s", "pppp")
	}, "")
}

func TestLogWarn(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Warn("test111111")
	}, "")
}

func TestLogWarnf(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Warnf("test111111: %s", "pppp")
	}, "")
}

func TestLogError(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Error("test111111")
	}, "")
}

func TestLogErrorf(t *testing.T) {
	assert.NotPanicsf(t, func() {
		Errorf("test111111: %s", "pppp")
	}, "")
}

func TestLogPanic(t *testing.T) {
	assert.Panicsf(t, func() {
		Panic("test111111")
	}, "")
}

func TestLogPanicf(t *testing.T) {
	assert.Panicsf(t, func() {
		Panicf("test111111: %s", "pppp")
	}, "")
}
