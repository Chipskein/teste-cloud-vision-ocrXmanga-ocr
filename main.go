package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func convertImageToBase64(imgPath string) string {
	data, err := os.ReadFile(imgPath)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

type Request struct {
	Requests []struct {
		Image struct {
			Content string `json:"content"`
		} `json:"image"`
		Features []struct {
			Type string `json:"type"`
		} `json:"features"`
	} `json:"requests"`
}

func callGCVAPI(apiKey string, base64Image string, outPath string) error {
	url := "https://vision.googleapis.com/v1/images:annotate?key=" + apiKey
	reqBody := map[string]any{
		"requests": []any{
			map[string]any{
				"image": map[string]any{
					"content": base64Image,
				},
				"features": []any{
					map[string]any{
						"type": "TEXT_DETECTION",
					},
				},
				"imageContext": map[string]any{
					"languageHints": []string{"ja"},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(outPath, body, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Status:", resp.Status)
	fmt.Printf("Response saved to %s\n", outPath)

	return nil
}

func main() {
	apiKeyFlagPtr := flag.String("GCV_API_KEY", "", "GOOGLE CLOUD VISION API")
	flag.Parse()
	if apiKeyFlagPtr == nil {
		log.Fatalln("Use GCV_API_KEY flag to set GOOGLE CLOUD VISION API")
	}
	apiKey := *apiKeyFlagPtr
	pathImgs := "img-teste/bubble-texts"
	outTextPath := "out-text/gcv"
	log.Printf("PathImgs: %s\n", pathImgs)
	log.Printf("outTextPath: %s\n", outTextPath)
	log.Printf("API KEY: %s\n", apiKey)
	entries, err := os.ReadDir(pathImgs)
	if err != nil {
		log.Fatalf("Couldn't read pathWithImgs: %s\n", pathImgs)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		pathStr := fmt.Sprintf("%s/%s", pathImgs, entry.Name())
		log.Printf("Reading File: %s\n", pathStr)
		imgBase64 := convertImageToBase64(pathStr)
		outJsonPath := fmt.Sprintf("%s/%s.json", outTextPath, entry.Name())
		err := callGCVAPI(apiKey, imgBase64, outJsonPath)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

}
