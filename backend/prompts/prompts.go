package prompts

import "fmt"

const SystemPrompt = `You are SkillGap AI, a career intelligence engine for analyzing student and early-career developer profiles against job descriptions.

Your analysis is precise, actionable, and prioritized. You always respond with valid JSON only — no markdown fences, no preamble, no explanation outside the JSON.`

func BuildAnalysisPrompt(resumeText, jobDescription, targetRole string) string {
	roleHint := ""
	if targetRole != "" {
		roleHint = fmt.Sprintf("\nTarget role context: %s", targetRole)
	}

	return fmt.Sprintf(`Analyze this resume against the job description and return a structured skill gap report.%s

--- RESUME ---
%s

--- JOB DESCRIPTION ---
%s

Return ONLY this JSON (no markdown, no backticks):
{
  "match_score": <integer 0-100>,
  "skills_matched": <integer>,
  "skills_missing": <integer>,
  "summary": "<2-sentence honest summary>",
  "strengths": ["<strength 1>", "<strength 2>", "<strength 3>"],
  "estimated_ready_weeks": <integer>,
  "gaps": [
    {
      "skill": "<skill name>",
      "current_level": <0-100>,
      "required_level": <0-100>,
      "priority": "high" | "medium" | "low",
      "reason": "<one sentence why this matters>",
      "category": "<Backend | Frontend | DevOps | Data | Soft Skills>"
    }
  ],
  "roadmap": [
    {
      "title": "<specific course/project name>",
      "type": "course" | "project" | "book" | "article",
      "url": "<real URL or null>",
      "duration": "<e.g. 12 hrs>",
      "free": true | false,
      "description": "<one sentence what they will gain>"
    }
  ],
  "questions": [
    {
      "question": "<interview question>",
      "category": "<Behavioral | Technical | System Design | DSA>",
      "difficulty": "easy" | "medium" | "hard"
    }
  ]
}

Rules:
- gaps: 5-7 items ordered by priority (high first)
- roadmap: 6-8 items ordered by what to tackle first, mix types
- questions: exactly 6 questions tailored to this candidate's gaps
- URLs must be real well-known resources (freeCodeCamp, MDN, roadmap.sh, Coursera) or null`, roleHint, resumeText, jobDescription)
}
