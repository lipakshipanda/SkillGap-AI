<script lang="ts">
  import type { MockQuestion } from '../types'
  import { difficultyColor } from '../helpers'

  export let questions: MockQuestion[]

  const CATEGORY_COLORS: Record<string, string> = {
    Technical: 'bg-blue-50 text-blue-700',
    Behavioral: 'bg-purple-50 text-purple-700',
    'System Design': 'bg-amber-50 text-amber-700',
    DSA: 'bg-red-50 text-red-700'
  }
  const BORDER_COLORS = ['#4f46e5', '#0891b2', '#7c3aed', '#059669', '#dc2626', '#d97706']
</script>

{#if !questions?.length}
  <p class="text-sm text-gray-400">No questions generated.</p>
{:else}
  <div class="space-y-3">
    <p class="text-xs text-gray-400 font-mono">
      Questions tailored to your gap areas. Use STAR format for behavioral questions.
    </p>
    {#each questions as q, i}
      <div
        class="border-l-[3px] bg-gray-50 rounded-r-xl px-4 py-3 flex gap-3"
        style="border-left-color: {BORDER_COLORS[i % BORDER_COLORS.length]}"
      >
        <span class="text-gray-400 shrink-0 mt-0.5">💬</span>
        <div class="flex-1">
          <div class="flex items-center gap-2 mb-1.5 flex-wrap">
            {#if q.category}
              <span class="text-xs font-mono px-2 py-0.5 rounded-md {CATEGORY_COLORS[q.category] || 'bg-gray-100 text-gray-500'}">
                {q.category}
              </span>
            {/if}
            {#if q.difficulty}
              <span class="text-xs font-mono {difficultyColor(q.difficulty)}">
                {q.difficulty}
              </span>
            {/if}
          </div>
          <p class="text-sm text-gray-700 leading-relaxed">{q.question}</p>
        </div>
      </div>
    {/each}
  </div>
{/if}
