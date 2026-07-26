<script lang="ts">
  import { page } from '$app/stores'
  $: pathname = $page.url.pathname
</script>

<nav class="border-b border-gray-100 bg-white sticky top-0 z-50">
  <div class="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
    <a href="/" class="flex items-center gap-2.5 font-extrabold text-base tracking-tight">
      <div class="w-7 h-7 bg-brand-900 rounded-lg flex items-center justify-center">
        <span class="text-white text-xs">⚡</span>
      </div>
      SkillGap AI
    </a>

    <div class="flex items-center gap-1">
      <a
        href="/"
        class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors {pathname === '/'
          ? 'bg-brand-50 text-brand-600'
          : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'}"
      >
        Analyze
      </a>
      <a
        href="/history"
        class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors flex items-center gap-1.5 {pathname ===
        '/history'
          ? 'bg-brand-50 text-brand-600'
          : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'}"
      >
        🕘 History
      </a>
    </div>
  </div>
</nav>
