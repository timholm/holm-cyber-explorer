<script>
  import { getCookies, saveCookies, deleteCookies } from '../lib/api.js';

  let cookies = $state('');
  let loading = $state(true);
  let saving = $state(false);
  let deleting = $state(false);
  let error = $state(null);
  let success = $state(null);
  let hasCookies = $state(false);

  async function loadCookies() {
    loading = true;
    error = null;

    try {
      const response = await getCookies();
      cookies = response.cookies || '';
      hasCookies = !!cookies;
    } catch (err) {
      if (err.status !== 404) {
        error = err.message;
      }
      cookies = '';
      hasCookies = false;
    } finally {
      loading = false;
    }
  }

  async function handleSave() {
    saving = true;
    error = null;
    success = null;

    try {
      await saveCookies(cookies);
      hasCookies = !!cookies;
      success = 'Cookies saved successfully';
      setTimeout(() => success = null, 3000);
    } catch (err) {
      error = err.message;
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!confirm('Are you sure you want to delete your cookies?')) {
      return;
    }

    deleting = true;
    error = null;
    success = null;

    try {
      await deleteCookies();
      cookies = '';
      hasCookies = false;
      success = 'Cookies deleted successfully';
      setTimeout(() => success = null, 3000);
    } catch (err) {
      error = err.message;
    } finally {
      deleting = false;
    }
  }

  $effect(() => {
    loadCookies();
  });
</script>

<div class="space-y-8">
  <!-- Header -->
  <div class="hidden lg:block">
    <h1 class="text-2xl font-bold text-dark-100">Settings</h1>
    <p class="text-dark-400 mt-1">Configure your YouTube archive settings</p>
  </div>

  <!-- Cookies Section -->
  <div class="card">
    <div class="px-6 py-4 border-b border-dark-800">
      <h2 class="text-lg font-semibold text-dark-100">YouTube Cookies</h2>
      <p class="text-dark-400 text-sm mt-1">
        Cookies are required to download age-restricted or members-only content
      </p>
    </div>

    <div class="p-6">
      {#if loading}
        <div class="flex items-center justify-center py-8">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-red-500"></div>
        </div>
      {:else}
        <!-- Status indicator -->
        <div class="flex items-center gap-3 mb-6">
          <div class="w-3 h-3 rounded-full {hasCookies ? 'bg-green-500' : 'bg-yellow-500'}"></div>
          <span class="text-dark-300">
            {#if hasCookies}
              Cookies configured
            {:else}
              No cookies configured
            {/if}
          </span>
        </div>

        <!-- Cookie input -->
        <div class="space-y-4">
          <div>
            <label for="cookies" class="block text-sm font-medium text-dark-300 mb-2">
              Cookie String (Netscape format)
            </label>
            <textarea
              id="cookies"
              bind:value={cookies}
              placeholder="Paste your cookies here in Netscape format...

Example:
.youtube.com	TRUE	/	TRUE	1234567890	VISITOR_INFO1_LIVE	value
.youtube.com	TRUE	/	TRUE	1234567890	YSC	value"
              rows="10"
              class="input font-mono text-sm"
            ></textarea>
          </div>

          <div class="bg-dark-800 rounded-lg p-4">
            <h3 class="text-sm font-medium text-dark-200 mb-2">How to get cookies:</h3>
            <ol class="text-sm text-dark-400 space-y-2 list-decimal list-inside">
              <li>Install a browser extension like "Get cookies.txt LOCALLY"</li>
              <li>Go to youtube.com and sign in</li>
              <li>Click the extension and export cookies for youtube.com</li>
              <li>Paste the contents here</li>
            </ol>
          </div>

          <!-- Error/Success messages -->
          {#if error}
            <div class="p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
              <p class="text-red-400 text-sm">{error}</p>
            </div>
          {/if}

          {#if success}
            <div class="p-4 bg-green-500/10 border border-green-500/20 rounded-lg">
              <p class="text-green-400 text-sm">{success}</p>
            </div>
          {/if}

          <!-- Actions -->
          <div class="flex gap-3">
            <button
              onclick={handleSave}
              disabled={saving}
              class="btn btn-primary"
            >
              {#if saving}
                <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2 inline-block"></div>
                Saving...
              {:else}
                <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                Save Cookies
              {/if}
            </button>

            {#if hasCookies}
              <button
                onclick={handleDelete}
                disabled={deleting}
                class="btn btn-secondary text-red-400 hover:text-red-300"
              >
                {#if deleting}
                  <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-red-400 mr-2 inline-block"></div>
                  Deleting...
                {:else}
                  <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  Delete Cookies
                {/if}
              </button>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  </div>

  <!-- About Section -->
  <div class="card p-6">
    <h2 class="text-lg font-semibold text-dark-100 mb-4">About</h2>
    <div class="space-y-3 text-sm">
      <div class="flex justify-between">
        <span class="text-dark-400">Version</span>
        <span class="text-dark-200">1.0.0</span>
      </div>
      <div class="flex justify-between">
        <span class="text-dark-400">Backend</span>
        <span class="text-dark-200">Go + SQLite FTS5</span>
      </div>
      <div class="flex justify-between">
        <span class="text-dark-400">Frontend</span>
        <span class="text-dark-200">Svelte 5 + Tailwind CSS</span>
      </div>
    </div>
  </div>
</div>
