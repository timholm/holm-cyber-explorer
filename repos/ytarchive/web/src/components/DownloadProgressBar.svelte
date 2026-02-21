<script>
  import { formatBytes } from '../lib/api.js';

  let {
    percentage = 0,
    downloadedBytes = 0,
    totalBytes = 0,
    speed = '',
    eta = '',
    fragment = '',
    showDetails = true,
    size = 'md'
  } = $props();

  const clampedProgress = $derived(Math.min(100, Math.max(0, percentage)));

  const sizeClasses = {
    sm: 'h-1.5',
    md: 'h-2.5',
    lg: 'h-3.5'
  };

  const displaySpeed = $derived(() => {
    if (speed && speed !== '0 B/s' && speed !== '0.0 B/s') {
      return speed;
    }
    return null;
  });

  const displayProgress = $derived(() => {
    if (totalBytes > 0) {
      return `${formatBytes(downloadedBytes)} / ${formatBytes(totalBytes)}`;
    }
    if (downloadedBytes > 0) {
      return formatBytes(downloadedBytes);
    }
    return null;
  });
</script>

<div class="w-full space-y-2">
  <!-- Progress bar -->
  <div class="flex items-center gap-3">
    <div class="flex-1 bg-dark-800 rounded-full overflow-hidden {sizeClasses[size]}">
      <div
        class="h-full bg-gradient-to-r from-red-600 to-red-500 rounded-full transition-all duration-300 ease-out relative overflow-hidden"
        style="width: {clampedProgress}%"
      >
        {#if clampedProgress > 0 && clampedProgress < 100}
          <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent animate-shimmer"></div>
        {/if}
      </div>
    </div>

    <span class="text-sm text-dark-300 font-medium tabular-nums min-w-[3.5rem] text-right">
      {clampedProgress.toFixed(1)}%
    </span>
  </div>

  <!-- Details row -->
  {#if showDetails}
    <div class="flex items-center justify-between text-xs text-dark-500">
      <div class="flex items-center gap-3">
        {#if displayProgress()}
          <span class="text-dark-400">{displayProgress()}</span>
        {/if}
        {#if fragment}
          <span class="text-dark-500">Segment {fragment}</span>
        {/if}
      </div>

      <div class="flex items-center gap-3">
        {#if displaySpeed()}
          <span class="flex items-center gap-1">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
            {displaySpeed()}
          </span>
        {/if}
        {#if eta && eta !== '00:00'}
          <span class="flex items-center gap-1">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {eta}
          </span>
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  @keyframes shimmer {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(200%);
    }
  }

  .animate-shimmer {
    animation: shimmer 2s infinite;
  }
</style>
