package main

import (
	"context"

	"go-ai-assistant/internal/azopenai"
)

// App struct
type App struct {
	ctx    context.Context
	client *azopenai.Client
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	c := azopenai.NewClient()
	a.client = c
}

// IsValidAssistantConfig 是否有效的AI助手配置
func (a *App) IsValidAssistantConfig() bool {
	return a.client.IsValidAssistantConfig()
}

// GetAssistantConfig 获取AI助手配置
func (a *App) GetAssistantConfig() *azopenai.AssistantConfig {
	return a.client.GetAssistantConfig()
}

// InitAssistant 初始化助理
func (a *App) InitAssistant(assistantName, description, instructions, fileName string, fileData []byte) error {
	// 加载文件
	var fileId string
	var err error
	if fileId, err = a.client.UploadFileWithBinaryData(a.ctx, fileName, fileData, true); err != nil {
		return err
	}

	return a.client.InitAssistant(a.ctx, assistantName, description, instructions, fileId)
}

// DestroyCurrentAssistant 销毁掉当前的助理
func (a *App) DestroyCurrentAssistant() {
	// 删除文件
	_ = a.client.DestroyCurrentFile(a.ctx)

	// 删除会话
	_ = a.client.DestroyCurrentThread(a.ctx)

	// 删除助手
	_ = a.client.DestroyCurrentAssistant(a.ctx)
}

// SendChatForAssistant 发送聊天到助理
func (a *App) SendChatForAssistant(chatMessage string) ([]string, error) {
	return a.client.SendChatForAssistant(chatMessage, nil)
}

// UploadFileAndAssociateAssistant 上传文件并关联到助理
func (a *App) UploadFileAndAssociateAssistant(uploadFilePath string) error {
	return a.client.UploadFileAndAssociateAssistant(a.ctx, uploadFilePath)
}
