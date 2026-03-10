package memory

import (
	"testing"

	"github.com/AntonPaus/GolangAdvanced/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Get(t *testing.T) {
	storage := NewMemoryStorage()
	t.Run("Gauge", func(t *testing.T) {
		err := storage.Set(repository.MetricTypeGauge, "g1", 3.1)
		require.NoError(t, err)
		got1, err := storage.Get(repository.MetricTypeGauge, "g1")
		require.NoError(t, err)
		require.Equal(t, 3.1, got1)
	})
	t.Run("Counter", func(t *testing.T) {
		err := storage.Set(repository.MetricTypeCounter, "c1", int64(3))
		require.NoError(t, err)
		err = storage.Set(repository.MetricTypeCounter, "c1", int64(3))
		require.NoError(t, err)
		got1, err := storage.Get(repository.MetricTypeCounter, "c1")
		require.NoError(t, err)
		require.Equal(t, int64(6), got1)
	})
	t.Run("Wrong type", func(t *testing.T) {
		err := storage.Set("wrongType", "c1", int64(3))
		require.Error(t, err)
	})
}
