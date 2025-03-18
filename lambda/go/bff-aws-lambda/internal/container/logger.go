package container

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/dcastellini/bff-lambda-service/internal/config"
	"io"
	"os"
	"strings"

	"context"
	"github.com/goccy/go-json"
	reflect "github.com/goccy/go-reflect"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	ZapLogger
}

type StdLogger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type ZapFields interface {
	AnyField(key string, value interface{}) zap.Field
	StringField(key, value string) zap.Field
	Any(val interface{}) zap.Field
	Struct(val interface{}) zap.Field
}

// ZapLogger GetError  interface.
type ZapLogger interface {
	StdLogger
	ZapFields
	Debug(msg string, keyvals ...zap.Field)
	Error(msg string, keyvals ...zap.Field)
	Info(msg string, keyvals ...zap.Field)
	Sync() error
	Warn(msg string, keyvals ...zap.Field)
	With(keyvals ...zap.Field) ZapLogger
}

// anyFieldWithCustomKey new key val field , value any type.
func anyFieldWithCustomKey(key string, value interface{}) zap.Field {
	defaultPrefix := "extra."

	return zap.Any(defaultPrefix+key, value)
}

type logger struct {
	log *zap.Logger
}

func (l *logger) Struct(val interface{}) zap.Field {
	return structField(val)
}

func (l *logger) Any(value interface{}) zap.Field {
	return anyField(value)
}

func (l *logger) StringField(key, value string) zap.Field {
	return zap.String(key, value)
}

func (l *logger) AnyField(key string, value interface{}) zap.Field {
	return anyFieldWithCustomKey(key, value)
}

func (l *logger) Debug(msg string, keyvals ...zap.Field) {
	l.log.Debug(msg, keyvals...)
}

func (l *logger) Error(msg string, keyvals ...zap.Field) {
	l.log.Error(msg, keyvals...)
}

func (l *logger) Info(msg string, keyvals ...zap.Field) {
	l.log.Info(msg, keyvals...)
}

func (l *logger) Warn(msg string, keyvals ...zap.Field) {
	l.log.Warn(msg, keyvals...)
}

func (l *logger) With(keyvals ...zap.Field) ZapLogger {
	return &logger{
		log: l.log.With(keyvals...),
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.log.Error(fmt.Sprint(format, v))
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.log.Warn(fmt.Sprint(format, v))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.log.Debug(fmt.Sprint(format, v))
}

func (l *logger) Sync() error {
	return l.log.Sync() //nolint: wrapcheck
}

// NewNopLambdaLogger returns a new LambdaLogger with a nop writer (no operations).
func NewNopLambdaLogger() ZapLogger {
	return &logger{log: zap.NewNop()}
}

// NewDefaultLogger configures logger.
func NewDefaultLogger() ZapLogger {
	return &logger{log: zap.NewExample()}
}

// NewLoggerWithBuildInfo configures logger.
func NewLoggerWithBuildInfo(cfg *config.LoggerConfiguration, buildInfo *BuildInfo) ZapLogger {
	return &logger{log: newZapLogger(buildInfo, cfg)}
}

func newZapLogger(infoBuild *BuildInfo, cfg *config.LoggerConfiguration, options ...zap.Option) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:          "data",
		LevelKey:            "logtype",
		TimeKey:             "date",
		NameKey:             "logger",
		CallerKey:           "custom.caller",
		FunctionKey:         "custom.function",
		StacktraceKey:       "exception.stackTrace",
		SkipLineEnding:      false,
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.RFC3339TimeEncoder,
		EncodeDuration:      zapcore.StringDurationEncoder,
		NewReflectedEncoder: defaultReflectedEncoder,
	}

	encoder := &customEncoder{Encoder: zapcore.NewJSONEncoder(encoderCfg)}

	addDefaultFields(encoder, infoBuild, cfg)

	zapLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	if strings.ToLower(cfg.General.LogLevel) == "debug" {
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	core := zapcore.NewCore(encoder, os.Stderr, zapLevel)

	options = append(options, zap.AddCallerSkip(1), zap.WithCaller(false))

	return zap.New(core).WithOptions(options...)
}

