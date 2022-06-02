package sallust

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// exactLevelEnabler is used to allow only (1) particular level of logging.
// This can be used to assert that the appropriate log method got called.
type exactLevelEnabler struct {
	Level zapcore.Level
}

func (ele exactLevelEnabler) Enabled(l zapcore.Level) bool {
	return ele.Level == l
}

func testPrinterPrintf(t *testing.T) {
	testData := []struct {
		printerLevel zapcore.Level
		levelEnabler zapcore.LevelEnabler
	}{
		{
			printerLevel: zapcore.DebugLevel,
			levelEnabler: exactLevelEnabler{Level: zapcore.DebugLevel},
		},
		{
			printerLevel: zapcore.InfoLevel,
			levelEnabler: exactLevelEnabler{Level: zapcore.InfoLevel},
		},
		{
			printerLevel: zapcore.WarnLevel,
			levelEnabler: exactLevelEnabler{Level: zapcore.WarnLevel},
		},
		{
			printerLevel: zapcore.ErrorLevel,
			levelEnabler: exactLevelEnabler{Level: zapcore.ErrorLevel},
		},
		{
			printerLevel: zapcore.Level(123),
			levelEnabler: exactLevelEnabler{Level: zapcore.InfoLevel},
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert = assert.New(t)
				buffer Buffer
				core   = zapcore.NewCore(
					zapcore.NewJSONEncoder(zapcore.EncoderConfig{
						MessageKey:  "msg",
						LevelKey:    "level",
						EncodeLevel: zapcore.LowercaseLevelEncoder,
					}),
					&buffer,
					record.levelEnabler,
				)

				logger  = zap.New(core)
				printer = Printer{SugaredLogger: logger.Sugar(), Level: record.printerLevel}
			)

			// just check to see if the output passed the levelEnabler
			printer.Printf("test: %s %d", "string", 123)
			logger.Sync()
			assert.Greater(buffer.Len(), 0)
		})
	}
}

func testPrinterPrint(t *testing.T) {
	testData := []struct {
		printerLevel zapcore.Level
		levelEnabler zapcore.LevelEnabler
	}{
		{
			printerLevel: zapcore.DebugLevel,
			levelEnabler: exactLevelEnabler{Level: zapcore.DebugLevel},
		},
		{
			printerLevel: zapcore.InfoLevel,
			levelEnabler: exactLevelEnabler{Level: zapcore.InfoLevel},
		},
		{
			printerLevel: zapcore.WarnLevel,
			levelEnabler: exactLevelEnabler{Level: zapcore.WarnLevel},
		},
		{
			printerLevel: zapcore.ErrorLevel,
			levelEnabler: exactLevelEnabler{Level: zapcore.ErrorLevel},
		},
		{
			printerLevel: zapcore.Level(123),
			levelEnabler: exactLevelEnabler{Level: zapcore.InfoLevel},
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert = assert.New(t)
				buffer Buffer
				core   = zapcore.NewCore(
					zapcore.NewJSONEncoder(zapcore.EncoderConfig{
						MessageKey:  "msg",
						LevelKey:    "level",
						EncodeLevel: zapcore.LowercaseLevelEncoder,
					}),
					&buffer,
					record.levelEnabler,
				)

				logger  = zap.New(core)
				printer = Printer{SugaredLogger: logger.Sugar(), Level: record.printerLevel}
			)

			// just check to see if the output passed the levelEnabler
			printer.Print("string", 123, -67, 4)
			logger.Sync()
			assert.Greater(buffer.Len(), 0)
		})
	}
}

func TestPrinter(t *testing.T) {
	t.Run("Printf", testPrinterPrintf)
	t.Run("Print", testPrinterPrint)
}
