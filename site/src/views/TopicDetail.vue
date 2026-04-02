<template>
  <div class="bg-white rounded-lg shadow-sm p-6">
    <div v-if="topic" class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-4">{{ topic.title }}</h1>
      <div v-if="topic.tags && topic.tags.length > 0" class="flex items-center flex-wrap gap-2 mb-4">
        <router-link v-for="tag in topic.tags" :key="tag.id"
          :to="`/?tag=${tag.id}`"
          class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded-full hover:bg-blue-200">
          #{{ tag.name }}
        </router-link>
      </div>
      <div class="flex items-center space-x-4 mb-6 pb-6 border-b">
        <router-link :to="`/user/${topic.user_id}`">
          <img :src="topic.user?.avatar || 'https://via.placeholder.com/48'" class="w-12 h-12 rounded-full">
        </router-link>
        <div>
          <router-link :to="`/user/${topic.user_id}`" class="font-medium text-gray-900 hover:text-blue-500">{{
            topic.user?.username }}</router-link>
          <div class="text-sm text-gray-500">{{ formatTime(topic.created_at) }} · {{ topic.view_count }} 浏览</div>
        </div>
      </div>
      <div class="prose max-w-none mb-6" v-html="topic.content"></div>
      <div class="flex items-center space-x-4 pt-4 border-t">
        <button class="flex items-center space-x-2 text-gray-500 hover:text-red-500">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
            </path>
          </svg>
          <span>{{ topic.like_count }}</span>
        </button>
        <button class="flex items-center space-x-2 text-gray-500 hover:text-blue-500">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
          </svg>
          <span>收藏</span>
        </button>
        <button class="flex items-center space-x-2 text-gray-500 hover:text-green-500">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z">
            </path>
          </svg>
          <span>分享</span>
        </button>
      </div>
    </div>
    <div class="mt-8">
      <h3 class="text-lg font-medium text-gray-900 mb-4">{{ posts.length }} 条评论</h3>
      <div v-if="userStore.isLoggedIn" class="mb-6">
        <textarea v-model="newPost" rows="3"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
          placeholder="写下你的评论..."></textarea>
        <div class="flex justify-end mt-2">
          <button @click="submitPost"
            class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600">发表评论</button>
        </div>
      </div>
      <div class="space-y-4">
        <div v-for="post in posts" :key="post.id" class="flex space-x-4 p-4 bg-gray-50 rounded-lg">
          <img :src="post.user?.avatar || 'https://via.placeholder.com/40'" class="w-10 h-10 rounded-full">
          <div class="flex-1">
            <div class="flex items-center space-x-2 mb-1">
              <span class="font-medium text-gray-900">{{ post.user?.username }}</span>
              <span class="text-sm text-gray-500">{{ formatTime(post.created_at) }}</span>
            </div>
            <p class="text-gray-700">{{ post.content }}</p>
            <div class="flex items-center space-x-4 mt-2 text-sm">
              <button class="text-gray-500 hover:text-red-500">❤️ {{ post.like_count }}</button>
              <button class="text-gray-500 hover:text-blue-500">回复</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import api from '@/api'

const route = useRoute()
const userStore = useUserStore()
const topic = ref(null)
const posts = ref([])
const newPost = ref('')

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return Math.floor(diff / 86400000) + '天前'
}

async function loadTopic() {
  try {
    const id = route.params.id
    topic.value = await api.get(`/topics/${id}`)
    posts.value = [
      {
        id: 1,
        content: '支持一下，期待新功能！',
        user: { username: '用户A', avatar: '' },
        like_count: 2,
        created_at: new Date(Date.now() - 3600000).toISOString()
      }
    ]
  } catch (e) {
    console.error(e)
  }
}

async function submitPost() {
  if (!newPost.value.trim()) return
  try {
    await api.post(`/topics/${route.params.id}/posts`, { content: newPost.value })
    newPost.value = ''
    loadTopic()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadTopic()
})
</script>
