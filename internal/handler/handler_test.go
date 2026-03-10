package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AntonPaus/GolangAdvanced/internal/repository"
	"github.com/AntonPaus/GolangAdvanced/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetric(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "success gauge",
			method:     http.MethodPost,
			path:       "/update/gauge/g1/3.1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "success counter",
			method:     http.MethodPost,
			path:       "/update/counter/c1/7",
			wantStatus: http.StatusOK,
		},
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			path:       "/update/gauge/g1/3.1",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "wrong path",
			method:     http.MethodPost,
			path:       "/update/gauge/g1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid metric type",
			method:     http.MethodPost,
			path:       "/update/unknown/m1/1",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid gauge value",
			method:     http.MethodPost,
			path:       "/update/gauge/g1/not-a-float",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid counter value",
			method:     http.MethodPost,
			path:       "/update/counter/c1/not-an-int",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Storage: memory.NewMemoryStorage(),
			}
			request := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			h.UpdateMetric(w, request)

			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, result.StatusCode)
		})
	}
}

func TestUpdateMetric_CounterAccumulates(t *testing.T) {
	storage := memory.NewMemoryStorage()
	h := &Handler{Storage: storage}

	firstReq := httptest.NewRequest(http.MethodPost, "/update/counter/c1/5", nil)
	firstW := httptest.NewRecorder()
	h.UpdateMetric(firstW, firstReq)
	firstResult := firstW.Result()
	err := firstResult.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, firstResult.StatusCode)

	secondReq := httptest.NewRequest(http.MethodPost, "/update/counter/c1/2", nil)
	secondW := httptest.NewRecorder()
	h.UpdateMetric(secondW, secondReq)
	secondResult := secondW.Result()
	err = secondResult.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, secondResult.StatusCode)

	value, err := storage.Get(repository.MetricTypeCounter, "c1")
	require.NoError(t, err)
	assert.Equal(t, int64(7), value)
}
