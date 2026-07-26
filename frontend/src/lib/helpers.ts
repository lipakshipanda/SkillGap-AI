export function priorityColor(priority: string): string {
  if (priority === 'high') return '#ef4444'
  if (priority === 'medium') return '#f59e0b'
  return '#10b981'
}

export function priorityBadgeClass(priority: string): string {
  if (priority === 'high') return 'badge-high'
  if (priority === 'medium') return 'badge-medium'
  return 'badge-low'
}

export function priorityLabel(priority: string): string {
  if (priority === 'high') return '🔴 High'
  if (priority === 'medium') return '🟡 Medium'
  return '🟢 Low'
}

export function scoreColor(score: number): string {
  if (score >= 70) return 'text-emerald-600'
  if (score >= 45) return 'text-amber-600'
  return 'text-red-500'
}

export function scoreHex(score: number): string {
  if (score >= 70) return '#059669'
  if (score >= 45) return '#d97706'
  return '#ef4444'
}

export function resourceIcon(type: string): string {
  const map: Record<string, string> = { course: '🎓', project: '💻', book: '📖', article: '📄' }
  return map[type] || '🔗'
}

export function resourceBg(type: string): string {
  const map: Record<string, string> = {
    course: 'bg-violet-50 text-violet-700',
    project: 'bg-emerald-50 text-emerald-700',
    book: 'bg-blue-50 text-blue-700',
    article: 'bg-gray-100 text-gray-600'
  }
  return map[type] || 'bg-gray-100 text-gray-600'
}

export function difficultyColor(d: string): string {
  if (d === 'hard') return 'text-red-600'
  if (d === 'medium') return 'text-amber-600'
  return 'text-emerald-600'
}

export function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
}

export const DEMO_RESUME = `B.Tech Computer Science, 3rd Year — NIT Rourkela
CGPA: 8.4/10

Skills: Python, React (basic), HTML/CSS, MySQL, Git
Projects:
  • Weather App — React + OpenWeatherMap API
  • College Portal — PHP, MySQL CRUD app
  • DSA practice repo on GitHub (150+ problems)

Experience:
  • 2-month web dev intern — built 3 internal dashboards using React`

export const DEMO_JD = `Senior Software Engineer – Full Stack
Company: Razorpay (Bangalore)

Requirements:
- 3+ years React, Node.js / TypeScript
- PostgreSQL, Redis — query optimization
- REST & GraphQL API design
- Docker, Kubernetes, AWS (EC2, S3, Lambda)
- System design — distributed systems, CAP theorem
- CI/CD pipelines (GitHub Actions / Jenkins)
- Strong DSA fundamentals`
