package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis(client *redis.Client) {
	redisClient = client
}

const QuestionQueueKey = "pending_questions"

type QuestionItem struct {
	ID              string   `json:"id"`
	Content         string   `json:"question"`
	ImagePath       string   `json:"answerImage"`
	Answer          string   `json:"parseImage"`
	Explanation     string   `json:"explanation"`
	QuestionType    string   `json:"type"`
	KnowledgePoints []string `json:"knowledgePoints"`
	SubjectID       int32    `json:"subjectId"`
}

func EnqueueQuestion(ctx context.Context, question QuestionItem) error {
	data, err := json.Marshal(question)
	if err != nil {
		return err
	}
	return redisClient.RPush(ctx, QuestionQueueKey, data).Err()
}

func DequeueQuestion(ctx context.Context) (*QuestionItem, error) {
	data, err := redisClient.LPop(ctx, QuestionQueueKey).Bytes()
	if err != nil {
		return nil, err
	}
	var question QuestionItem
	err = json.Unmarshal(data, &question)
	return &question, err
}

func GetQueueStatus(ctx context.Context) (map[string]interface{}, error) {
	length, err := redisClient.LLen(ctx, QuestionQueueKey).Result()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"queue_length": length,
	}, nil
}

func BLPOPQuestion(ctx context.Context) (*QuestionItem, error) {
	result, err := redisClient.BLPop(ctx, 0, QuestionQueueKey).Result()
	if err != nil {
		return nil, err
	}

	if len(result) != 2 {
		return nil, fmt.Errorf("unexpected result length from BLPOP")
	}

	var question QuestionItem
	err = json.Unmarshal([]byte(result[1]), &question)
	if err != nil {
		return nil, err
	}

	return &question, nil
}
