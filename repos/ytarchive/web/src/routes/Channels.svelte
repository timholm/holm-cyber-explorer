<script>
  import { getChannels, addChannel, deleteChannel } from '../lib/api.js';

  let { navigate } = $props();

  let channels = $state([]);
  let loading = $state(true);
  let error = $state(null);

  // Modal state
  let showAddModal = $state(false);
  let channelUrl = $state('');
  let addingChannel = $state(false);
  let addError = $state(null);

  // Delete confirmation
  let deleteConfirm = $state(null);
  let deleting = $state(false);

  async function loadChannels() {
    loading = true;
    error = null;

    try {
      const response = await getChannels();
      channels = response.channels || response || [];
    } catch (err) {
      error = err.message;
      channels = [];
    } finally {
      loading = false;
    }
  }

  async function handleAddChannel() {
    if (!channelUrl.trim()) return;

    addingChannel = true;
    addError = null;

    try {
      await addChannel(channelUrl.trim());
      channelUrl = '';
      showAddModal = false;
      await loadChannels();
    } catch (err) {
      addError = err.message;
    } finally {
      addingChannel = false;
    }
  }

  async function handleDeleteChannel(channel) {
    deleting = true;

    try {
      await deleteChannel(channel.id);
      deleteConfirm = null;
      await loadChannels();
    } catch (err) {
      error = err.message;
    } finally {
      deleting = false;
    }
  }

  function openChannel(channel) {
    navigate('channel', { id: channel.id });
  }

  $effect(() => {
    loadChannels();
  });
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div class="hidden lg:block">
      <h1 class="text-2xl font-bold text-dark-100">Channels</h1>
      <p class="text-dark-400 mt-1">Manage your archived YouTube channels</p>
    </div>
    <button onclick={() => showAddModal = true} class="btn btn-primary">
      <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Channel
    </button>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
    </div>
  {:else if error}
    <div class="card p-6">
      <p class="text-red-400">Error loading channels: {error}</p>
      <button onclick={loadChannels} class="btn btn-secondary mt-4">Retry</button>
    </div>
  {:else if channels.length === 0}
    <div class="card p-12 text-center">
      <div class="w-16 h-16 mx-auto bg-dark-800 rounded-full flex items-center justify-center mb-4">
        <svg class="w-8 h-8 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-dark-200 mb-2">No channels yet</h3>
      <p class="text-dark-400 mb-6">Add your first YouTube channel to start archiving</p>
      <button onclick={() => showAddModal = true} class="btn btn-primary">
        Add Channel
      </button>
    </div>
  {:else}
    <!-- Channel Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {#each channels as channel}
        <div class="card group cursor-pointer hover:border-dark-700 transition-colors" onclick={() => openChannel(channel)} onkeydown={(e) => e.key === 'Enter' && openChannel(channel)} role="button" tabindex="0">
          <!-- Banner -->
          <div class="h-20 bg-gradient-to-r from-dark-800 to-dark-700 relative">
            {#if channel.bannerUrl}
              <img
                src={channel.bannerUrl}
                alt=""
                class="w-full h-full object-cover"
              />
            {/if}
            <!-- Avatar -->
            <div class="absolute -bottom-8 left-4">
              <div class="w-16 h-16 rounded-full border-4 border-dark-900 overflow-hidden bg-dark-700">
                {#if channel.avatarUrl}
                  <img
                    src={channel.avatarUrl}
                    alt={channel.name}
                    class="w-full h-full object-cover"
                  />
                {:else}
                  <div class="w-full h-full flex items-center justify-center text-dark-400 text-xl font-bold">
                    {channel.name?.[0]?.toUpperCase() || '?'}
                  </div>
                {/if}
              </div>
            </div>
          </div>

          <!-- Content -->
          <div class="pt-10 pb-4 px-4">
            <h3 class="font-semibold text-dark-100 truncate group-hover:text-red-400 transition-colors">
              {channel.name || channel.channelId}
            </h3>
            <p class="text-sm text-dark-400 mt-1">
              {channel.videoCount || 0} videos
            </p>

            <!-- Stats -->
            <div class="flex items-center gap-4 mt-3 text-xs text-dark-500">
              <span class="flex items-center gap-1">
                <span class="w-2 h-2 rounded-full bg-green-500"></span>
                {channel.completedCount || 0} done
              </span>
              <span class="flex items-center gap-1">
                <span class="w-2 h-2 rounded-full bg-yellow-500"></span>
                {channel.pendingCount || 0} pending
              </span>
            </div>
          </div>

          <!-- Actions -->
          <div class="px-4 pb-4 flex gap-2">
            <button
              onclick={(e) => { e.stopPropagation(); openChannel(channel); }}
              class="flex-1 btn btn-secondary text-sm py-1.5"
            >
              View
            </button>
            <button
              onclick={(e) => { e.stopPropagation(); deleteConfirm = channel; }}
              class="btn btn-ghost text-sm py-1.5 px-3 text-red-400 hover:text-red-300 hover:bg-red-900/20"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<!-- Add Channel Modal -->
{#if showAddModal}
  <div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" onclick={() => showAddModal = false} onkeydown={(e) => e.key === 'Escape' && (showAddModal = false)} role="button" tabindex="0">
    <div class="bg-dark-900 rounded-xl border border-dark-800 w-full max-w-md p-6" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
      <h2 class="text-xl font-semibold text-dark-100 mb-4">Add Channel</h2>

      <form onsubmit={(e) => { e.preventDefault(); handleAddChannel(); }}>
        <label class="block mb-4">
          <span class="text-sm text-dark-400 mb-2 block">Channel URL</span>
          <input
            type="text"
            bind:value={channelUrl}
            placeholder="https://www.youtube.com/@channelname"
            class="input"
            disabled={addingChannel}
          />
        </label>

        {#if addError}
          <p class="text-red-400 text-sm mb-4">{addError}</p>
        {/if}

        <div class="flex gap-3 justify-end">
          <button
            type="button"
            onclick={() => showAddModal = false}
            class="btn btn-secondary"
            disabled={addingChannel}
          >
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            disabled={addingChannel || !channelUrl.trim()}
          >
            {#if addingChannel}
              <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2 inline-block"></div>
              Adding...
            {:else}
              Add Channel
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Delete Confirmation Modal -->
{#if deleteConfirm}
  <div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" onclick={() => deleteConfirm = null} onkeydown={(e) => e.key === 'Escape' && (deleteConfirm = null)} role="button" tabindex="0">
    <div class="bg-dark-900 rounded-xl border border-dark-800 w-full max-w-md p-6" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
      <h2 class="text-xl font-semibold text-dark-100 mb-2">Delete Channel</h2>
      <p class="text-dark-400 mb-6">
        Are you sure you want to delete <strong class="text-dark-200">{deleteConfirm.name}</strong>? This will also remove all downloaded videos for this channel.
      </p>

      <div class="flex gap-3 justify-end">
        <button
          onclick={() => deleteConfirm = null}
          class="btn btn-secondary"
          disabled={deleting}
        >
          Cancel
        </button>
        <button
          onclick={() => handleDeleteChannel(deleteConfirm)}
          class="btn bg-red-600 hover:bg-red-700 text-white"
          disabled={deleting}
        >
          {#if deleting}
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2 inline-block"></div>
            Deleting...
          {:else}
            Delete
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
