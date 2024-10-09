package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const ImageOCRQueueKey = "image_ocr_queue"

type ImageOCRTask struct {
	ImageURL string `json:"imageUrl"`
	Cookie   string `json:"cookie"`
	Referer  string `json:"referer"`
}

func EnqueueImageOCRTask(ctx context.Context, task ImageOCRTask) error {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return redisClient.RPush(ctx, ImageOCRQueueKey, jsonTask).Err()
}

func DequeueImageOCRTask(ctx context.Context) (*ImageOCRTask, error) {
	result, err := redisClient.LPop(ctx, ImageOCRQueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 队列为空
		}
		return nil, err
	}

	var task ImageOCRTask
	err = json.Unmarshal([]byte(result), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func BLPOPImageOCRTask(ctx context.Context) (*ImageOCRTask, error) {
	result, err := redisClient.BLPop(ctx, 0, ImageOCRQueueKey).Result()
	if err != nil {
		return nil, err
	}

	if len(result) != 2 {
		return nil, fmt.Errorf("unexpected result length from BLPOP")
	}

	var task ImageOCRTask
	err = json.Unmarshal([]byte(result[1]), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func GetImageOCRQueueLength(ctx context.Context) (int64, error) {
	return redisClient.LLen(ctx, ImageOCRQueueKey).Result()
}
