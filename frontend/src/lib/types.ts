export type Priority = 'high' | 'medium' | 'low'
export type ResourceType = 'course' | 'project' | 'book' | 'article'

export interface SkillGap {
  skill: string
  current_level: number
  required_level: number
  priority: Priority
  reason: string
  category: string
}

export interface RoadmapItem {
  title: string
  type: ResourceType
  url?: string | null
  duration: string
  free: boolean
  description: string
}

export interface MockQuestion {
  question: string
  category: string
  difficulty: string
}

export interface AnalysisResult {
  match_score: number
  skills_matched: number
  skills_missing: number
  summary: string
  strengths: string[]
  estimated_ready_weeks: number
  gaps: SkillGap[]
  roadmap: RoadmapItem[]
  questions: MockQuestion[]
}

export interface AnalysisResponse {
  id: string
  result: AnalysisResult
  created_at: string
  resume_snippet: string
  target_role: string
}

export interface HistoryItem {
  id: string
  target_role: string
  match_score: number
  created_at: string
  resume_snippet: string
}
