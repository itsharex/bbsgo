<template>
  <div class="bg-white rounded-lg shadow-sm p-6">
    <h1 class="text-2xl font-bold text-gray-900 mb-6">搜索: {{ keyword }}</h1>
    <div v-if="topics.length > 0" class="space-y-4">
      <div v-for="topic in topics" :key="topic.id" class="border-b pb-4 last:border-b-0">
        <router-link :to="`/topic/${topic.id}`" class="block">
          <h3 class="text-lg font-medium text-gray-900 mb-2 hover:text-blue-500">{{ topic.title }}</h3>
        </router-link>
        <p class="text-gray-600 text-sm mb-2 line-clamp-2">{{ topic.content.substring(0, 150) }}</p>
        <div class="flex items-center space-x-4 text-sm text-gray-500">
          <span>{{ topic.user?.username }}</span>
          <span>{{ formatTime(topic.created_at) }}</span>
          <span>{{ topic.view_count }} 浏览</span>
        </div>
      </div>
    </div>
    <div v-else class="text-center py-12 text-gray-500">
      没有找到相关内容
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const keyword = computed(() => route.query.keyword || '')
const topics = ref([])

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return Math.floor(diff / 86400000) + '天前'
}

onMounted(() => {
  topics.value = [
    {
      id: 1,
      title: 'bbs-go 相关搜索结果',
      content: '这是关于 ' + keyword.value + ' 的搜索结果...',
      user: { username: '小码哥' },
      created_at: new Date().toISOString(),
      view_count: 100
    }
  ]
})
</script>
