package file

import (
	"testing"

	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Get(t *testing.T) {
	storage, err := NewFileStorage(true, "test.json", 1)
	require.NoError(t, err)
	t.Run("Gauge", func(t *testing.T) {
		err := storage.Set(interfaces.MetricTypeGauge, "g1", 3.1)
		require.NoError(t, err)
		got1, err := storage.Get(interfaces.MetricTypeGauge, "g1")
		require.NoError(t, err)
		require.Equal(t, 3.1, got1)
	})
	t.Run("Counter", func(t *testing.T) {
		err := storage.Set(interfaces.MetricTypeCounter, "c1", int64(3))
		require.NoError(t, err)
		err = storage.Set(interfaces.MetricTypeCounter, "c1", int64(3))
		require.NoError(t, err)
		got1, err := storage.Get(interfaces.MetricTypeCounter, "c1")
		require.NoError(t, err)
		require.Equal(t, int64(6), got1)
	})
	t.Run("Wrong type", func(t *testing.T) {
		err := storage.Set("wrongType", "c1", int64(3))
		require.Error(t, err)
	})
}
