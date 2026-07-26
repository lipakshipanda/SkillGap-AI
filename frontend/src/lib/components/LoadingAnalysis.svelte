<script lang="ts">
  const STEPS = [
    { label: 'Parsing your resume…', delay: 0 },
    { label: 'Analyzing job requirements…', delay: 800 },
    { label: 'Identifying skill gaps…', delay: 2000 },
    { label: 'Building your roadmap…', delay: 3500 },
    { label: 'Generating interview questions…', delay: 5000 }
  ]
</script>

<div class="flex flex-col items-center justify-center py-16 gap-6">
  <div class="relative w-12 h-12">
    <div class="w-12 h-12 border-4 border-gray-100 rounded-full"></div>
    <div class="w-12 h-12 border-4 border-brand-600 border-t-transparent rounded-full animate-spin absolute inset-0"></div>
  </div>

  <div class="text-center">
    <p class="font-semibold text-sm">Analyzing your profile</p>
    <p class="text-xs text-gray-400 font-mono mt-1">Mistral AI is reading your resume + JD…</p>
  </div>

  <div class="w-full max-w-xs space-y-2">
    {#each STEPS as step}
      <div
        class="flex items-center gap-2.5 text-xs font-mono text-gray-400"
        style="opacity: 0; animation: fadeIn 0.4s ease-out {step.delay}ms forwards"
      >
        <div class="w-1.5 h-1.5 rounded-full bg-brand-400 shrink-0"></div>
        {step.label}
      </div>
    {/each}
  </div>
</div>

<style>
  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(4px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
