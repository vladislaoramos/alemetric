package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
				Name:           agentName + "1",
				PollInterval:   time.Second * 3,
				ReportInterval: time.Second * 9,
				ServerURL:      "0.0.0.0:8888",
				MetricsNames:   metricsNames,
				RateLimit:      rateLimit * 2,
				Key:            "key",
			},
		}

		cfg.updateAgentConfigs(flags)

		require.Equal(t, "0.0.0.0:8888", cfg.Agent.ServerURL)
		require.Equal(t, time.Second*3, cfg.Agent.PollInterval)
		require.Equal(t, time.Second*9, cfg.Agent.ReportInterval)
		require.Equal(t, uint(2), cfg.Agent.RateLimit)
		require.Equal(t, "alemetric-agent1", cfg.Agent.Name)
		require.Equal(t, "key", cfg.Agent.Key)
	})

	t.Run("server update", func(t *testing.T) {
		cfg := defaultServerCfg()

		flags := &Config{
			Server: Server{
				Name:          serverName + "1",
				Address:       "0.0.0.0:8888",
				StoreInterval: time.Second * 60,
				StoreFile:     "file.json",
				Restore:       false,
				Key:           "key",
			},
			Database: Database{
				URL: "url",
			},
			Logger: Logger{
				Level: "debug",
			},
		}

		cfg.updateServerConfigs(flags)

		require.Equal(t, "0.0.0.0:8888", cfg.Server.Address)
		require.Equal(t, time.Second*60, cfg.Server.StoreInterval)
		require.Equal(t, "file.json", cfg.Server.StoreFile)
		require.Equal(t, false, cfg.Server.Restore)
		require.Equal(t, "alemetric-server1", cfg.Server.Name)
		require.Equal(t, "key", cfg.Server.Key)
		require.Equal(t, "url", cfg.Database.URL)
		require.Equal(t, "debug", cfg.Logger.Level)
	})
}

func TestServerJSONConfig(t *testing.T) {
	t.Run("envs with json", func(t *testing.T) {
		err := cleanEnvs()
		require.NoError(t, err)

		err = os.Setenv("NAME", "my-server")
		require.NoError(t, err)

		err = os.Setenv("STORE_INTERVAL", "250s")
		require.NoError(t, err)

		err = os.Setenv("RESTORE", "true")
		require.NoError(t, err)

		err = os.Setenv("CONFIG", serverJSONConfig)
		require.NoError(t, err)

		cfg := new(Config)

		jsonConfig, err := loadServerJSONConfig("server.json")
		require.NoError(t, err)
		cfg.updateServerConfigs(jsonConfig)

		envs := new(Config)
		err = cleanenv.ReadEnv(envs)
		require.NoError(t, err)
		cfg.updateServerConfigs(envs)

		require.Equal(t, time.Second*250, cfg.Server.StoreInterval)
		require.Equal(t, "my-server", cfg.Server.Name)
		require.Equal(t, true, cfg.Server.Restore)
		require.Equal(t, "/path/to/file.db", cfg.Server.StoreFile)
		require.Equal(t, "/path/to/key.pem", cfg.Server.CryptoKey)
	})

	t.Run("flags with json", func(t *testing.T) {
		cfg := new(Config)

		jsonConfig, err := loadServerJSONConfig("server.json")
		require.NoError(t, err)
		cfg.updateServerConfigs(jsonConfig)

		flags := &Config{
			Server: Server{
				Name:          serverName + "1",
				StoreInterval: time.Second * 250,
				Restore:       true,
				Key:           "key",
				StoreFile:     "file.db",
			},
		}

		cfg.updateServerConfigs(flags)

		require.Equal(t, time.Second*250, cfg.Server.StoreInterval)
		require.Equal(t, serverName+"1", cfg.Server.Name)
		require.Equal(t, true, cfg.Server.Restore)
		require.Equal(t, "file.db", cfg.Server.StoreFile)
		require.Equal(t, "/path/to/key.pem", cfg.Server.CryptoKey)
	})
}

func TestAgentJSONConfig(t *testing.T) {
	t.Run("envs with json", func(t *testing.T) {
		err := cleanEnvs()
		require.NoError(t, err)

		err = os.Setenv("NAME", "my-agent")
		require.NoError(t, err)

		err = os.Setenv("REPORT_INTERVAL", "25s")
		require.NoError(t, err)

		err = os.Setenv("CONFIG", agentJSONConfig)
		require.NoError(t, err)

		cfg := new(Config)

		jsonConfig, err := loadAgentJSONConfig("agent.json")
		require.NoError(t, err)
		cfg.updateAgentConfigs(jsonConfig)

		envs := new(Config)
		err = cleanenv.ReadEnv(envs)
		require.NoError(t, err)
		cfg.updateAgentConfigs(envs)

		require.Equal(t, time.Second*25, cfg.Agent.ReportInterval)
		require.Equal(t, time.Second, cfg.Agent.PollInterval)
		require.Equal(t, "localhost:8080", cfg.Agent.ServerURL)
		require.Equal(t, "/path/to/key.pem", cfg.Agent.CryptoKey)
	})

	t.Run("flags with json", func(t *testing.T) {
		cfg := new(Config)

		jsonConfig, err := loadAgentJSONConfig("agent.json")
		require.NoError(t, err)
		cfg.updateAgentConfigs(jsonConfig)

		flags := &Config{
			Agent: Agent{
				ReportInterval: time.Second * 9,
				ServerURL:      "0.0.0.0:8888",
				Key:            "key",
			},
		}

		cfg.updateAgentConfigs(flags)

		require.Equal(t, time.Second*9, cfg.Agent.ReportInterval)
		require.Equal(t, time.Second, cfg.Agent.PollInterval)
		require.Equal(t, "0.0.0.0:8888", cfg.Agent.ServerURL)
		require.Equal(t, "/path/to/key.pem", cfg.Agent.CryptoKey)
		require.Equal(t, "key", cfg.Agent.Key)
	})
}

func cleanEnvs() error {
	stringEnvs := []string{
		"NAME", "ADDRESS", "STORE_INTERVAL", "STORE_FILE",
		"KEY", "CRYPTO_KEY", "POLL_INTERVAL", "DATABASE_DSN",
	}

	for _, name := range stringEnvs {
		err := os.Setenv(name, "")
		if err != nil {
			return fmt.Errorf("could not erase env var with name %s", name)
		}
	}

	timeEnvs := []string{
		"STORE_INTERVAL", "POLL_INTERVAL", "REPORT_INTERVAL",
	}

	for _, name := range timeEnvs {
		err := os.Setenv(name, "0s")
		if err != nil {
			return fmt.Errorf("could not erase env var with name %s", name)
		}
	}

	err := os.Setenv("RATE_LIMIT", "1")
	if err != nil {
		return fmt.Errorf("could not erase env var with name %s", "RATE_LIMIT")
	}

	err = os.Setenv("RESTORE", "f")
	if err != nil {
		return fmt.Errorf("could not erase env var with name %s", "RATE_LIMIT")
	}

	return nil
}
