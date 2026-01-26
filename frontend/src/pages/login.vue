<script lang="ts" setup>
import { Label } from '@/components/ui/label'
import { Input as UiInput } from '@/components/ui/input'
import { Button as UiButton } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { Checkbox } from '@/components/ui/checkbox'
import { loginSchema, type LoginSchema } from '@/types/schema'
import { useForm } from 'vee-validate'
import { backendFetch, MAX_PASSWORD_LEN } from '@/utils/constants'
import { capitalize, delay } from '@/utils/utils'
import { EyeIcon, EyeOffIcon, LoaderIcon } from 'lucide-vue-next'
import { ref } from 'vue'

const { handleSubmit, errors, defineField } = useForm<LoginSchema>({
  validationSchema: loginSchema,
  initialValues: {},
})

const isSubmitting = ref(false)
const isPassVisible = ref(false)

const [username] = defineField('username')
const [password] = defineField('password')
// const [stayLoggedIn] = defineField('stayLoggedIn')

const onSubmit = handleSubmit(async (_, ac) => {
  if (isSubmitting.value) return
  isSubmitting.value = true

  const formData = new FormData(ac.evt?.target as HTMLFormElement)
  // formData.append('username', values.username)
  // formData.append('password', values.password)
  // formData.append('stayLoggedIn', values.stayLoggedIn ? 'true' : 'false')

  try {
    const data = await backendFetch('/auth/login', {
      method: 'POST',
      body: formData,
    })

    console.log(data)
  } catch (error) {
  } finally {
    await delay(2000)
    isSubmitting.value = false
  }
})
</script>

<template>
  <div class="h-full">
    <div class="h-full w-full flex items-center justify-center">
      <div>
        <h1 class="font-bold text-2xl text-center mb-5">Welcome Back</h1>
        <div class="flex justify-center mb-5">
          <img src="/logo.webp" class="h-25 rounded-full" />
        </div>
        <form
          @submit.prevent="onSubmit"
          class="min-w-90 shadow shadow-gray-400 dark:shadow-none bg-gray-100/50 dark:bg-gray-900 p-5 pt-10 rounded-2xl"
        >
          <div class="flex flex-col gap-1.5 mb-5">
            <Label for="username">Username</Label>
            <UiInput
              type="text"
              name="username"
              id="username"
              placeholder="Enter your username"
              v-model="username"
            />
            <span v-if="errors.username" class="text-error text-xs">
              {{ capitalize(errors.username) }}</span
            >
          </div>

          <div class="flex flex-col gap-1.5 mb-5">
            <Label for="password">Password</Label>
            <div class="relative">
              <span
                class="absolute cursor-pointer text-accent-foreground h-5 right-3.5 -mt-0.5 top-1/2 -translate-y-1/2"
                @click="isPassVisible = !isPassVisible"
              >
                <EyeOffIcon v-if="isPassVisible" class="" />
                <EyeIcon v-else />
              </span>
              <UiInput
                :type="isPassVisible ? 'text' : 'password'"
                name="password"
                id="password"
                placeholder="Enter your password"
                v-model="password"
              />
            </div>
            <span v-if="errors.password" class="text-error text-xs">
              {{ capitalize(errors.password) }}
            </span>
          </div>

          <!-- <div class="flex items-center justify-start gap-2.5 mb-5">
            <Checkbox class="rounded-xs cursor-pointer" type="checkbox" name="stayLoggedIn" id="stayLoggedIn"
              value="true" v-model="stayLoggedIn" />
            <Label class="cursor-pointer" for="stayLoggedIn"> Stay logged in </Label>
          </div> -->

          <span>{{ errors[''] }}</span>
          <div class="flex justify-end">
            <UiButton type="submit" class="cursor-pointer py-5 px-7.5">
              <template v-if="isSubmitting">
                <Spinner />
              </template>
              <template v-else> Login </template>
            </UiButton>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
