<script lang="ts">
  import '../app.css'
  import Navbar from '$lib/components/Navbar.svelte'
</script>

<div class="min-h-screen flex flex-col">
  <Navbar />
  <main class="flex-1">
    <slot />
  </main>
</div>