func addDefaultFields(encoder zapcore.Encoder, infoBuild *BuildInfo, cfg *config.LoggerConfiguration) {
	encoder.AddString("build.version", infoBuild.Version)
	encoder.AddString("build.version", infoBuild.Version)
	encoder.AddString("build.date", infoBuild.Date)
	encoder.AddString("build.hash", infoBuild.Hash)
	encoder.AddString("build.runtime", infoBuild.Runtime)
	encoder.AddString("function.country", cfg.General.Country)
	encoder.AddString("function.env", cfg.General.Environment)
	encoder.AddString("function.flow", cfg.General.Flow)
	encoder.AddString("function.version", cfg.General.Version)
	encoder.AddString("function.name", cfg.General.Name)
}

func defaultReflectedEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := json.NewEncoder(w)

	enc.SetEscapeHTML(false)

	return enc
}

type customEncoder struct {
	zapcore.Encoder
}

func (e *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	return e.Encoder.EncodeEntry(entry, fields) //nolint: wrapcheck
}

// structField Log structField Field.
func structField(val interface{}) zap.Field {
	return customField(val)
}

func anyField(val interface{}) zap.Field {
	return customField(val)
}

// customField log any Sensitive Data.
func customField(val interface{}) zap.Field {
	reflectType := reflect.TypeOf(val)

	var key string

	if reflectType.Kind() == reflect.Ptr {
		key = "extra.*" + reflectType.Elem().Name()
	} else {
		key = "extra." + reflectType.Name()
	}

	reflectTypeValue := reflect.ValueOf(val)

	switch reflectTypeValue.Kind() { //nolint: exhaustive, nolintlint
	case reflect.Struct:
		return zap.Any(key, maskSensitiveFields(reflectType, reflectTypeValue))
	default:
		return zap.Any(key, val)
	}
}

func maskSensitiveFields(reflectType reflect.Type, reflectTypeValue reflect.Value) map[string]interface{} {
	out := map[string]interface{}{}

	for index := 0; index < reflectType.NumField(); index++ {
		reflectTypeField := reflectType.Field(index)

		jsonTag := reflectTypeField.Tag.Get("json")

		if reflectTypeField.Tag.Get("zap") == "sensitive" {
			out[jsonTag] = "[filtered]"
		} else {
			if jsonTag != "" {
				out[strings.Split(jsonTag, ",")[0]] = reflectTypeValue.Field(index).Interface()
			} else {
				out[reflectTypeField.Name] = reflectTypeValue.Field(index).Interface()
			}
		}
	}

	return out
}

// BuildInfo Build Information.
type BuildInfo struct {
	Version string
	Date    string
	Hash    string
	Runtime string
}

// NewBuildInfo create new build info.
func NewBuildInfo(version, date, hash, runtime string) *BuildInfo {
	return &BuildInfo{Version: version, Date: date, Hash: hash, Runtime: runtime}
}

type loggerKey struct{}

type ContextLogger struct {
	logger    ZapLogger
	lambdaArn string
}

func NewContextLogger(logger ZapLogger) *ContextLogger {
	return &ContextLogger{logger: logger}
}

type Middleware interface {
	WithLogger(parentCtx context.Context, logger ZapLogger) context.Context
}

func (l *ContextLogger) WithLogger(parentCtx context.Context, keyvals ...zap.Field) (context.Context, ZapLogger) {
	lambdaContext, found := lambdacontext.FromContext(parentCtx)
	if found {
		if l.lambdaArn == "" {
			l.lambdaArn = lambdaContext.InvokedFunctionArn
		}
		keyvals = append(keyvals,
			zap.String("request.awsRequestId", lambdaContext.AwsRequestID),
			zap.String("function.arn", l.lambdaArn),
		)
	}

	lambdaLogger := l.logger.With(keyvals...)

	return context.WithValue(parentCtx, loggerKey{}, lambdaLogger), lambdaLogger
}

func GetLogger(ctx context.Context) ZapLogger { //nolint: ireturn
	lambdaLogger, found := ctx.Value(loggerKey{}).(ZapLogger)
	if !found {
		return NewNopLambdaLogger()
	}

	return lambdaLogger
}
