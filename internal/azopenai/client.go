package azopenai

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/sashabaranov/go-openai"

	"go-ai-assistant/internal/utils"
)

type Client struct {
	cli *openai.Client

	cfg *AssistantConfig
	opt *AIOptions
}

func NewClient() *Client {
	cli := &Client{
		cfg: NewAssistantConfig(),
		opt: NewAIOptions(),
	}

	_ = cli.opt.LoadConfig()
	_ = cli.cfg.LoadConfig()

	cli.cli = openai.NewClient(cli.opt.OpenAI.Key)

	return cli
}

// IsValidAssistantConfig 是否有效的AI助手配置
func (c *Client) IsValidAssistantConfig() bool {
	return c.cfg.IsValidConfig()
}

// GetAssistantConfig 获取AI助手配置
func (c *Client) GetAssistantConfig() *AssistantConfig {
	return c.cfg
}

// UploadFileWithBinaryData 上传文件
func (c *Client) UploadFileWithBinaryData(ctx context.Context, fileName string, fileData []byte, isAssistant bool) (string, error) {
	purpose := "assistants"
	if !isAssistant {
		purpose = "fine-tune"
	}
	// go-openai v1.41.2 只支持 FilePath 字段，需写入临时文件
	tmpFile, err := utils.WriteTempFile(fileName, fileData)
	if err != nil {
		log.Println("write temp file failed: ", err.Error())
		return "", err
	}
	defer utils.RemoveFile(tmpFile)

	resp, err := c.cli.CreateFile(ctx, openai.FileRequest{
		FileName: fileName,
		Purpose:  purpose,
		FilePath: tmpFile,
	})
	if err != nil {
		log.Println("upload file failed: ", err.Error())
		return "", err
	}
	log.Printf("upload file: [%s]", resp.ID)
	c.cfg.FileId = resp.ID
	_ = c.cfg.SaveConfig()
	return resp.ID, nil
}

// UploadFile 上传文件
func (c *Client) UploadFile(ctx context.Context, uploadFilePath string, isAssistant bool) (string, error) {
	if len(uploadFilePath) == 0 {
		return "", nil
	}

	fileData, err := utils.ReadBinaryFile(uploadFilePath)
	if err != nil {
		log.Println("read file failed: ", err.Error())
		return "", err
	}
	fileName := filepath.Base(uploadFilePath)

	return c.UploadFileWithBinaryData(ctx, fileName, fileData, isAssistant)
}

// DeleteFile 删除文件
func (c *Client) DeleteFile(ctx context.Context, fileId string) error {
	err := c.cli.DeleteFile(ctx, fileId)
	return err
}

func (c *Client) DestroyCurrentFile(ctx context.Context) error {
	err := c.DeleteFile(ctx, c.cfg.FileId)
	c.cfg.FileId = ""
	_ = c.cfg.SaveConfig()
	return err
}

// DeleteThread 删除对话
func (c *Client) DeleteThread(ctx context.Context, threadID string) error {
	_, err := c.cli.DeleteThread(ctx, threadID)
	return err
}

func (c *Client) DestroyCurrentThread(ctx context.Context) error {
	err := c.DeleteThread(ctx, c.cfg.ThreadId)
	c.cfg.ThreadId = ""
	_ = c.cfg.SaveConfig()
	return err
}

// DeleteAssistant 删除助手
func (c *Client) DeleteAssistant(ctx context.Context, assistantID string) error {
	_, err := c.cli.DeleteAssistant(ctx, assistantID)
	return err
}

func (c *Client) DestroyCurrentAssistant(ctx context.Context) error {
	err := c.DeleteAssistant(ctx, c.cfg.AssistantId)
	c.cfg.AssistantId = ""
	_ = c.cfg.SaveConfig()
	return err
}

// UploadFileAndAssociateAssistant 上传文件并关联到助理
func (c *Client) UploadFileAndAssociateAssistant(ctx context.Context, uploadFilePath string) error {
	var err error

	// 解除之前的关联
	if c.cfg.FileId != "" {
		_ = c.DeleteFile(ctx, c.cfg.FileId)

		_ = c.deleteAssistantFile(c.cfg.AssistantId, c.cfg.FileId)
		c.cfg.FileId = ""
	}

	// 加载文件
	var fileId string
	if fileId, err = c.UploadFile(ctx, uploadFilePath, true); err != nil {
		return err
	}

	c.cfg.FileId = fileId
	_ = c.cfg.SaveConfig()

	// 关联文件
	return c.associateAssistantFile(c.cfg.AssistantId, c.cfg.FileId)
}

