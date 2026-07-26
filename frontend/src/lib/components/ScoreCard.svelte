<script lang="ts">
  import { onMount, onDestroy, afterUpdate } from 'svelte'
  import type { AnalysisResult } from '../types'
  import { scoreColor, scoreHex } from '../helpers'

  export let result: AnalysisResult

  let chartEl: HTMLDivElement
  let chart: any = null
  let ApexChartsCtor: any = null
  let mounted = false

  async function renderChart() {
    if (!chartEl || !mounted) return

    if (!ApexChartsCtor) {
      const mod = await import('apexcharts')
      ApexChartsCtor = mod.default
    }

    if (chart) {
      chart.destroy()
      chart = null
    }

    chart = new ApexChartsCtor(chartEl, {
      chart: { type: 'radialBar', height: 140, sparkline: { enabled: true } },
      series: [result.match_score],
      colors: [scoreHex(result.match_score)],
      plotOptions: {
        radialBar: {
          hollow: { size: '60%' },
          track: { background: '#f3f4f6' },
          dataLabels: {
            name: { show: false },
            value: {
              fontSize: '20px',
              fontWeight: 800,
              offsetY: 6,
              formatter: (val: number) => `${val}%`
            }
          }
        }
      }
    })
    chart.render()
  }

  onMount(() => {
    mounted = true
    renderChart()
  })

  afterUpdate(() => {
    if (mounted) renderChart()
  })

  onDestroy(() => chart?.destroy())
</script>

<div class="space-y-4">
  <div class="grid grid-cols-3 gap-3 items-center">
    <div class="bg-gray-50 rounded-xl p-2 flex flex-col items-center">
      <div bind:this={chartEl} class="w-full"></div>
      <div class="text-xs font-mono text-gray-400 -mt-2">match score</div>
    </div>
    <div class="bg-gray-50 rounded-xl p-4 text-center">
      <div class="text-3xl font-extrabold tracking-tight text-emerald-600 flex items-center justify-center gap-1.5">
        ✓ {result.skills_matched}
      </div>
      <div class="text-xs font-mono text-gray-400 mt-1">skills matched</div>
    </div>
    <div class="bg-gray-50 rounded-xl p-4 text-center">
      <div class="text-3xl font-extrabold tracking-tight text-red-500 flex items-center justify-center gap-1.5">
        ✕ {result.skills_missing}
      </div>
      <div class="text-xs font-mono text-gray-400 mt-1">gaps to close</div>
    </div>
  </div>

  {#if result.estimated_ready_weeks > 0}
    <div class="flex items-center gap-2.5 bg-brand-50 rounded-lg px-4 py-3 border border-brand-100">
      <span class="text-brand-600 shrink-0">⏱</span>
      <span class="text-sm text-brand-800">
        <strong>{result.estimated_ready_weeks} weeks</strong> of focused study to be job-ready
      </span>
    </div>
  {/if}

  {#if result.summary}
    <p class="text-sm text-gray-600 leading-relaxed border-l-2 border-brand-200 pl-3">{result.summary}</p>
  {/if}

  {#if result.strengths?.length}
    <div>
      <p class="text-xs font-mono text-gray-400 uppercase tracking-wide mb-2">Your strengths</p>
      <div class="flex flex-wrap gap-2">
        {#each result.strengths as s}
          <span class="text-xs px-2.5 py-1 bg-emerald-50 text-emerald-700 border border-emerald-100 rounded-full font-mono">
            ✓ {s}
          </span>
        {/each}
      </div>
    </div>
  {/if}
</div>