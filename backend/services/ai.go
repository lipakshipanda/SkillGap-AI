package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"skillgap-ai/models"
	"skillgap-ai/prompts"
)

const mistralURL = "https://api.mistral.ai/v1/chat/completions"

type mistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type mistralRequest struct {
	Model          string            `json:"model"`
	Messages       []mistralMessage  `json:"messages"`
	Temperature    float64           `json:"temperature"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

type mistralChoice struct {
	Message mistralMessage `json:"message"`
}

type mistralResponse struct {
	Choices []mistralChoice `json:"choices"`
}

// rawGap / rawRoadmap / rawQuestion mirror the loose JSON the model returns
// before we validate/clamp it into strongly typed models, same role as
// _parse_result() in the old Python analyzer.py.
type rawResult struct {
	MatchScore          int      `json:"match_score"`
	SkillsMatched       int      `json:"skills_matched"`
	SkillsMissing       int      `json:"skills_missing"`
	Summary             string   `json:"summary"`
	Strengths           []string `json:"strengths"`
	EstimatedReadyWeeks int      `json:"estimated_ready_weeks"`
	Gaps                []struct {
		Skill         string `json:"skill"`
		CurrentLevel  int    `json:"current_level"`
		RequiredLevel int    `json:"required_level"`
		Priority      string `json:"priority"`
		Reason        string `json:"reason"`
		Category      string `json:"category"`
	} `json:"gaps"`
	Roadmap []struct {
		Title       string `json:"title"`
		Type        string `json:"type"`
		URL         string `json:"url"`
		Duration    string `json:"duration"`
		Free        bool   `json:"free"`
		Description string `json:"description"`
	} `json:"roadmap"`
	Questions []struct {
		Question   string `json:"question"`
		Category   string `json:"category"`
		Difficulty string `json:"difficulty"`
	} `json:"questions"`
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func safePriority(p string) string {
	switch strings.ToLower(p) {
	case "high", "medium", "low":
		return strings.ToLower(p)
	default:
		return "medium"
	}
}

func safeResourceType(t string) string {
	switch strings.ToLower(t) {
	case "course", "project", "book", "article":
		return strings.ToLower(t)
	default:
		return "course"
	}
}

// AnalyzeResume calls the Mistral chat completions API and parses the
// response into a strongly-typed AnalysisResult, equivalent to the old
// Gemini-backed analyze_resume() in analyzer.py.
func AnalyzeResume(resumeText, jobDescription, targetRole string) (*models.AnalysisResult, error) {
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MISTRAL_API_KEY is not set")
	}

	prompt := prompts.SystemPrompt + "\n\n" + prompts.BuildAnalysisPrompt(resumeText, jobDescription, targetRole)

	reqBody := mistralRequest{
		Model:       "mistral-large-latest",
		Temperature: 0.3,
		Messages: []mistralMessage{
			{Role: "user", Content: prompt},
		},
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(http.MethodPost, mistralURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mistral request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mistral returned %d: %s", resp.StatusCode, string(body))
	}

	var mr mistralResponse
	if err := json.Unmarshal(body, &mr); err != nil {
		return nil, fmt.Errorf("failed to decode mistral envelope: %w", err)
	}
	if len(mr.Choices) == 0 {
		return nil, fmt.Errorf("mistral returned no choices")
	}

	raw := strings.TrimSpace(mr.Choices[0].Message.Content)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var rr rawResult
	if err := json.Unmarshal([]byte(raw), &rr); err != nil {
		return nil, fmt.Errorf("mistral returned invalid JSON: %w", err)
	}

	result := &models.AnalysisResult{
		MatchScore:          clamp(rr.MatchScore, 0, 100),
		SkillsMatched:       rr.SkillsMatched,
		SkillsMissing:       rr.SkillsMissing,
		Summary:             rr.Summary,
		Strengths:           rr.Strengths,
		EstimatedReadyWeeks: rr.EstimatedReadyWeeks,
	}
	if result.Strengths == nil {
		result.Strengths = []string{}
	}

	for _, g := range rr.Gaps {
		result.Gaps = append(result.Gaps, models.SkillGap{
			Skill:         g.Skill,
			CurrentLevel:  clamp(g.CurrentLevel, 0, 100),
			RequiredLevel: clamp(g.RequiredLevel, 0, 100),
			Priority:      safePriority(g.Priority),
			Reason:        g.Reason,
			Category:      g.Category,
		})
	}
	for _, r := range rr.Roadmap {
		result.Roadmap = append(result.Roadmap, models.RoadmapItem{
			Title:       r.Title,
			Type:        safeResourceType(r.Type),
			URL:         r.URL,
			Duration:    r.Duration,
			Free:        r.Free,
			Description: r.Description,
		})
	}
	for _, q := range rr.Questions {
		result.Questions = append(result.Questions, models.MockQuestion{
			Question:   q.Question,
			Category:   q.Category,
			Difficulty: q.Difficulty,
		})
	}

	return result, nil
}
