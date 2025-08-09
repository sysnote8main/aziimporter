package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type FluxDBSetting struct {
	Url          string `yaml:"url"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
	Bucket       string `yaml:"bucket"`
}

type Config struct {
	Url   string        `yaml:"url"`
	Token string        `yaml:"token"`
	DB    FluxDBSetting `yaml:"db"`
}

var defaultConfig = Config{
	Url:   "https://api-ktor.azisaba.net",
	Token: "Token is here!",
	DB: FluxDBSetting{
		Url:          "http://localhost:8086",
		Token:        "ExampleToken",
		Organization: "organization",
		Bucket:       "bucket",
	},
}

func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("設定ファイルが見つかりません。デフォルト設定で作成します:", path)
		if err := saveConfigToYAML(path, defaultConfig); err != nil {
			return nil, fmt.Errorf("設定ファイル作成失敗: %w", err)
		}
		panic("設定してから、再実行してください。")
	}

	return loadConfigFromYAML(path)
}

func loadConfigFromYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfigToYAML(path string, cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	header := []byte("# Minecraft server monitor config\n\n")
	return os.WriteFile(path, append(header, data...), 0644)
}
