package logger

import (
	"net/http"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}

// type Logger struct {
// 	Sugar *zap.SugaredLogger
// }

// func NewLogger() (*Logger, error) {
// 	logger, err := zap.NewDevelopment()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer logger.Sync()
// 	s := logger.Sugar()
// 	sugar := &Logger{
// 		Sugar: s,
// 	}
// 	return sugar, nil
// }

// func (l *Logger) WithLogging(h http.Handler) http.Handler {
// 	logFn := func(w http.ResponseWriter, r *http.Request) {
// 		// функция Now() возвращает текущее время
// 		start := time.Now()

// 		// эндпоинт /ping
// 		uri := r.RequestURI
// 		// метод запроса
// 		method := r.Method

// 		// точка, где выполняется хендлер pingHandler
// 		h.ServeHTTP(w, r) // обслуживание оригинального запроса

// 		// Since возвращает разницу во времени между start
// 		// и моментом вызова Since. Таким образом можно посчитать
// 		// время выполнения запроса.
// 		duration := time.Since(start)

// 		// отправляем сведения о запросе в zap
// 		sugar.Infoln(
// 			"uri", uri,
// 			"method", method,
// 			"duration", duration,
// 		)

// 	}
// 	// возвращаем функционально расширенный хендлер
// 	return http.HandlerFunc(logFn)
// }
