package openai

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

func CreateThread(
	ctx context.Context,
	openaiClient openai.Client,
	threadMd map[string]string,
) (*openai.Thread, error) {
	return openaiClient.Beta.Threads.New(
		ctx,
		openai.BetaThreadNewParams{
			Metadata: threadMd,
		},
	)
}

func DeleteThread(
	ctx context.Context,
	openaiClient openai.Client,
	threadID string,
) {
	res, err := openaiClient.Beta.Threads.Delete(
		ctx,
		threadID,
	)

	if err != nil {
		log.Println("error while deleting thread: ", err)
	}

	log.Printf("thread with ID(%s) is deleted? %v", threadID, res.Deleted)
}

func AddMessageToThread(
	ctx context.Context,
	openaiClient openai.Client,
	threadID string,
	message string,
	role openai.BetaThreadMessageNewParamsRole) error {
	_, err := openaiClient.Beta.Threads.Messages.New(
		ctx,
		threadID,
		openai.BetaThreadMessageNewParams{
			Content: openai.BetaThreadMessageNewParamsContentUnion{
				OfString: openai.String(message),
			},
			Role: role,
		},
	)

	return err
}

func RunThread(
	ctx context.Context,
	openaiClient openai.Client,
	threadID string,
	assistantID string,
	instructions string,
) error {
	run, err := openaiClient.Beta.Threads.Runs.New(
		ctx,
		threadID,
		openai.BetaThreadRunNewParams{
			AssistantID:  assistantID,
			Instructions: openai.String(instructions),
		},
	)

	if err != nil {
		return err
	}

	timer := time.NewTimer(60 * time.Second)

	for !isFinished(run.Status) {
		time.Sleep(3 * time.Second)

		select {
		case <-timer.C:
			return fmt.Errorf("failed to finish run")
		default:
			run, err = openaiClient.Beta.Threads.Runs.Get(ctx, threadID, run.ID)
			if err != nil {
				return err
			}

			log.Println("waiting for 3 seconds, for run to finish. Current status: ", run.Status)
		}
	}

	timer.Stop()

	return nil
}

func GetThreadMessages(
	ctx context.Context,
	openaiClient openai.Client,
	threadID string) ([]string, error) {
	msgs := make([]string, 0)

	threadMsgs, err := openaiClient.Beta.Threads.Messages.List(
		ctx,
		threadID,
		openai.BetaThreadMessageListParams{})
	if err != nil {
		return nil, err
	}

	for _, msg := range threadMsgs.Data {
		msgs = append(msgs, msg.Content[0].Text.Value)
	}

	after := ""
	if len(threadMsgs.Data) > 0 {
		after = threadMsgs.Data[len(threadMsgs.Data)-1].ID
	}

	for {
		if !threadMsgs.HasMore {
			break
		}

		nextPage, err := openaiClient.Beta.Threads.Messages.List(
			ctx,
			threadID,
			openai.BetaThreadMessageListParams{
				After: param.NewOpt(after),
			})
		if err != nil {
			return nil, err
		}

		if len(nextPage.Data) == 0 {
			break
		}

		for _, msg := range nextPage.Data {
			msgs = append(msgs, msg.Content[0].Text.Value)
		}

		after = nextPage.Data[len(nextPage.Data)-1].ID
		threadMsgs = nextPage
	}

	return msgs, nil
}

func isFinished(runStatus openai.RunStatus) bool {
	return runStatus == openai.RunStatusCancelled || runStatus == openai.RunStatusCompleted
}
