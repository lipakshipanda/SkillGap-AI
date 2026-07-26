<script lang="ts">
  import type { SkillGap } from '../types'
  import { priorityBadgeClass, priorityLabel, priorityColor } from '../helpers'

  export let gaps: SkillGap[]
</script>

{#if !gaps?.length}
  <p class="text-sm text-gray-400">No skill gaps identified.</p>
{:else}
  <div class="space-y-2.5">
    {#each gaps as gap}
      <div class="border border-gray-100 rounded-xl p-4 hover:border-gray-200 transition-colors">
        <div class="flex items-start justify-between gap-3 mb-3">
          <div>
            <div class="font-semibold text-sm">{gap.skill}</div>
            {#if gap.category}
              <span class="text-xs font-mono text-gray-400">{gap.category}</span>
            {/if}
          </div>
          <span class="{priorityBadgeClass(gap.priority)} shrink-0">
            {priorityLabel(gap.priority)}
          </span>
        </div>

        <div class="space-y-1.5 mb-2.5">
          <div class="flex items-center gap-3">
            <span class="text-xs font-mono text-gray-400 w-16 shrink-0">current</span>
            <div class="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-700"
                style="width: {gap.current_level}%; background: {priorityColor(gap.priority)}"
              ></div>
            </div>
            <span class="text-xs font-mono text-gray-500 w-8 text-right">{gap.current_level}%</span>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-xs font-mono text-gray-400 w-16 shrink-0">required</span>
            <div class="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-gray-300" style="width: {gap.required_level}%"></div>
            </div>
            <span class="text-xs font-mono text-gray-500 w-8 text-right">{gap.required_level}%</span>
          </div>
        </div>

        <p class="text-xs text-gray-500 leading-relaxed">{gap.reason}</p>
      </div>
    {/each}
  </div>
{/if}
