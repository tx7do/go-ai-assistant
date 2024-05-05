package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-ai-assistant/internal/utils"
)

type Config struct {
	ID string `json:"id,omitempty" yaml:"id,omitempty"`
}

func TestWriteBinaryFile(t *testing.T) {
}

func TestYAML(t *testing.T) {
	cfg := Config{
		ID: "123",
	}

	filePath := "test.yaml"

	var err error

	err = utils.SaveYaml(filePath, &cfg)
	assert.Nil(t, err)

	var cfg1 Config
	err = utils.LoadYaml(filePath, &cfg1)
	assert.Nil(t, err)
	assert.Equal(t, cfg1.ID, cfg.ID)
}
