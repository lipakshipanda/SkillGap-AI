<script lang="ts">
  export let value: string
  export let onChange: (v: string) => void
  export let targetRole: string
  export let onRoleChange: (v: string) => void

  function handleRoleInput(e: Event) {
    onRoleChange((e.currentTarget as HTMLInputElement).value)
  }

  function handleValueInput(e: Event) {
    onChange((e.currentTarget as HTMLTextAreaElement).value)
  }
</script>

<div class="space-y-3">
  <div>
    <label for="target-role" class="text-xs font-mono text-gray-500 uppercase tracking-wide mb-2 block">
      Target role (optional)
    </label>
    <input
      id="target-role"
      type="text"
      class="w-full text-sm bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:border-brand-500 transition-colors"
      placeholder="e.g. Software Engineer, Backend Developer..."
      value={targetRole}
      on:input={handleRoleInput}
    />
  </div>
  <div>
    <label for="job-description" class="text-xs font-mono text-gray-500 uppercase tracking-wide mb-2 block">
      Job Description
    </label>
    <textarea
      id="job-description"
      class="input-base"
      rows="8"
      placeholder="Paste the full job description — requirements, tech stack, responsibilities..."
      value={value}
      on:input={handleValueInput}
    ></textarea>
  </div>
</div>