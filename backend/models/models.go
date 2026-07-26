package models

import "time"

type SkillGap struct {
	Skill         string `json:"skill"`
	CurrentLevel  int    `json:"current_level"`
	RequiredLevel int    `json:"required_level"`
	Priority      string `json:"priority"` // high | medium | low
	Reason        string `json:"reason"`
	Category      string `json:"category"`
}

type RoadmapItem struct {
	Title       string `json:"title"`
	Type        string `json:"type"` // course | project | book | article
	URL         string `json:"url"`
	Duration    string `json:"duration"`
	Free        bool   `json:"free"`
	Description string `json:"description"`
}

type MockQuestion struct {
	Question   string `json:"question"`
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
}

type AnalysisResult struct {
	MatchScore          int            `json:"match_score"`
	SkillsMatched       int            `json:"skills_matched"`
	SkillsMissing       int            `json:"skills_missing"`
	Summary             string         `json:"summary"`
	Strengths           []string       `json:"strengths"`
	EstimatedReadyWeeks int            `json:"estimated_ready_weeks"`
	Gaps                []SkillGap     `json:"gaps"`
	Roadmap             []RoadmapItem  `json:"roadmap"`
	Questions           []MockQuestion `json:"questions"`
}

// AnalysisRecord is what gets persisted in BadgerDB.
type AnalysisRecord struct {
	ID            string         `json:"id"`
	ResumeSnippet string         `json:"resume_snippet"`
	TargetRole    string         `json:"target_role"`
	Result        AnalysisResult `json:"result"`
	CreatedAt     time.Time      `json:"created_at"`
}

// AnalysisResponse is what the API returns for a single analysis.
type AnalysisResponse struct {
	ID            string         `json:"id"`
	Result        AnalysisResult `json:"result"`
	CreatedAt     time.Time      `json:"created_at"`
	ResumeSnippet string         `json:"resume_snippet"`
	TargetRole    string         `json:"target_role"`
}

// HistoryItem is the lightweight shape returned in history listings.
type HistoryItem struct {
	ID            string    `json:"id"`
	TargetRole    string    `json:"target_role"`
	MatchScore    int       `json:"match_score"`
	CreatedAt     time.Time `json:"created_at"`
	ResumeSnippet string    `json:"resume_snippet"`
}

func (r AnalysisRecord) ToResponse() AnalysisResponse {
	return AnalysisResponse{
		ID:            r.ID,
		Result:        r.Result,
		CreatedAt:     r.CreatedAt,
		ResumeSnippet: r.ResumeSnippet,
		TargetRole:    r.TargetRole,
	}
}

func (r AnalysisRecord) ToHistoryItem() HistoryItem {
	return HistoryItem{
		ID:            r.ID,
		TargetRole:    r.TargetRole,
		MatchScore:    r.Result.MatchScore,
		CreatedAt:     r.CreatedAt,
		ResumeSnippet: r.ResumeSnippet,
	}
}
