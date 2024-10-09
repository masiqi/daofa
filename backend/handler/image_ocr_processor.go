package handler

import (
	"bytes"
	"context"
	"crypto/md5"
	"daofa/backend/dal"
	"daofa/backend/model"
	"daofa/backend/queue"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func ProcessImageOCRTasks(ctx context.Context) {
	for {
		task, err := queue.BLPOPImageOCRTask(ctx)
		if err != nil {
			fmt.Printf("从队列获取任务失败: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		// 处理任务
		err = processImageOCRTask(ctx, task)
		if err != nil {
			fmt.Printf("处理图片OCR任务失败: %v\n", err)
		}
	}
}

func processImageOCRTask(ctx context.Context, task *queue.ImageOCRTask) error {
	// 创建任务记录
	status := "processing"
	ocrTask := &model.ImageOcrTask{
		ImageURL: task.ImageURL,
		Cookie:   &task.Cookie,
		Referer:  &task.Referer,
		Status:   &status,
	}
	err := dal.Q.ImageOcrTask.WithContext(ctx).Create(ocrTask)
	if err != nil {
		return fmt.Errorf("创建任务记录失败: %v", err)
	}

	// 下载图片
	localPath, err := downloadImage(task.ImageURL, task.Cookie, task.Referer)
	if err != nil {
		return updateTaskStatus(ctx, ocrTask, "failed", fmt.Errorf("下载图片失败: %v", err))
	}
	ocrTask.LocalFilePath = &localPath

	// 执行OCR
	ocrResult, err := performImageOCR(localPath)
	if err != nil {
		return updateTaskStatus(ctx, ocrTask, "failed", fmt.Errorf("执行OCR失败: %v", err))
	}

	// 更新任务记录
	ocrTask.OcrResult = &ocrResult
	return updateTaskStatus(ctx, ocrTask, "completed", nil)
}

func updateTaskStatus(ctx context.Context, task *model.ImageOcrTask, status string, err error) error {
	task.Status = &status
	_, updateErr := dal.Q.ImageOcrTask.WithContext(ctx).Where(
		dal.Q.ImageOcrTask.ID.Eq(task.ID),
	).UpdateSimple(
		dal.Q.ImageOcrTask.Status.Value(status),
		dal.Q.ImageOcrTask.OcrResult.Value(*task.OcrResult),
		dal.Q.ImageOcrTask.LocalFilePath.Value(*task.LocalFilePath),
	)
	if updateErr != nil {
		fmt.Printf("更新任务状态失败: %v\n", updateErr)
	}
	return err
}

func downloadImage(url, cookie, referer string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载图片失败,状态码: %d", resp.StatusCode)
	}

	// 读取图片内容
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 计算文件哈希
	hash := md5.Sum(imageData)
	hashString := hex.EncodeToString(hash[:])

	// 获取文件扩展名
	urlParts := strings.Split(url, "?")
	fileName := filepath.Base(urlParts[0])
	fileExt := filepath.Ext(fileName)

	// 创建目录结构
	storagePath := viper.GetString("IMAGE_STORAGE_PATH")
	dirPath := filepath.Join(storagePath, hashString[:2], hashString[2:4])
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建本地文件
	localFileName := fmt.Sprintf("%s%s", hashString, fileExt)
	localPath := filepath.Join(dirPath, localFileName)
	err = os.WriteFile(localPath, imageData, 0644)
	if err != nil {
		return "", err
	}

	return localPath, nil
}

func performImageOCR(imagePath string) (string, error) {
	ocrURL := viper.GetString("OCR_URL")

	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	part, err := writer.CreateFormFile("file", filepath.Base(imagePath))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	// 添加其他字段
	_ = writer.WriteField("multi_page", "false")
	_ = writer.WriteField("render", "false")
	_ = writer.WriteField("output_format", "plain")

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", ocrURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OCR请求失败,状态码: %d", resp.StatusCode)
	}

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析JSON响应
	var ocrResponse struct {
		Result string `json:"result"`
	}
	err = json.Unmarshal(respBody, &ocrResponse)
	if err != nil {
		return "", err
	}

	return ocrResponse.Result, nil
}
