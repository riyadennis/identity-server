package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	openaiURL = "https://api.groq.com/openai/v1/chat/completions"
	githubURL = "https://api.github.com"

	openaiModel  = "llama-3.3-70b-versatile"
	maxDiffLines = 500
	maxTokens    = 1024
	httpTimeout  = 30 * time.Second
)

var httpClient = &http.Client{Timeout: httpTimeout}

type openaiRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []openaiMessage `json:"messages"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiResponse struct {
	Choices []struct {
		Message openaiMessage `json:"message"`
	} `json:"choices"`
}

type githubCommentRequest struct {
	Body string `json:"body"`
}

func main() {
	apiKey := mustEnv("OPENAI_API_KEY")
	githubToken := mustEnv("GITHUB_TOKEN")
	prNumber := mustEnv("PR_NUMBER")
	repo := mustEnv("REPO")

	diff, err := os.ReadFile("/tmp/pr_diff.txt")
	if err != nil {
		log.Fatalf("failed to read diff: %v", err)
	}
	if len(bytes.TrimSpace(diff)) == 0 {
		log.Println("diff is empty, nothing to review")
		return
	}

	review, err := callOpenAI(apiKey, truncateDiff(string(diff)))
	if err != nil {
		log.Fatalf("openai review failed: %v", err)
	}

	if err := postGitHubComment(githubToken, repo, prNumber, review); err != nil {
		log.Fatalf("failed to post github comment: %v", err)
	}

	log.Println("review posted successfully")
}

func truncateDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	if len(lines) <= maxDiffLines {
		return diff
	}
	log.Printf("diff truncated from %d to %d lines to stay within API limits", len(lines), maxDiffLines)
	return strings.Join(lines[:maxDiffLines], "\n") + "\n\n[diff truncated — showing first 500 lines only]"
}

func callOpenAI(apiKey, diff string) (string, error) {
	prompt := fmt.Sprintf(`You are a senior software engineer performing a code review.
Review the following git diff and provide concise, actionable feedback.
Focus on: bugs, security issues, performance problems, and code clarity.
Format your response in markdown.

Diff:
%s`, diff)

	body, err := json.Marshal(openaiRequest{
		Model:     openaiModel,
		MaxTokens: maxTokens,
		Messages: []openaiMessage{
			{Role: "system", Content: "You are a senior software engineer performing code reviews."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	const maxRetries = 3
	retryDelays := []time.Duration{10 * time.Second, 30 * time.Second, 60 * time.Second}

	var resp *http.Response
	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest(http.MethodPost, openaiURL, bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err = httpClient.Do(req)
		if err != nil {
			return "", err
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxRetries {
			log.Printf("rate limited (429): %s — retrying in %s (attempt %d/%d)",
				string(respBody), retryDelays[attempt], attempt+1, maxRetries)
			time.Sleep(retryDelays[attempt])
			continue
		}

		return "", fmt.Errorf("openai API returned status %d: %s", resp.StatusCode, string(respBody))
	}
	if resp == nil {
		return "", fmt.Errorf("no response received after %d attempts", maxRetries)
	}
	defer resp.Body.Close()

	var result openaiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai returned no content")
	}

	return result.Choices[0].Message.Content, nil
}

func postGitHubComment(token, repo, prNumber, body string) error {
	url := fmt.Sprintf("%s/repos/%s/issues/%s/comments", githubURL, repo, prNumber)

	payload, err := json.Marshal(githubCommentRequest{Body: "## OpenAI Code Review\n\n" + body})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("github API returned status %d", resp.StatusCode)
	}
	return nil
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}