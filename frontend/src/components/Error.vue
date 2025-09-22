<script setup lang="ts">
import { useRouter } from 'vue-router'

const props = defineProps<{
  error: Error
}>()
const router = useRouter()
const handleError = () => {
  router.replace({ path: '/' })
}
</script>

<template>
  <div>
    <h2>{{ error?.message }}</h2>
    <button @click="handleError">Clear errors</button>
  </div>
</template>
