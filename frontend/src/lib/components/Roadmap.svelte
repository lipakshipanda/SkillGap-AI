<script lang="ts">
  import type { RoadmapItem } from '../types'
  import { resourceIcon, resourceBg } from '../helpers'

  export let roadmap: RoadmapItem[]
</script>

{#if !roadmap?.length}
  <p class="text-sm text-gray-400">No roadmap generated.</p>
{:else}
  <div class="space-y-2.5">
    {#each roadmap as item, i}
      <div class="flex gap-3 border border-gray-100 rounded-xl p-4 hover:border-gray-200 transition-colors">
        <div class="flex flex-col items-center gap-1.5 shrink-0">
          <div class="w-9 h-9 rounded-lg flex items-center justify-center text-base {resourceBg(item.type)}">
            {resourceIcon(item.type)}
          </div>
          <span class="text-xs font-mono text-gray-300">#{i + 1}</span>
        </div>

        <div class="flex-1 min-w-0">
          <div class="flex items-start justify-between gap-2">
            <div class="font-semibold text-sm leading-snug">{item.title}</div>
            {#if item.url}
              <a
                href={item.url}
                target="_blank"
                rel="noopener noreferrer"
                class="shrink-0 text-gray-300 hover:text-brand-600 transition-colors"
              >
                ↗
              </a>
            {/if}
          </div>
          <div class="flex items-center gap-2 mt-1 flex-wrap">
            <span class="text-xs px-2 py-0.5 rounded-md font-mono {resourceBg(item.type)}">
              {item.type}
            </span>
            {#if item.duration}
              <span class="text-xs font-mono text-gray-400">⏱ {item.duration}</span>
            {/if}
            <span
              class="text-xs font-mono px-2 py-0.5 rounded-md {item.free
                ? 'bg-emerald-50 text-emerald-700'
                : 'bg-gray-100 text-gray-500'}"
            >
              {item.free ? 'free' : 'paid'}
            </span>
          </div>
          {#if item.description}
            <p class="text-xs text-gray-500 mt-1.5 leading-relaxed">{item.description}</p>
          {/if}
        </div>
      </div>
    {/each}
  </div>
{/if}
