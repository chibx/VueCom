<script lang="ts" setup>
import { Toaster } from '@/components/ui/sonner'
import { useDark } from '@vueuse/core'

const dark = useDark()
</script>

<template>
  <div class="h-full bg-slate-50 dark:bg-gray-950">
    <RouterView />
  </div>
  <Toaster />
</template>

<style>
html,
body,
#app {
  height: 100%;
}
</style>
