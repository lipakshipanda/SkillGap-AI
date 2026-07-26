<script lang="ts">
  import ResumeInput from '$lib/components/ResumeInput.svelte'
  import JDInput from '$lib/components/JDInput.svelte'
  import ScoreCard from '$lib/components/ScoreCard.svelte'
  import SkillGapList from '$lib/components/SkillGapList.svelte'
  import Roadmap from '$lib/components/Roadmap.svelte'
  import MockInterview from '$lib/components/MockInterview.svelte'
  import LoadingAnalysis from '$lib/components/LoadingAnalysis.svelte'
  import HistorySidebar from '$lib/components/HistorySidebar.svelte'
  import { analyzeText, analyzePDF, getAnalysis } from '$lib/api'
  import { DEMO_RESUME, DEMO_JD } from '$lib/helpers'
  import type { AnalysisResponse } from '$lib/types'

  const TABS = [
    { key: 'skills', label: 'Skill Gaps' },
    { key: 'roadmap', label: 'Roadmap' },
    { key: 'interview', label: 'Mock Interview' }
  ] as const

  let resumeText = ''
  let resumeFile: File | null = null
  let resumeMode: 'text' | 'pdf' = 'text'
  let jd = ''
  let targetRole = ''

  let result: AnalysisResponse | null = null
  let loading = false
  let error: string | null = null
  let activeTab: (typeof TABS)[number]['key'] = 'skills'

  $: canSubmit = (resumeText.trim().length > 50 || resumeFile) && jd.trim().length > 50

  async function handleSubmit() {
    loading = true
    error = null
    result = null
    try {
      result = resumeFile
        ? await analyzePDF({ resumeFile, jobDescription: jd, targetRole })
        : await analyzeText({ resumeText, jobDescription: jd, targetRole })
      activeTab = 'skills'
    } catch (err: any) {
      error = err?.response?.data?.detail || err?.message || 'Analysis failed'
    } finally {
      loading = false
    }
  }

  function loadDemo() {
    resumeText = DEMO_RESUME
    jd = DEMO_JD
    targetRole = 'Software Engineer'
    resumeMode = 'text'
  }

  function reset() {
    result = null
    error = null
    activeTab = 'skills'
  }

  async function handleHistorySelect(id: string) {
    try {
      result = await getAnalysis(id)
      activeTab = 'skills'
    } catch {
      // ignore
    }
  }
</script>

<div class="max-w-6xl mx-auto px-6 py-8">
  <div class="mb-8 max-w-xl">
    <p class="text-xs font-mono text-gray-400 uppercase tracking-widest mb-2">career intelligence platform</p>
    <h1 class="text-4xl font-extrabold tracking-tight leading-tight mb-2">
      Close the gap.<br />
      <span class="text-brand-600">Get the role.</span>
    </h1>
    <p class="text-sm text-gray-500 leading-relaxed">
      Paste your resume + a job description. AI builds a personalized skill roadmap with courses, GitHub projects,
      and mock interview questions.
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-[420px_1fr] gap-6 items-start">
    <!-- Left panel -->
    <div class="space-y-4">
      <div class="card">
        <div class="flex items-center gap-2 mb-4 text-xs font-mono text-gray-400 uppercase tracking-wide">
          <div class="w-5 h-5 rounded-full bg-brand-900 text-white flex items-center justify-center text-[10px]">1</div>
          Input
        </div>

        <div class="space-y-4">
          <ResumeInput
            value={resumeText}
            onChange={(v) => (resumeText = v)}
            onFileChange={(f) => (resumeFile = f)}
            file={resumeFile}
            mode={resumeMode}
            setMode={(m) => (resumeMode = m)}
          />
          <JDInput value={jd} onChange={(v) => (jd = v)} targetRole={targetRole} onRoleChange={(v) => (targetRole = v)} />
        </div>

        <div class="flex gap-2 mt-4">
          <button class="btn-primary flex-1" on:click={handleSubmit} disabled={!canSubmit || loading}>
            ⚙ {loading ? 'Analyzing…' : 'Analyze my profile'}
          </button>
          <button class="btn-secondary" on:click={loadDemo}> ✨ Demo </button>
        </div>

        {#if error}
          <div class="mt-3 text-xs text-red-600 bg-red-50 border border-red-100 rounded-lg p-3 font-mono">
            {error}
          </div>
        {/if}
      </div>

      <HistorySidebar onSelect={handleHistorySelect} />
    </div>

    <!-- Right panel -->
    <div class="card min-h-[400px]">
      {#if loading}
        <LoadingAnalysis />
      {:else if result}
        <div>
          <div class="flex items-center gap-2 mb-5 text-xs font-mono text-gray-400 uppercase tracking-wide">
            <div class="w-5 h-5 rounded-full bg-brand-900 text-white flex items-center justify-center text-[10px]">2</div>
            Results
            {#if result.target_role}
              <span class="ml-auto text-gray-300 normal-case">{result.target_role}</span>
            {/if}
          </div>

          <ScoreCard result={result.result} />

          <div class="flex border-b border-gray-100 mt-6 mb-4">
            {#each TABS as tab}
              <button
                on:click={() => (activeTab = tab.key)}
                class="px-4 py-2.5 text-xs font-mono uppercase tracking-wide border-b-2 transition-colors {activeTab ===
                tab.key
                  ? 'border-brand-600 text-brand-600'
                  : 'border-transparent text-gray-400 hover:text-gray-700'}"
              >
                {tab.label}
              </button>
            {/each}
          </div>

          {#if activeTab === 'skills'}
            <SkillGapList gaps={result.result.gaps} />
          {:else if activeTab === 'roadmap'}
            <Roadmap roadmap={result.result.roadmap} />
          {:else if activeTab === 'interview'}
            <MockInterview questions={result.result.questions} />
          {/if}

          <button on:click={reset} class="btn-secondary mt-6 w-full justify-center text-xs"> ← Analyze another role </button>
        </div>
      {:else}
        <div class="flex flex-col items-center justify-center py-16 text-center gap-3">
          <div class="text-4xl">🎯</div>
          <p class="font-semibold text-sm text-gray-700">Your analysis will appear here</p>
          <p class="text-xs text-gray-400 max-w-xs leading-relaxed">
            Paste your resume and a job description, then hit <strong>Analyze</strong>.
          </p>
        </div>
      {/if}
    </div>
  </div>
</div>
