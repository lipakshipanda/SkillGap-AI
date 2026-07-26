<script lang="ts">
  import { onMount } from 'svelte'
  import { goto } from '$app/navigation'
  import { getHistory, deleteAnalysis } from '$lib/api'
  import { scoreColor, timeAgo } from '$lib/helpers'
  import type { HistoryItem } from '$lib/types'

  let items: HistoryItem[] = []
  let loading = true

  onMount(async () => {
    try {
      items = await getHistory(20)
    } catch {
      // could surface a toast here
    } finally {
      loading = false
    }
  })

  async function handleDelete(id: string, e: MouseEvent) {
    e.stopPropagation()
    try {
      await deleteAnalysis(id)
      items = items.filter((i) => i.id !== id)
    } catch {
      // ignore
    }
  }
</script>

{#if loading}
  <div class="flex items-center justify-center py-24">
    <div class="w-6 h-6 border-2 border-gray-200 border-t-brand-600 rounded-full animate-spin"></div>
  </div>
{:else}
  <div class="max-w-3xl mx-auto px-6 py-8">
    <h1 class="text-2xl font-extrabold tracking-tight mb-1">Analysis History</h1>
    <p class="text-sm text-gray-400 mb-6">Your past resume analyses</p>

    {#if items.length === 0}
      <div class="card text-center py-16">
        <p class="text-3xl mb-3">📋</p>
        <p class="text-sm text-gray-500">No analyses yet.</p>
        <button class="btn-primary mt-4 mx-auto" on:click={() => goto('/')}> Run your first analysis </button>
      </div>
    {:else}
      <div class="space-y-3">
        {#each items as item (item.id)}
          <div
            role="button"
            tabindex="0"
            on:click={() => goto('/')}
            on:keydown={(e) => e.key === 'Enter' && goto('/')}
            class="card flex items-center gap-4 cursor-pointer hover:border-gray-200 transition-all group"
          >
            <div class="text-2xl font-extrabold tracking-tight w-14 shrink-0 {scoreColor(item.match_score)}">
              {item.match_score}%
            </div>
            <div class="flex-1 min-w-0">
              <div class="font-semibold text-sm truncate">{item.target_role || 'Untitled'}</div>
              <div class="text-xs text-gray-400 font-mono mt-0.5">{item.resume_snippet}</div>
              <div class="text-xs text-gray-300 font-mono mt-0.5">{timeAgo(item.created_at)}</div>
            </div>
            <div class="flex items-center gap-2">
              <button
                on:click={(e) => handleDelete(item.id, e)}
                class="opacity-0 group-hover:opacity-100 transition-opacity p-1.5 rounded-lg hover:bg-red-50 text-gray-400 hover:text-red-500"
              >
                🗑
              </button>
              <span class="text-gray-300 group-hover:text-brand-500 transition-colors">→</span>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
{/if}
