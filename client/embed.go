package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

type EmbeddingRequest struct {
	Prompt string `json:"prompt"`
}

// get the embeddings given a piece of text
func getEmbeddings(text string) ([]float32, error) {
	req := &EmbeddingRequest{
		Prompt: text,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("err in marshaling:", err)
		return []float32{}, err
	}

	r := bytes.NewReader(reqJson)
	httpResp, err := http.Post("http://localhost:11333/api/embeddings",
		"application/json", r)
	if err != nil {
		fmt.Println("err in calling embedding API server:", err)
		return []float32{}, err
	}
	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println("err in reading:", err)
		return []float32{}, err
	}

	resp := &EmbeddingResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		fmt.Println("err in unmarshaling:", err)
		return []float32{}, err
	}

	return resp.Embedding, nil
}

// cosine similarity of 2 float64 slices
func similarity(a, b []float64) float64 {
	return dotproduct(a, b) / (magnitude(a) * magnitude(b))
}

// dot product of 2 float64 slices
func dotproduct(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	var dp float64
	for i := 0; i < len(a); i++ {
		dp += a[i] * b[i]
	}
	return dp
}

// magnitude of a float64 slice
func magnitude(a []float64) float64 {
	var mag float64
	for i := 0; i < len(a); i++ {
		mag += math.Pow(a[i], 2.0)
	}
	return mag
}

func test() {
	cvfile, _ := os.Open("cvs/dev01.txt")
	cv, _ := io.ReadAll(cvfile)
	cvEmbeddings, _ := getEmbeddings(string(cv))
	float64CV := toF64(cvEmbeddings)

	jobfile, _ := os.Open("jobs/dev04.txt")
	job, _ := io.ReadAll(jobfile)
	jobEmbeddings, _ := getEmbeddings(string(job))
	float64Job := make([]float64, len(jobEmbeddings))
	for i, val := range jobEmbeddings {
		float64Job[i] = float64(val)
	}

	sim := similarity(float64CV, float64Job)
	fmt.Println(sim)
}

func toF64(cvEmbeddings []float32) []float64 {
	float64CV := make([]float64, len(cvEmbeddings))
	for i, val := range cvEmbeddings {
		float64CV[i] = float64(val)
	}
	return float64CV
}
