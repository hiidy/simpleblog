package log

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debugw(msg string, kvs ...any)

	Infow(msg string, kvs ...any)

	Warnw(msg string, kvs ...any)

	Errorw(msg string, kvs ...any)

	Panicw(msg string, kvs ...any)

	Fatalw(msg string, kvs ...any)

	Sync()
}

type zapLogger struct {
	z *zap.Logger
}

var (
	mu sync.Mutex

	std = New(NewOptions())
)

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	std = New(opts)
}

func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level

	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	// zap.Logger 구축에 필요한 구성 생성
	cfg := &zap.Config{
		// 로그에 로그를 호출한 파일과 줄 번호 표시 여부, 예: `"caller":"apiserver/server.go:75"`
		DisableCaller: opts.DisableCaller,
		// panic 이상 레벨에서 스택 정보 출력 금지 여부
		DisableStacktrace: opts.DisableStacktrace,
		// 로그 레벨 지정
		Level: zap.NewAtomicLevelAt(zapLevel),
		// 로그 표시 형식 지정, 선택 가능한 값: console, json
		Encoding:      opts.Format,
		EncoderConfig: encoderConfig,
		// 로그 출력 위치 지정
		OutputPaths: opts.OutputPaths,
		// zap 내부 에러 출력 위치 설정
		ErrorOutputPaths: []string{"stderr"},
	}

	// cfg를 사용하여 *zap.Logger 객체 생성
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	// 표준 라이브러리의 log 출력을 zap.Logger로 리디렉션
	zap.RedirectStdLog(z)

	return &zapLogger{z: z}
}

func Sync() {
	std.Sync()
}

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

func Debugw(msg string, kvs ...any) {
	std.Debugw(msg, kvs...)
}

func (l *zapLogger) Debugw(msg string, kvs ...any) {
	l.z.Sugar().Debugw(msg, kvs...)
}

func Infow(msg string, kvs ...any) {
	std.Infow(msg, kvs...)
}

func (l *zapLogger) Infow(msg string, kvs ...any) {
	l.z.Sugar().Infow(msg, kvs...)
}

func Warnw(msg string, kvs ...any) {
	std.Warnw(msg, kvs...)
}

func (l *zapLogger) Warnw(msg string, kvs ...any) {
	l.z.Sugar().Warnw(msg, kvs...)
}

func Errorw(msg string, kvs ...any) {
	std.Errorw(msg, kvs...)
}

func (l *zapLogger) Errorw(msg string, kvs ...any) {
	l.z.Sugar().Errorw(msg, kvs...)
}

func Panicw(msg string, kvs ...any) {
	std.Panicw(msg, kvs...)
}

func (l *zapLogger) Panicw(msg string, kvs ...any) {
	l.z.Sugar().Panicw(msg, kvs...)
}

func Fatalw(msg string, kvs ...any) {
	std.Fatalw(msg, kvs...)
}

func (l *zapLogger) Fatalw(msg string, kvs ...any) {
	l.z.Sugar().Fatalw(msg, kvs...)
}
