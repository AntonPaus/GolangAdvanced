package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func UncompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		decompressedBody, err := io.ReadAll(gz)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(decompressedBody))
		r.ContentLength = int64(len(decompressedBody))
		r.Header.Del("Content-Encoding")
		next.ServeHTTP(w, r)
	})
}

func CompressGzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gzip writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		w.Close()
		return nil, fmt.Errorf("failed to write data to gzip writer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to finalize gzip compression: %v", err)
	}
	return b.Bytes(), nil
}

// func CompressHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
// 			var b bytes.Buffer
// 			gz, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			body, err := io.ReadAll(r.Body)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusBadRequest)
// 				return
// 			}
// 			_, err = gz.Write(body)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			err = gz.Close()
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			w.Header().Set("Content-Encoding", "gzip")
// 			w.Write(b.Bytes())
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 	})
// }

// func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ow := w

// 		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
// 		acceptEncoding := r.Header.Get("Accept-Encoding")
// 		supportsGzip := strings.Contains(acceptEncoding, "gzip")
// 		if supportsGzip {
// 			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
// 			cw := newCompressWriter(w)
// 			// меняем оригинальный http.ResponseWriter на новый
// 			ow = cw
// 			// не забываем отправить клиенту все сжатые данные после завершения middleware
// 			defer cw.Close()
// 		}

// 		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
// 		contentEncoding := r.Header.Get("Content-Encoding")
// 		sendsGzip := strings.Contains(contentEncoding, "gzip")
// 		if sendsGzip {
// 			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
// 			cr, err := newCompressReader(r.Body)
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 			// меняем тело запроса на новое
// 			r.Body = cr
// 			defer cr.Close()
// 		}

// 		// передаём управление хендлеру
// 		h.ServeHTTP(ow, r)
// 	}
// }

// // ...

// func run() error {
// 	if err := logger.Initialize(flagLogLevel); err != nil {
// 		return err
// 	}

// 	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
// 	// оборачиваем хендлер webhook в middleware с логированием и поддержкой gzip
// 	return http.ListenAndServe(flagRunAddr, logger.RequestLogger(gzipMiddleware(webhook)))
// }

// func Compress(w http.ResponseWriter, r *http.Request) {

// 	gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer gz.Close()
// 	w.Header().Set("Content-Encoding", "gzip")
// 	io.WriteString(gz, "Hello, World!")

// }

// func Decompress(w http.ResponseWriter, r *http.Request) {
// 	gz, err := gzip.NewReader(r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer gz.Close()
// 	io.Copy(w, gz)
// }
