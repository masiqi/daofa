package queue

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis(client *redis.Client) {
	redisClient = client
}

const questionQueueKey = "pending_questions"

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
	return redisClient.RPush(ctx, questionQueueKey, data).Err()
}

func DequeueQuestion(ctx context.Context) (*QuestionItem, error) {
	data, err := redisClient.LPop(ctx, questionQueueKey).Bytes()
	if err != nil {
		return nil, err
	}
	var question QuestionItem
	err = json.Unmarshal(data, &question)
	return &question, err
}

func GetQueueStatus(ctx context.Context) (map[string]interface{}, error) {
	length, err := redisClient.LLen(ctx, questionQueueKey).Result()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"queue_length": length,
	}, nil
}
