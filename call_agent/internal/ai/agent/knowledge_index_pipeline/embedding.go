package knowledge_index_pipeline

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/embedding"
	"gopkg.in/yaml.v3"
)

type embeddingModelConfig struct {
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url"`
	Model   string `yaml:"model"`
}

type appConfig struct {
	DoubaoEmbeddingModel embeddingModelConfig `yaml:"doubao_embedding_model"`
}

func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	cfg, err := loadAppConfig("config/config.yml")
	if err != nil {
		return nil, err
	}
	m := cfg.DoubaoEmbeddingModel
	if m.APIKey == "" || m.Model == "" {
		return nil, fmt.Errorf("doubao_embedding_model incomplete in config")
	}

	// doubao-embedding-vision-* 不支持 /embeddings，需走 multimodal 接口。
	apiType := ark.APITypeMultiModal
	timeout := 60 * time.Second
	config := &ark.EmbeddingConfig{
		APIKey:  m.APIKey,
		BaseURL: m.BaseURL,
		Model:   m.Model,
		APIType: &apiType,
		Timeout: &timeout,
	}
	eb, err = ark.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}

func loadAppConfig(path string) (*appConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	var cfg appConfig
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	return &cfg, nil
}
