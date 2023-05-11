package configs

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("agent config", func(t *testing.T) {
		t.Parallel()

		err := os.Setenv("POLL_INTERVAL", "1s")
		require.NoError(t, err)

		err = os.Setenv("REPORT_INTERVAL", "2s")
		require.NoError(t, err)

		err = os.Setenv("ADDRESS", "http://localhost:8080")
		require.NoError(t, err)

		cfg := defaultAgentCfg()

		require.Equal(t, "127.0.0.1:8080", cfg.Agent.ServerURL)
		require.Equal(t, time.Second*2, cfg.Agent.PollInterval)
		require.Equal(t, time.Second*10, cfg.Agent.ReportInterval)

		require.Equal(t, agentName, cfg.Agent.Name)
		require.Equal(t, uint(rateLimit), cfg.Agent.RateLimit)
	})

	t.Run("server config", func(t *testing.T) {
		t.Parallel()

		err := os.Setenv("STORE_INTERVAL", "300s")
		require.NoError(t, err)

		err = os.Setenv("STORE_FILE", "/tmp/devops-metrics-db.json")
		require.NoError(t, err)

		err = os.Setenv("RESTORE", "true")
		require.NoError(t, err)

		err = os.Setenv("ADDRESS", "http://localhost:8080")
		require.NoError(t, err)

		cfg := defaultServerCfg()

		require.Equal(t, "127.0.0.1:8080", cfg.Server.Address)
		require.Equal(t, "/tmp/devops-metrics-db.json", cfg.Server.StoreFile)
		require.Equal(t, time.Second*300, cfg.Server.StoreInterval)
		require.Equal(t, true, cfg.Server.Restore)

		require.Equal(t, serverName, cfg.Server.Name)
	})
}

func TestUpdateConfig(t *testing.T) {
	t.Run("agent update", func(t *testing.T) {
		cfg := defaultAgentCfg()

		flags := &Config{
			Agent: Agent{
				Name:           agentName,
				PollInterval:   time.Second * 3,
				ReportInterval: time.Second * 9,
				ServerURL:      "0.0.0.0:8888",
				MetricsNames:   metricsNames,
				RateLimit:      rateLimit * 2,
			},
		}

		cfg.updateAgentConfigs(flags)

		require.Equal(t, "0.0.0.0:8888", cfg.Agent.ServerURL)
		require.Equal(t, time.Second*3, cfg.Agent.PollInterval)
		require.Equal(t, time.Second*9, cfg.Agent.ReportInterval)
		require.Equal(t, uint(2), cfg.Agent.RateLimit)
	})

	t.Run("server update", func(t *testing.T) {
		cfg := defaultServerCfg()

		flags := &Config{
			Server: Server{
				Name:          serverName,
				Address:       "0.0.0.0:8888",
				StoreInterval: time.Second * 60,
				StoreFile:     "file.json",
				Restore:       false,
			},
		}

		cfg.updateServerConfigs(flags)

		require.Equal(t, "0.0.0.0:8888", cfg.Server.Address)
		require.Equal(t, time.Second*60, cfg.Server.StoreInterval)
		require.Equal(t, "file.json", cfg.Server.StoreFile)
		require.Equal(t, false, cfg.Server.Restore)
	})
}
