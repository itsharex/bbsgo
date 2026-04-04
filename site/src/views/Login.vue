<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6">{{ t('login.title') }}</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-medium mb-2">{{ t('login.username') }}</label>
          <input type="text" v-model="form.username"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500" required>
        </div>
        <div class="mb-6">
          <label class="block text-gray-700 text-sm font-medium mb-2">{{ t('login.password') }}</label>
          <input type="password" v-model="form.password"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500" required>
        </div>
        <button type="submit" class="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition-colors">
          {{ t('login.loginBtn') }}
        </button>
      </form>
      <p class="text-center mt-4 text-gray-600">
        {{ t('login.noAccount') }}<router-link to="/register" class="text-blue-500 hover:underline">{{ t('login.goRegister') }}</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const form = ref({
  username: '',
  password: ''
})

async function handleLogin() {
  try {
    await userStore.login(form.value)
    router.push('/')
  } catch (e) {
    console.error(e)
  }
}
</script>
