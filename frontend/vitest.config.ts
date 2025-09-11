import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    name: 'nuxt',
    include: ['tests/{e2e,unit}/*.{test,spec}.ts'],
    environment: 'nuxt',
  },
})
