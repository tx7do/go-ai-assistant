package azopenai

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"

	"go-ai-assistant/internal/utils"
)

type Client struct {
	cli *azopenaiassistants.Client

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

	cli.cli = cli.createAssistantClient(cli.opt)

	return cli
}

// createAssistantClient 创建助手客户端
func (c *Client) createAssistantClient(opt *AIOptions) *azopenaiassistants.Client {
	if opt.EnableAzure {
		return c.createAssistantClientFromAzure(opt)
	} else {
		return c.createAssistantClientFromOpenAI(opt)
	}
}

func (c *Client) createAssistantClientFromAzure(opt *AIOptions) *azopenaiassistants.Client {
	opts := &azopenaiassistants.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Logging: policy.LogOptions{
				IncludeBody: true,
			},
		},
	}

	tmpClient, err := azopenaiassistants.NewClientWithKeyCredential(opt.Azure.Endpoint, azcore.NewKeyCredential(opt.Azure.Key), opts)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return tmpClient
}

func (c *Client) createAssistantClientFromOpenAI(opt *AIOptions) *azopenaiassistants.Client {
	opts := &azopenaiassistants.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Logging: policy.LogOptions{
				IncludeBody: true,
			},
		},
	}

	tmpClient, err := azopenaiassistants.NewClientForOpenAI(opt.OpenAI.Endpoint, azcore.NewKeyCredential(opt.OpenAI.Key), opts)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return tmpClient
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
	var filePurpose azopenaiassistants.FilePurpose
	if isAssistant {
		filePurpose = azopenaiassistants.FilePurposeAssistants
	} else {
		filePurpose = azopenaiassistants.FilePurposeFineTune
	}

	// 上传文件
	uploadResp, err := c.cli.UploadFile(ctx, bytes.NewReader(fileData),
		filePurpose,
		&azopenaiassistants.UploadFileOptions{
			Filename: to.Ptr(fileName),
		})
	if err != nil {
		log.Println("upload file failed: ", err.Error())
		return "", err
	}
	log.Printf("upload file: [%s][%s]", *uploadResp.ID, *uploadResp.Filename)

	c.cfg.FileId = *uploadResp.ID
	_ = c.cfg.SaveConfig()

	return *uploadResp.ID, nil
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

	_, err := c.cli.DeleteFile(ctx, fileId, nil)
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
	_, err := c.cli.DeleteThread(ctx, threadID, nil)
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
	_, err := c.cli.DeleteAssistant(ctx, assistantID, nil)
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
	_, err := c.cli.CreateAssistantFile(ctx, assistantId, azopenaiassistants.CreateAssistantFileBody{
		FileID: &fileId,
	}, nil)
	return err
}

// deleteAssistantFile 删除掉助手绑定的文件
func (c *Client) deleteAssistantFile(assistantId, fileId string) error {
	ctx := context.Background()
	_, err := c.cli.DeleteAssistantFile(ctx, assistantId, fileId, nil)
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
	var err error

	var fileIDs []string
	if c.cfg.FileId != "" {
		fileIDs = append(fileIDs, fileId)
	}

	// 创建助手
	createAssistantResp, err := c.cli.CreateAssistant(ctx, azopenaiassistants.AssistantCreationBody{
		Name:           to.Ptr(assistantName),
		DeploymentName: to.Ptr(c.opt.GetDeploymentName()),
		Description:    to.Ptr(description),
		Instructions:   to.Ptr(instructions),
		FileIDs:        fileIDs,
		Tools: []azopenaiassistants.ToolDefinitionClassification{
			&azopenaiassistants.CodeInterpreterToolDefinition{},
			//&assistants.RetrievalToolDefinition{},
		},
	}, nil)
	if err != nil {
		log.Println("create assistant failed: ", err.Error())
		return err
	}
	log.Printf("create assistant: [%s][%s]", *createAssistantResp.ID, *createAssistantResp.DeploymentName)

	// 创建线程
	createThreadResp, err := c.cli.CreateThread(ctx, azopenaiassistants.AssistantThreadCreationOptions{}, nil)
	if err != nil {
		log.Println("create thread failed: ", err.Error())
		return err
	}
	log.Printf("create thread: [%s]", *createThreadResp.ID)

	c.cfg.AssistantId = *createAssistantResp.ID
	c.cfg.ThreadId = *createThreadResp.ID
	c.cfg.FileId = fileId

	c.cfg.Name = assistantName
	c.cfg.Description = description
	c.cfg.Instructions = instructions
	_ = c.cfg.SaveConfig()

	return nil
}

func (c *Client) pollRunEnd(ctx context.Context, client *azopenaiassistants.Client, threadID string, runID string) error {
	for {
		lastGetRunResp, err := client.GetRun(ctx, threadID, runID, nil)

		if err != nil {
			return err
		}

		if *lastGetRunResp.Status != azopenaiassistants.RunStatusQueued && *lastGetRunResp.Status != azopenaiassistants.RunStatusInProgress {
			if *lastGetRunResp.Status == azopenaiassistants.RunStatusCompleted {
				return nil
			}

			return fmt.Errorf("run ended but status was not complete: %s", *lastGetRunResp.Status)
		}

		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// SendChatForAssistant 发送聊天到助理
func (c *Client) SendChatForAssistant(message string, threadID *string) ([]string, error) {
	ctx := context.Background()

	var fileIDs []string
	if c.cfg.FileId != "" {
		fileIDs = append(fileIDs, c.cfg.FileId)
	}

	if threadID == nil || len(*threadID) == 0 {
		threadID = &c.cfg.ThreadId
	}

	_, err := c.cli.CreateMessage(ctx, *threadID, azopenaiassistants.CreateMessageBody{
		Content: to.Ptr(message),
		Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
		FileIDs: fileIDs,
	}, nil)
	if err != nil {
		return nil, err
	}
	//log.Printf("create message: [%s]", *createMessageResp.ID)

	createRunResp, err := c.cli.CreateRun(ctx, *threadID, azopenaiassistants.CreateRunBody{
		AssistantID: &c.cfg.AssistantId,
	}, nil)
	if err != nil {
		return nil, err
	}
	//log.Printf("create run: [%s]", *createRunResp.ID)
	runId := *createRunResp.ID

	err = c.pollRunEnd(ctx, c.cli, *threadID, runId)

	var chatMessages []string
	listMessagesPager := c.cli.NewListRunStepsPager(*threadID, runId, nil)
	for listMessagesPager.More() {
		page, err := listMessagesPager.NextPage(ctx)
		if err != nil {
			continue
		}

		for _, runStep := range page.Data {
			_, err := c.cli.GetRunStep(ctx, *threadID, runId, *runStep.ID, nil)
			if err != nil {
				continue
			}
			//log.Printf("get run step: [%s][%s][%s]", *rereadRunStep.ID, *rereadRunStep.ThreadID, *rereadRunStep.RunID)

			if runStep.StepDetails == nil {
				continue
			}

			switch t := runStep.StepDetails.(type) {
			case *azopenaiassistants.RunStepMessageCreationDetails:
				messageResp, err := c.cli.GetMessage(ctx, *threadID, *t.MessageCreation.MessageID, nil)
				if err != nil {
					continue
				}
				if *messageResp.Role == azopenaiassistants.MessageRoleAssistant {
					body := *messageResp.Content[0].(*azopenaiassistants.MessageTextContent).Text.Value

					fmt.Printf("Assistant response: %s\n", body)

					chatMessages = append(chatMessages, body)
				}

			case *azopenaiassistants.RunStepToolCallDetails:
			}

		}
	}

	return chatMessages, nil
}
