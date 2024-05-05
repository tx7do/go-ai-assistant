package azopenai

import (
	"fmt"
	"log"

	"go-ai-assistant/internal/utils"
)

// AssistantConfig AI智能助理配置
type AssistantConfig struct {
	AssistantId string `json:"assistantId,omitempty" yaml:"assistantId,omitempty"`
	FileId      string `json:"fileId,omitempty" yaml:"fileId,omitempty"`
	ThreadId    string `json:"threadId,omitempty" yaml:"threadId,omitempty"`

	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string `json:"description,omitempty" yaml:"description,omitempty"`
	Instructions string `json:"instructions,omitempty" yaml:"instructions,omitempty"`
}

// Reset 重置
func (c *AssistantConfig) Reset() {
	c.AssistantId = ""
	c.FileId = ""
	c.ThreadId = ""

	c.Name = ""
	c.Description = ""
	c.Instructions = ""
}

func NewAssistantConfig() *AssistantConfig {
	return &AssistantConfig{}
}

// IsValidConfig 是否有效的AI助手配置
func (c *AssistantConfig) IsValidConfig() bool {
	return c.AssistantId != "" && c.ThreadId != "" && c.FileId != ""
}

// getConfigFilePath 获取配置文件的文件路径
func (c *AssistantConfig) getConfigFilePath() string {
	return fmt.Sprintf("%s/assistant.yaml", utils.GetExePath())
}

// LoadConfig 加载配置
func (c *AssistantConfig) LoadConfig() error {
	if err := utils.LoadYaml(c.getConfigFilePath(), &c); err != nil {
		return err
	}

	log.Println("AssistantId :", c.AssistantId)
	log.Println("FileId :", c.FileId)
	log.Println("ThreadId :", c.ThreadId)

	return nil
}

// SaveConfig 保存配置文件
func (c *AssistantConfig) SaveConfig() error {
	return utils.SaveYaml(c.getConfigFilePath(), &c)
}

type OpenAIPlatform struct {
	Key            string `json:"key,omitempty" yaml:"key,omitempty"`
	Endpoint       string `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	DeploymentName string `json:"deployment,omitempty" yaml:"deployment,omitempty"`
}

// AIOptions AI配置
type AIOptions struct {
	EnableAzure bool           `json:"enableAzure,omitempty" yaml:"enableAzure,omitempty"`
	OpenAI      OpenAIPlatform `json:"openai,omitempty" yaml:"openai,omitempty"`
	Azure       OpenAIPlatform `json:"azure,omitempty" yaml:"azure,omitempty"`
}

func NewAIOptions() *AIOptions {
	o := &AIOptions{
		EnableAzure: true,

		OpenAI: OpenAIPlatform{
			Endpoint:       "https://api.openai.com/v1",
			Key:            "",
			DeploymentName: "gpt-4-1106-preview",
		},
		Azure: OpenAIPlatform{
			Endpoint:       "",
			Key:            "",
			DeploymentName: "gpt-4-1106-preview",
		},
	}
	return o
}

func (o *AIOptions) GetDeploymentName() string {
	if o.EnableAzure {
		return o.Azure.DeploymentName
	} else {
		return o.OpenAI.DeploymentName
	}
}

// getConfigFilePath 获取配置文件的文件路径
func (o *AIOptions) getConfigFilePath() string {
	filePath := fmt.Sprintf("%s/openai.yaml", utils.GetExePath())
	return filePath
}

// LoadConfig 加载配置
func (o *AIOptions) LoadConfig() error {
	if err := utils.LoadYaml(o.getConfigFilePath(), &o); err != nil {
		return err
	}
	return nil
}

// SaveConfig 保存配置文件
func (o *AIOptions) SaveConfig() error {
	return utils.SaveYaml(o.getConfigFilePath(), &o)
}