// associateAssistantFile 助手绑定文件
func (c *Client) associateAssistantFile(assistantId, fileId string) error {
	ctx := context.Background()
	// go-openai 没有 AttachFileToAssistant，需用 UpdateAssistant 追加文件
	assistant, err := c.cli.RetrieveAssistant(ctx, assistantId)
	if err != nil {
		return err
	}
	fileIDs := append(assistant.FileIDs, fileId)
	_, err = c.cli.ModifyAssistant(ctx, assistantId, openai.AssistantRequest{
		Name:         assistant.Name,
		Description:  assistant.Description,
		Instructions: assistant.Instructions,
		Tools:        assistant.Tools,
		FileIDs:      fileIDs,
	})
	return err
}

// deleteAssistantFile 删除掉助手绑定的文件
func (c *Client) deleteAssistantFile(assistantId, fileId string) error {
	ctx := context.Background()
	assistant, err := c.cli.RetrieveAssistant(ctx, assistantId)
	if err != nil {
		return err
	}
	var newFileIDs []string
	for _, fid := range assistant.FileIDs {
		if fid != fileId {
			newFileIDs = append(newFileIDs, fid)
		}
	}
	_, err = c.cli.ModifyAssistant(ctx, assistantId, openai.AssistantRequest{
		Name:         assistant.Name,
		Description:  assistant.Description,
		Instructions: assistant.Instructions,
		Tools:        assistant.Tools,
		FileIDs:      newFileIDs,
	})
	return err
}

// InitAssistantWithUploadFile 创建助手
func (c *Client) InitAssistantWithUploadFile(ctx context.Context, assistantName, description, instructions, uploadFilePath string) error {
	// 加载文件
	var fileId string
	var err error
	if fileId, err = c.UploadFile(ctx, uploadFilePath, true); err != nil {
		return err
	}

	c.cfg.FileId = fileId
	_ = c.cfg.SaveConfig()

	return c.InitAssistant(ctx, assistantName, description, instructions, fileId)
}

// InitAssistant 创建助手
func (c *Client) InitAssistant(ctx context.Context, assistantName, description, instructions, fileId string) error {
	var fileIDs []string
	if c.cfg.FileId != "" {
		fileIDs = append(fileIDs, fileId)
	}
	assistantReq := openai.AssistantRequest{
		Name:         &assistantName,
		Description:  &description,
		Instructions: &instructions,
		FileIDs:      fileIDs,
		Tools:        []openai.AssistantTool{{Type: "code_interpreter"}},
	}
	assistant, err := c.cli.CreateAssistant(ctx, assistantReq)
	if err != nil {
		log.Println("create assistant failed: ", err.Error())
		return err
	}
	log.Printf("create assistant: [%s]", assistant.ID)

	thread, err := c.cli.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		log.Println("create thread failed: ", err.Error())
		return err
	}
	log.Printf("create thread: [%s]", thread.ID)

	c.cfg.AssistantId = assistant.ID
	c.cfg.ThreadId = thread.ID
	c.cfg.FileId = fileId

	c.cfg.Name = assistantName
	c.cfg.Description = description
	c.cfg.Instructions = instructions
	_ = c.cfg.SaveConfig()

	return nil
}

// pollRunEnd 已废弃，go-openai SDK 不需要此方法

// SendChatForAssistant 发送聊天到助理
func (c *Client) SendChatForAssistant(message string, threadID *string) ([]string, error) {
	ctx := context.Background()

	if threadID == nil || len(*threadID) == 0 {
		threadID = &c.cfg.ThreadId
	}

	// 创建消息
	_, err := c.cli.CreateMessage(ctx, *threadID, openai.MessageRequest{
		Role:    "user",
		Content: message,
	})
	if err != nil {
		return nil, err
	}

	// 创建 Run
	run, err := c.cli.CreateRun(ctx, *threadID, openai.RunRequest{
		AssistantID: c.cfg.AssistantId,
	})
	if err != nil {
		return nil, err
	}

	// 轮询 Run 状态直到完成
	for {
		runStatus, err := c.cli.RetrieveRun(ctx, *threadID, run.ID)
		if err != nil {
			return nil, err
		}
		if runStatus.Status == "completed" {
			break
		}
		if runStatus.Status == "failed" || runStatus.Status == "cancelled" || runStatus.Status == "expired" {
			return nil, fmt.Errorf("run ended but status was not complete: %s", runStatus.Status)
		}
		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// 获取消息
	resp, err := c.cli.ListMessage(ctx, *threadID, nil, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var chatMessages []string
	for _, msg := range resp.Messages {
		if msg.Role == "assistant" {
			for _, c := range msg.Content {
				if c.Type == "text" && c.Text != nil {
					chatMessages = append(chatMessages, c.Text.Value)
				}
			}
		}
	}
	return chatMessages, nil
}
