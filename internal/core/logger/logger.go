package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	// не даем название этому полю, чтобы встроить его в структуру логгера. Таким образом, мы можем пользоваться функциональностью заплоггера напрямую через свою структуру.
	*zap.Logger

	// файл, в который будут записываться логи. Создаем указать, чтобы его закрыть при окончании работы 
	file *os.File
}

// Функция для получения логгера из контекста. Необходима отдельная функция, т.к. делать это приходится довольно часто. Первое использование в common.go
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		// с помощбю %w оборачиваем ошибку добавляя контекст. Благодаря этому получатель этой ошибки сможет получить исходную ошибку без контекста с помощью errors.Is()
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	// 0755 - выдача прав
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	// лог файл будет называться согласна таймстемпу создания файла, поэтому получаем его
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
		)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
		)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file: logFile,
	}, nil
}

// переопределили стандартный метод With у встроенного зап логгера. Это нужно, потому что стандартный 
// метод возвращает указатель на zap.Logger, а нам нужно, чтобы он возвращал указатель на наш созданный логер
// теперь он это делает.
// данная конструкция необходима для middleware Logger в /internal/core/transport/http/middleware/middleware.go
// подробнее можно узнать в комментариях там
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file: l.file,
	}
}


func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}
