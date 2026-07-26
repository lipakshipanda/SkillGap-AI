<script lang="ts">
  import { onMount } from 'svelte'
  import { getHistory } from '../api'
  import { timeAgo, scoreColor } from '../helpers'
  import type { HistoryItem } from '../types'

  export let onSelect: (id: string) => void

  let items: HistoryItem[] = []

  onMount(async () => {
    try {
      items = await getHistory(5)
    } catch {
      items = []
    }
  })
</script>

{#if items.length}
  <div class="border border-gray-100 rounded-xl p-4">
    <div class="flex items-center gap-2 mb-3">
      <span class="text-gray-400">🕘</span>
      <span class="text-xs font-mono text-gray-400 uppercase tracking-wide">Recent analyses</span>
    </div>
    <div class="space-y-1">
      {#each items as item (item.id)}
        <button
          on:click={() => onSelect(item.id)}
          class="w-full text-left flex items-center gap-3 p-2.5 rounded-lg hover:bg-gray-50 transition-colors group"
        >
          <div class="text-base font-extrabold tracking-tight w-10 shrink-0 {scoreColor(item.match_score)}">
            {item.match_score}%
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-xs font-medium truncate text-gray-700">{item.target_role || 'Untitled'}</div>
            <div class="text-xs font-mono text-gray-400">{timeAgo(item.created_at)}</div>
          </div>
        </button>
      {/each}
    </div>
  </div>
{/if}
