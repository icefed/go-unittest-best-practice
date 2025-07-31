package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// dir, err := os.MkdirTemp("", "testloadconfig")
	// if err != nil {
	// 	t.Fatalf(err.Error())
	// }
	// defer os.RemoveAll(dir)

	// after go1.15
	dir := t.TempDir()

	configPath := filepath.Join(dir, "config.yaml")

	t.Run("file not found", func(t *testing.T) {
		_, err := LoadConfig(filepath.Join(dir, "config2.yaml"))
		if err == nil {
			t.Errorf("expect file not found error, but got nil error")
		} else {
			if !strings.Contains(err.Error(), "no such file") {
				t.Errorf("expect file not found error, but got %v", err)
			}
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		err := os.WriteFile(configPath,
			[]byte(`
dbhost: 127.0.0.1
dbport: 3306
dbuser: root
dbpassword: 123456
dbname= test
listenPort: 6666`), 0644)
		if err != nil {
			t.Fatal(err)
		}
		_, err = LoadConfig(configPath)
		if err == nil {
			t.Errorf("expect invalid yaml error, but got nil error")
		} else {
			if !strings.Contains(err.Error(), "parse yaml failed") {
				t.Errorf("expect invalid yaml error, but got \"%v\"", err)
			}
		}
	})

	t.Run("success", func(t *testing.T) {
		err := os.WriteFile(configPath,
			[]byte(`
dbhost: 127.0.0.1
dbport: 3306
dbuser: root
dbpassword: 123456
dbname: test
listenPort: 6666
pprofAddr: :7777`), 0644)
		if err != nil {
			t.Fatal(err)
		}
		conf, err := LoadConfig(configPath)
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}
		assertConfig(t, &Config{
			DBHost:     "127.0.0.1",
			DBPort:     3306,
			DBUser:     "root",
			DBPassword: "123456",
			DBName:     "test",
			ListenPort: 6666,
			PprofAddr:  ":7777",
		}, conf)
	})
}

func assertConfig(t *testing.T, expectedConfig *Config, conf *Config) {
	t.Helper()

	if expectedConfig == nil || conf == nil {
		if expectedConfig != nil || conf != nil {
			t.Errorf("assert config not equal, expected: %v, actual: %v", expectedConfig, conf)
		}
		return
	}
	if expectedConfig.DBHost != conf.DBHost ||
		expectedConfig.DBPort != conf.DBPort ||
		expectedConfig.DBUser != conf.DBUser ||
		expectedConfig.DBPassword != conf.DBPassword ||
		expectedConfig.DBName != conf.DBName ||
		expectedConfig.ListenPort != conf.ListenPort ||
		expectedConfig.PprofAddr != conf.PprofAddr {
		t.Errorf("assert config not equal, expected: %v, actual: %v", expectedConfig, conf)
	}
}