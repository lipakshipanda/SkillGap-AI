<script lang="ts">
  export let value: string
  export let onChange: (v: string) => void
  export let onFileChange: (f: File | null) => void
  export let file: File | null
  export let mode: 'text' | 'pdf'
  export let setMode: (m: 'text' | 'pdf') => void

  let isDragActive = false
  let fileInput: HTMLInputElement

  function handleFiles(files: FileList | null) {
    if (!files || !files[0]) return
    const f = files[0]
    if (f.type === 'application/pdf' || f.name.toLowerCase().endsWith('.pdf')) {
      onFileChange(f)
    }
  }

  function onDrop(e: DragEvent) {
    e.preventDefault()
    isDragActive = false
    handleFiles(e.dataTransfer?.files ?? null)
  }

  function handleTextareaInput(e: Event) {
    onChange((e.currentTarget as HTMLTextAreaElement).value)
  }

  function handleFileInputChange(e: Event) {
    handleFiles((e.currentTarget as HTMLInputElement).files)
  }
</script>

<div>
  <div class="flex items-center justify-between mb-2">
    <span id="resume-label" class="text-xs font-mono text-gray-500 uppercase tracking-wide">Your Resume</span>
    <div class="flex rounded-md border border-gray-200 overflow-hidden text-xs font-mono">
      <button
        on:click={() => {
          setMode('text')
          onFileChange(null)
        }}
        class="px-3 py-1 transition-colors {mode === 'text' ? 'bg-brand-900 text-white' : 'text-gray-500 hover:bg-gray-50'}"
      >
        paste text
      </button>
      <button
        on:click={() => setMode('pdf')}
        class="px-3 py-1 transition-colors {mode === 'pdf' ? 'bg-brand-900 text-white' : 'text-gray-500 hover:bg-gray-50'}"
      >
        upload pdf
      </button>
    </div>
  </div>

  {#if mode === 'text'}
    <textarea
      aria-labelledby="resume-label"
      class="input-base"
      rows="8"
      placeholder="Paste your resume here — education, skills, projects, experience..."
      value={value}
      on:input={handleTextareaInput}
    ></textarea>
  {:else if file}
    <div class="flex items-center gap-3 p-4 bg-emerald-50 border border-emerald-200 rounded-lg">
      <span class="text-emerald-600 shrink-0">📄</span>
      <div class="flex-1 min-w-0">
        <p class="text-sm font-medium text-emerald-800 truncate">{file.name}</p>
        <p class="text-xs text-emerald-600 font-mono">{(file.size / 1024).toFixed(1)} KB</p>
      </div>
      <button on:click={() => onFileChange(null)} class="text-emerald-600 hover:text-emerald-800">✕</button>
    </div>
  {:else}
    <div
      role="button"
      tabindex="0"
      on:click={() => fileInput.click()}
      on:keydown={(e) => e.key === 'Enter' && fileInput.click()}
      on:dragover={(e) => {
        e.preventDefault()
        isDragActive = true
      }}
      on:dragleave={() => (isDragActive = false)}
      on:drop={onDrop}
      class="border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors {isDragActive
        ? 'border-brand-500 bg-brand-50'
        : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'}"
    >
      <input
        bind:this={fileInput}
        type="file"
        accept="application/pdf,.pdf"
        class="hidden"
        on:change={handleFileInputChange}
      />
      <div class="mx-auto mb-2 text-gray-400">⬆</div>
      <p class="text-sm text-gray-600 font-medium">
        {isDragActive ? 'Drop your PDF here' : 'Drag & drop your resume PDF'}
      </p>
      <p class="text-xs text-gray-400 mt-1">or click to browse</p>
    </div>
  {/if}
</div>