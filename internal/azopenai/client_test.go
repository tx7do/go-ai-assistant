package azopenai_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-ai-assistant/internal/azopenai"
)

func TestClient(t *testing.T) {
	c := azopenai.NewClient()
	assert.NotNil(t, c)

	_, _ = c.SendChatForAssistant("What was my last question?", nil)
}

func TestConfig(t *testing.T) {
	var opt = azopenai.NewAIOptions()
	_ = opt.SaveConfig()
}
