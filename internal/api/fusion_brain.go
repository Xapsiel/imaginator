package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"imageBot/internal/config"
)

type Text2ImageAPI struct {
	URL    string
	APIKey string
	Secret string
}

func New(cfg config.FB) *Text2ImageAPI {
	return &Text2ImageAPI{
		URL:    cfg.URL,
		APIKey: cfg.APIKey,
		Secret: cfg.Secret,
	}
}

type ModelResponse struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Version json.Number `json:"version"`
	Type    string      `json:"type"`
}
type GenerateResponse struct {
	UUID string `json:"uuid"`
}

type StatusResponse struct {
	Status string `json:"status"`
	//Images  []string `json:"images"`
	Error  string       `json:"errorDescription"`
	Result StatusResult `json:"result"`
	//Message string   `json:"message"`
}
type StatusResult struct {
	Files    []string `json:"files"`
	Censored bool     `json:"censored"`
}

func (api *Text2ImageAPI) GetModel() (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", api.URL+"key/api/v1/pipelines", nil)
	if err != nil {
		return "", fmt.Errorf("request creation failed: %v", err)
	}

	req.Header.Add("X-Key", "Key "+api.APIKey)
	req.Header.Add("X-Secret", "Secret "+api.Secret)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error [%d]: %s", resp.StatusCode, string(body))
	}

	var modelResponse []ModelResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("response reading failed: %v", err)
	}

	if err := json.Unmarshal(body, &modelResponse); err != nil {
		return "", fmt.Errorf("JSON decode failed: %v\nResponse: %s", err, string(body))
	}

	if len(modelResponse) == 0 {
		return "", fmt.Errorf("no available models")
	}

	return modelResponse[0].ID, nil

}

func (api *Text2ImageAPI) Generate(prompt, model string, images, width, height int) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormField("pipeline_id")
	if err != nil {
		return "", fmt.Errorf("form creation failed: %v", err)
	}
	part.Write([]byte(model))

	params := map[string]interface{}{
		"type":      "GENERATE",
		"numImages": images,
		"width":     width,
		"height":    height,
		"generateParams": map[string]interface{}{
			"query": prompt,
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("params marshaling failed: %v", err)
	}

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Type", "application/json")
	partHeader.Set("Content-Disposition", `form-data; name="params"`)

	part, err = writer.CreatePart(partHeader)
	if err != nil {
		return "", fmt.Errorf("params part creation failed: %v", err)
	}
	part.Write(paramsJSON)

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("form closing failed: %v", err)
	}

	req, err := http.NewRequest("POST", api.URL+"key/api/v1/pipeline/run", body)
	if err != nil {
		return "", fmt.Errorf("request creation failed: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-Key", "Key "+api.APIKey)
	req.Header.Add("X-Secret", "Secret "+api.Secret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	var generateResponse GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&generateResponse); err != nil {
		return "", fmt.Errorf("response decode failed: %v", err)
	}

	return generateResponse.UUID, nil
}
func (api *Text2ImageAPI) CheckGeneration(requestID string, attempts int, delay time.Duration) ([]string, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	for i := 0; i < attempts; i++ {
		req, err := http.NewRequest("GET", api.URL+"key/api/v1/pipeline/status/"+requestID, nil)
		if err != nil {
			return nil, fmt.Errorf("request creation failed: %v", err)
		}

		req.Header.Add("X-Key", "Key "+api.APIKey)
		req.Header.Add("X-Secret", "Secret "+api.Secret)

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("API error [%d]: %s", resp.StatusCode, string(body))
		}

		var statusResponse StatusResponse
		if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("response decode failed: %v", err)
		}
		resp.Body.Close()

		switch statusResponse.Status {
		case "DONE":
			return statusResponse.Result.Files, nil
		case "FAIL":
			return nil, fmt.Errorf("generation failed: %s", statusResponse.Error)
		}

		time.Sleep(delay)
	}

	return nil, fmt.Errorf("maximum attempts reached (%d)", attempts)
}

func (api *Text2ImageAPI) DecodeImage(images []string) ([]byte, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("no images to save")
	}

	imageData, err := base64.StdEncoding.DecodeString(images[0])
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %v", err)
	}

	return imageData, nil
}

func (api *Text2ImageAPI) Draw(prompt string, width, height int) ([]byte, error) {
	modelID, err := api.GetModel()
	if err != nil {
		slog.Error(fmt.Sprintf("🚨 Model fetch error: %v", err))
		return nil, err
	}

	uuid, err := api.Generate(prompt, modelID, 1, width, height)
	if err != nil {
		slog.Info(fmt.Sprintf("🚨 Generation error: %v", err))
		return nil, err
	}

	images, err := api.CheckGeneration(uuid, 15, 10*time.Second)
	if err != nil {
		slog.Error(fmt.Sprintf("🚨 Status check error: %v", err))
		return nil, err
	}

	data, err := api.DecodeImage(images)
	if err != nil {
		slog.Error(fmt.Sprintf("Decoding image error: %v", err))
	}
	//slog.Error(fmt.Sprintf("🚨 Image save error: %v", err))
	//	return nil, err
	//}
	///
	//slog.Info("✅ Image successfully generated and saved as generated_image.jpg")
	return data, nil
}
