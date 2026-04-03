<template>
  <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
    <h1 class="text-lg sm:text-2xl font-bold text-gray-900 mb-4 sm:mb-6">搜索: {{ keyword }}</h1>
    <div v-if="loading" class="text-center py-8 sm:py-12">
      <div class="text-gray-500 text-sm">搜索中...</div>
    </div>
    <div v-else-if="topics.length > 0" class="space-y-3 sm:space-y-4">
      <div v-for="topic in topics" :key="topic.id" class="border-b pb-3 sm:pb-4 last:border-b-0">
        <router-link :to="`/topic/${topic.id}`" class="block">
          <h3 class="text-base sm:text-lg font-medium text-gray-900 mb-1.5 hover:text-blue-500">{{ topic.title }}</h3>
        </router-link>
        <p class="text-gray-600 text-xs sm:text-sm mb-2 line-clamp-2">{{ stripMarkdown(topic.content).substring(0, 150) }}</p>
        <div class="flex items-center flex-wrap gap-2 text-xs sm:text-sm text-gray-500">
          <span>{{ topic.user?.username }}</span>
          <span>{{ formatTime(topic.created_at) }}</span>
          <span>{{ topic.view_count }} 浏览</span>
        </div>
      </div>
    </div>
    <div v-else class="text-center py-8 sm:py-12 text-gray-500 text-sm">
      没有找到相关内容
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import { stripMarkdown } from '@/utils/markdown'

const route = useRoute()
const keyword = computed(() => route.query.keyword || '')
const topics = ref([])
const loading = ref(false)

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return Math.floor(diff / 86400000) + '天前'
}

async function searchTopics() {
  if (!keyword.value) {
    topics.value = []
    return
  }

  loading.value = true
  try {
    const res = await api.get('/search', {
      params: { keyword: keyword.value }
    })
    topics.value = res?.list || []
  } catch (e) {
    console.error('搜索失败', e)
    topics.value = []
  } finally {
    loading.value = false
  }
}

watch(keyword, () => {
  searchTopics()
}, { immediate: true })
</script>
