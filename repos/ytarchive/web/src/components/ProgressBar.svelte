<script>
  let { progress = 0, showPercentage = true, size = 'md', animated = true } = $props();

  const clampedProgress = $derived(Math.min(100, Math.max(0, progress)));

  const sizeClasses = {
    sm: 'h-1',
    md: 'h-2',
    lg: 'h-3'
  };
</script>

<div class="w-full">
  <div class="flex items-center gap-3">
    <div class="flex-1 bg-dark-800 rounded-full overflow-hidden {sizeClasses[size]}">
      <div
        class="h-full bg-gradient-to-r from-red-600 to-red-500 rounded-full transition-all duration-300 ease-out {animated && clampedProgress < 100 ? 'animate-pulse' : ''}"
        style="width: {clampedProgress}%"
      >
        {#if animated && clampedProgress > 0 && clampedProgress < 100}
          <div class="h-full w-full relative overflow-hidden">
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent animate-shimmer"></div>
          </div>
        {/if}
      </div>
    </div>

    {#if showPercentage}
      <span class="text-sm text-dark-400 font-medium tabular-nums min-w-[3rem] text-right">
        {clampedProgress.toFixed(0)}%
      </span>
    {/if}
  </div>
</div>

<style>
  @keyframes shimmer {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(100%);
    }
  }

  .animate-shimmer {
    animation: shimmer 1.5s infinite;
  }
</style>
