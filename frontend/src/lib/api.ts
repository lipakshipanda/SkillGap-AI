import axios from 'axios'
import type { AnalysisResponse, HistoryItem } from './types'

const BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const api = axios.create({ baseURL: BASE, timeout: 60000 })

interface AnalyzeTextArgs {
  resumeText: string
  jobDescription: string
  targetRole?: string
}

interface AnalyzePDFArgs {
  resumeFile: File
  jobDescription: string
  targetRole?: string
}

export async function analyzeText({ resumeText, jobDescription, targetRole = '' }: AnalyzeTextArgs) {
  const form = new FormData()
  form.append('resume_text', resumeText)
  form.append('job_description', jobDescription)
  if (targetRole) form.append('target_role', targetRole)
  const { data } = await api.post<AnalysisResponse>('/api/analyze', form)
  return data
}

export async function analyzePDF({ resumeFile, jobDescription, targetRole = '' }: AnalyzePDFArgs) {
  const form = new FormData()
  form.append('resume_file', resumeFile)
  form.append('job_description', jobDescription)
  if (targetRole) form.append('target_role', targetRole)
  const { data } = await api.post<AnalysisResponse>('/api/analyze', form)
  return data
}

export async function getHistory(limit = 10) {
  const { data } = await api.get<HistoryItem[]>('/api/history', { params: { limit } })
  return data
}

export async function getAnalysis(id: string) {
  const { data } = await api.get<AnalysisResponse>(`/api/history/${id}`)
  return data
}

export async function deleteAnalysis(id: string) {
  const { data } = await api.delete(`/api/history/${id}`)
  return data
}
