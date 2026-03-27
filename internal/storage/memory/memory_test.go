package memory

import (
	"context"
	"testing"

	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_Get(t *testing.T) {
	storage, err := NewMemoryStorage()
	require.NoError(t, err)
	t.Run("Gauge", func(t *testing.T) {
		err := storage.Set(context.TODO(), interfaces.MetricTypeGauge, "g1", 3.1)
		require.NoError(t, err)
		got1, err := storage.Get(context.TODO(), interfaces.MetricTypeGauge, "g1")
		require.NoError(t, err)
		require.Equal(t, 3.1, got1)
	})
	t.Run("Counter", func(t *testing.T) {
		err := storage.Set(context.TODO(), interfaces.MetricTypeCounter, "c1", int64(3))
		require.NoError(t, err)
		err = storage.Set(context.TODO(), interfaces.MetricTypeCounter, "c1", int64(3))
		require.NoError(t, err)
		got1, err := storage.Get(context.TODO(), interfaces.MetricTypeCounter, "c1")
		require.NoError(t, err)
		require.Equal(t, int64(6), got1)
	})
	t.Run("Wrong type", func(t *testing.T) {
		err := storage.Set(context.TODO(), "wrongType", "c1", int64(3))
		require.Error(t, err)
	})
}
