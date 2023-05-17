package repo

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/alemetric/internal/entity"
)

func TestMetricsRepo_GetMetrics(t *testing.T) {
	metricsRepo := &MetricsRepo{
		Mu:      &sync.Mutex{},
		storage: make(map[string]entity.Metrics),
	}

	var value entity.Gauge = 100.500
	metricsRepo.storage["Frees"] = entity.Metrics{
		ID:    "Frees",
		MType: "gauge",
		Value: &value,
	}

	metricsRepo.storage["Alloc"] = entity.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &value,
	}

	ctx := context.Background()

	for k := range metricsRepo.storage {
		m, err := metricsRepo.GetMetrics(ctx, k)
		require.NoError(t, err)
		require.Equal(t, k, m.ID)
		require.Equal(t, value, *m.Value)
	}

	_, err := metricsRepo.GetMetrics(ctx, "some name")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestMetricsRepo_StoreMetrics(t *testing.T) {
	metricsRepo := &MetricsRepo{
		Mu:      &sync.Mutex{},
		storage: make(map[string]entity.Metrics),
	}

	var value entity.Gauge = 100.500
	frees := entity.Metrics{
		ID:    "Frees",
		MType: "gauge",
		Value: &value,
	}

	alloc := entity.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &value,
	}

	metrics := []entity.Metrics{frees, alloc}

	ctx := context.Background()

	for _, m := range metrics {
		err := metricsRepo.StoreMetrics(ctx, m)
		require.NoError(t, err)
	}

	for _, name := range []string{"Frees", "Alloc"} {
		got, ok := metricsRepo.storage[name]
		require.True(t, ok)
		require.Equal(t, name, got.ID)
	}
}

func TestMetricsRepo_GetMetricsNames(t *testing.T) {
	metricsRepo := &MetricsRepo{
		Mu:      &sync.Mutex{},
		storage: make(map[string]entity.Metrics),
	}

	var value entity.Gauge = 100.500
	metricsRepo.storage["Frees"] = entity.Metrics{
		ID:    "Frees",
		MType: "gauge",
		Value: &value,
	}

	metricsRepo.storage["Alloc"] = entity.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: &value,
	}

	ctx := context.Background()

	expected := make([]string, 0, 2)
	for k := range metricsRepo.storage {
		expected = append(expected, k)
	}

	actual := metricsRepo.GetMetricsNames(ctx)
	require.Equal(t, len(expected), len(actual))

	sort.Strings(expected)
	sort.Strings(actual)
	require.EqualValues(t, expected, actual)
}

func TestUpload(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "metrics")
	require.NoError(t, err)

	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "metrics.json")
	storage := map[string]entity.Metrics{
		"metric1": {ID: "metric1", MType: "type1"},
		"metric2": {ID: "metric2", MType: "type2"},
	}

	data, err := json.Marshal(storage)
	data = append(data, '\n')
	require.NoError(t, err)

	err = os.WriteFile(tempFile, data, 0644)
	require.NoError(t, err)

	repo := &MetricsRepo{
		storage:       make(map[string]entity.Metrics),
		Mu:            &sync.Mutex{},
		StoreFilePath: tempFile,
	}

	err = repo.Upload(context.Background())
	require.NoError(t, err)

	require.Equal(t, len(repo.storage), len(storage))

	for key, value := range storage {
		_, ok := repo.storage[key]
		require.True(t, ok)
		require.Equal(t, value, repo.storage[key])
	}

	emptyFile := filepath.Join(tempDir, "empty.json")
	_, err = os.Create(emptyFile)
	require.NoError(t, err)

	repo.StoreFilePath = emptyFile
	err = repo.Upload(context.Background())
	require.NoError(t, err)
}

func TestStoreAll(t *testing.T) {
	tempFile, err := os.CreateTemp("", "metrics")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	repo := &MetricsRepo{
		storage:       map[string]entity.Metrics{"metric1": {ID: "metric1", MType: "type1"}},
		Mu:            &sync.Mutex{},
		StoreFilePath: tempFile.Name(),
	}

	err = repo.StoreAll()
	require.NoError(t, err)

	_, err = os.ReadFile(tempFile.Name())
	require.NoError(t, err)
}

func TestNewMetricsRepo(t *testing.T) {
	restoreOption := func(r *MetricsRepo) {
		r.Restore = true
	}

	tempFile, err := os.CreateTemp("", "metrics")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	storeFilePathOption := func(r *MetricsRepo) {
		r.StoreFilePath = tempFile.Name()
	}

	repo, err := NewMetricsRepo(restoreOption, storeFilePathOption)
	require.NoError(t, err)

	require.NotNil(t, repo.Mu)
	require.NotNil(t, repo.storage)
	require.Empty(t, repo.storage)
	require.True(t, repo.Restore)
}
