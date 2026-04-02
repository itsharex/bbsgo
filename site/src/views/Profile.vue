<template>
  <div>
    <div class="bg-white rounded-lg shadow-sm mb-6 overflow-hidden">
      <div class="h-48 bg-cover bg-center" :style="{ backgroundImage: 'url(https://picsum.photos/1200/400)' }">
        <div class="flex justify-end p-4">
          <button class="bg-white/90 px-3 py-1 rounded text-sm text-gray-600 hover:bg-white">
            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z">
              </path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
            </svg>
            设置背景
          </button>
        </div>
      </div>
      <div class="px-6 pb-6">
        <div class="flex items-end -mt-16">
          <img :src="user?.avatar || 'https://via.placeholder.com/128'"
            class="w-32 h-32 rounded-full border-4 border-white shadow-lg">
          <div class="ml-4 mb-2">
            <h1 class="text-xl font-bold text-gray-900">{{ user?.username || '用户名' }}</h1>
            <p class="text-gray-500 text-sm">{{ user?.signature || '这家伙很懒，什么都没留下...' }}</p>
          </div>
        </div>
      </div>
    </div>
    <div class="flex gap-6">
      <aside class="w-72 flex-shrink-0">
        <div class="space-y-4">
          <div class="bg-white rounded-lg shadow-sm p-4">
            <h3 class="font-medium text-gray-900 mb-4 border-b pb-2">个人成就</h3>
            <div class="grid grid-cols-4 gap-4 text-center">
              <div>
                <div class="text-2xl font-bold text-gray-700">{{ user?.credits || 0 }}</div>
                <div class="text-xs text-gray-400">积分</div>
              </div>
              <div>
                <div class="text-2xl font-bold text-gray-700">{{ userStats.topic_count || 0 }}</div>
                <div class="text-xs text-gray-400">帖子</div>
              </div>
              <div>
                <div class="text-2xl font-bold text-gray-700">{{ userStats.post_count || 0 }}</div>
                <div class="text-xs text-gray-400">评论</div>
              </div>
              <div>
                <div class="text-2xl font-bold text-gray-700">{{ userStats.rank || 0 }}</div>
                <div class="text-xs text-gray-400">注册排名</div>
              </div>
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-4">
            <div class="flex justify-between items-center mb-4">
              <h3 class="font-medium text-gray-900">个人资料</h3>
              <button class="text-blue-500 text-sm hover:underline">编辑资料</button>
            </div>
            <div class="space-y-3">
              <div class="flex">
                <span class="w-20 text-gray-500 text-sm">昵称</span>
                <span class="text-gray-900 text-sm">{{ user?.nickname || user?.username }}</span>
              </div>
              <div class="flex">
                <span class="w-20 text-gray-500 text-sm">签名</span>
                <span class="text-gray-900 text-sm">{{ user?.signature || '-' }}</span>
              </div>
              <div class="flex">
                <span class="w-20 text-gray-500 text-sm">主页</span>
                <span class="text-blue-500 text-sm">{{ user?.intro ? user.intro : 'https://mlog.club/user/' + (user?.id || '') }}</span>
              </div>
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-4">
            <div class="flex justify-between items-center mb-4">
              <h3 class="font-medium text-gray-900">粉丝 {{ followers.length }}</h3>
              <button class="text-blue-500 text-sm hover:underline">更多</button>
            </div>
            <div v-if="followers.length > 0" class="space-y-3">
              <div v-for="follower in followers" :key="follower.id" class="flex items-center space-x-3">
                <img :src="follower.avatar || 'https://via.placeholder.com/40'" class="w-10 h-10 rounded-full">
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-medium text-gray-900 truncate">{{ follower.username }}</div>
                  <div class="text-xs text-gray-400 truncate">{{ follower.signature || '这家伙很懒，什么都没留下' }}</div>
                </div>
                <button class="bg-blue-500 text-white text-xs px-3 py-1 rounded hover:bg-blue-600">+ 关注</button>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-4 text-sm">
              暂无粉丝
            </div>
          </div>
        </div>
      </aside>
      <div class="flex-1 min-w-0">
        <div class="bg-white rounded-lg shadow-sm">
          <div class="p-4">
            <div v-if="userTopics.length > 0" class="space-y-4">
              <div v-for="topic in userTopics" :key="topic.id" class="border-b pb-4 last:border-b-0">
                <div class="flex items-center justify-between mb-1">
                  <span class="text-sm text-gray-500">{{ topic.user?.username }}</span>
                  <span class="text-xs text-gray-400">{{ formatTime(topic.created_at) }}</span>
                </div>
                <router-link :to="`/topic/${topic.id}`" class="block">
                  <h3 class="text-lg font-medium text-gray-900 mb-2 hover:text-blue-500">{{ topic.title }}</h3>
                </router-link>
                <p class="text-gray-600 text-sm mb-3 line-clamp-3">{{ topic.content.substring(0, 200) }}</p>
                <div class="flex items-center space-x-4 text-xs text-gray-500">
                  <span>心赞 {{ topic.like_count || 0 }}</span>
                  <span>评论 {{ topic.reply_count || 0 }}</span>
                  <span>浏览 {{ topic.view_count || 0 }}</span>
                  <span v-if="topic.forum" class="bg-gray-100 px-2 py-0.5 rounded">{{ topic.forum.name }}</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-12">
              暂无帖子
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
import api from '@/api'

const route = useRoute()
const user = ref(null)
const userStats = ref({
  topic_count: 0,
  post_count: 0,
  rank: 0
})
const userTopics = ref([])
const followers = ref([])

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return Math.floor(diff / 86400000) + '天前'
}

async function loadUser() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}`)
    user.value = res
  } catch (e) {
    console.error('加载用户信息失败', e)
  }
}

async function loadUserTopics() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/topics`)
    userTopics.value = res?.list || []
  } catch (e) {
    console.error('加载用户帖子失败', e)
  }
}

async function loadFollowers() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/followers`)
    followers.value = res?.list || []
  } catch (e) {
    console.error('加载粉丝列表失败', e)
  }
}

async function loadUserStats() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/stats`)
    userStats.value = res || {
      topic_count: 0,
      post_count: 0,
      rank: 0
    }
  } catch (e) {
    console.error('加载用户统计失败', e)
  }
}

onMounted(async () => {
  await Promise.all([
    loadUser(),
    loadUserTopics(),
    loadFollowers(),
    loadUserStats()
  ])
})
</script>
