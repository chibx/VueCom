<script lang="ts" setup>
import MyFooter from '@/components/MyFooter.vue'
import MyHeader from '@/components/MyHeader.vue'
</script>

<template>
  <div>
    <h1>Admin Layout as you can see</h1>
    <MyHeader />
    <main>
      <slot />
    </main>
    <MyFooter />
  </div>
</template>

<style scoped></style>
