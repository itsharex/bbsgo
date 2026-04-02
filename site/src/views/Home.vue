<template>
  <div class="flex gap-6 max-w-6xl mx-auto">
    <aside class="w-56 flex-shrink-0 hidden lg:block">
      <div class="bg-white rounded-lg shadow-sm overflow-hidden">
        <div class="px-4 py-3 border-b bg-gray-50">
          <h3 class="font-semibold text-gray-700">热门话题</h3>
        </div>
        <router-link v-for="tag in tags" :key="tag.id"
          :to="tag.id ? `/?tag=${tag.id}` : '/'"
          :class="['px-4 py-3 flex items-center justify-between transition-colors',
            currentTagId === tag.id ? 'bg-blue-50 text-blue-600' : 'text-gray-600 hover:bg-gray-50']">
          <div class="flex items-center space-x-2">
            <span v-if="tag.icon" class="text-lg">{{ tag.icon }}</span>
            <span class="font-medium">{{ tag.name }}</span>
          </div>
          <span class="text-xs text-gray-400">{{ tag.usage_count }}</span>
        </router-link>
      </div>
    </aside>
    <div class="flex-1 min-w-0">
      <div class="space-y-4">
        <div v-for="topic in topics" :key="topic.id"
          class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex space-x-4">
            <router-link :to="`/user/${topic.user_id}`">
              <img :src="topic.user?.avatar || 'https://via.placeholder.com/48'" class="w-12 h-12 rounded-full">
            </router-link>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between mb-1">
                <div class="flex items-center space-x-2">
                  <router-link :to="`/user/${topic.user_id}`" class="font-medium text-gray-900 hover:text-blue-500">
                    {{ topic.user?.nickname || topic.user?.username }}
                  </router-link>
                  <span v-if="topic.forum" class="text-xs bg-blue-100 text-blue-600 px-2 py-0.5 rounded">
                    {{ topic.forum.name }}
                  </span>
                </div>
                <span class="text-xs text-gray-400">{{ formatTime(topic.created_at) }}</span>
              </div>
              <router-link :to="`/topic/${topic.id}`" class="block">
                <h3 class="text-lg font-semibold text-gray-900 mb-2 hover:text-blue-500 line-clamp-2">
                  {{ topic.title }}
                </h3>
                <p class="text-gray-600 text-sm mb-3 line-clamp-3" v-html="topic.content.substring(0, 200)"></p>
              </router-link>
              <div class="flex items-center flex-wrap gap-2 mb-2" v-if="topic.tags && topic.tags.length > 0">
                <router-link v-for="tag in topic.tags" :key="tag.id"
                  :to="`/?tag=${tag.id}`"
                  class="px-2 py-0.5 text-xs bg-gray-100 text-gray-600 rounded hover:bg-blue-100 hover:text-blue-600">
                  #{{ tag.name }}
                </router-link>
              </div>
              <div class="flex items-center space-x-6 text-sm text-gray-500">
                <button @click="toggleLike(topic)" :class="['flex items-center space-x-1 transition-colors', topic.liked ? 'text-red-500' : 'hover:text-red-500']">
                  <svg class="w-4 h-4" :fill="topic.liked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                    </path>
                  </svg>
                  <span>{{ topic.like_count }}</span>
                </button>
                <router-link :to="`/topic/${topic.id}`" class="flex items-center space-x-1 hover:text-blue-500">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                    </path>
                  </svg>
                  <span>{{ topic.reply_count }}</span>
                </router-link>
                <button class="flex items-center space-x-1 hover:text-green-500">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z">
                    </path>
                  </svg>
                  <span>{{ topic.view_count }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <div ref="loadMoreTrigger" class="py-8 text-center">
          <div v-if="loading" class="flex items-center justify-center space-x-2">
            <svg class="animate-spin w-5 h-5 text-blue-500" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span class="text-gray-500">加载中...</span>
          </div>
          <div v-else-if="noMore" class="text-gray-400 text-sm">
            已经到底啦~
          </div>
        </div>
      </div>
    </div>
    <aside class="w-52 flex-shrink-0 hidden xl:block">
      <div class="bg-white rounded-lg shadow-sm p-4 mb-4" v-if="userStore.isLoggedIn">
        <h3 class="font-semibold text-gray-900 mb-3">签到</h3>
        <div class="text-center">
          <div class="text-3xl font-bold text-gray-800 mb-2">{{ userStore.user?.credits || 0 }}</div>
          <div class="text-xs text-gray-500 mb-4">积分</div>
          <button @click="handleSignIn" :disabled="signInLoading || signInStatus.signed_today"
            :class="['w-full py-2 rounded-lg font-medium transition-colors', 
              signInStatus.signed_today 
                ? 'bg-gray-100 text-gray-400 cursor-not-allowed' 
                : 'bg-blue-500 text-white hover:bg-blue-600']">
            {{ signInLoading ? '签到中...' : (signInStatus.signed_today ? '今日已签到' : '立即签到') }}
          </button>
          <div v-if="signInStatus.last_sign_at" class="mt-3 text-xs text-gray-500">
            你已经连续签到 {{ getStreakDays() }} 天啦！
          </div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 mb-4">
        <h3 class="font-semibold text-gray-900 mb-3">热门帖子</h3>
        <div class="space-y-3">
          <router-link v-for="t in hotTopics" :key="t.id" :to="`/topic/${t.id}`" class="block group">
            <div class="text-sm text-gray-700 group-hover:text-blue-500 line-clamp-2">{{ t.title }}</div>
            <div class="text-xs text-gray-400 mt-1">{{ t.view_count }} 浏览</div>
          </router-link>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <h3 class="font-semibold text-gray-900 mb-3">活跃用户</h3>
        <div class="space-y-3">
          <div v-for="(user, index) in creditUsers" :key="user.id" class="flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <img :src="user.avatar || 'https://via.placeholder.com/24'" class="w-6 h-6 rounded-full">
              <span class="text-sm text-gray-700">{{ user.nickname || user.username }}</span>
            </div>
            <span class="text-xs font-medium text-gray-600">{{ user.credits }}</span>
          </div>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useIntersectionObserver } from '@vueuse/core'
import api from '@/api'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const tags = ref([])
const topics = ref([])
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)
const noMore = ref(false)
const loadMoreTrigger = ref(null)

const hotTopics = ref([])
const creditUsers = ref([])
const signInStatus = ref({
  signed_today: false,
  last_sign_at: null,
  credits: 0
})
const signInLoading = ref(false)

const currentTagId = computed(() => {
  const tagId = route.query.tag
  return tagId ? parseInt(tagId) : null
})

const currentForum = computed(() => {
  const forumId = route.query.forum
  return forumId ? parseInt(forumId) : null
})

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return Math.floor(diff / 86400000) + '天前'
}

async function loadTags() {
  try {
    const res = await api.get('/tags')
    tags.value = res || []
  } catch (e) {
    console.error(e)
  }
}

async function loadTopics(isLoadMore = false) {
  if (loading.value || noMore.value) return
  
  loading.value = true
  
  try {
    const params = {
      page: page.value,
      page_size: pageSize
    }
    if (currentForum.value) {
      params.forum_id = currentForum.value
    }
    if (currentTagId.value) {
      params.tag_id = currentTagId.value
    }
    const res = await api.get('/topics', { params })
    
    if (isLoadMore) {
      topics.value = [...topics.value, ...(res.list || [])]
    } else {
      topics.value = res.list || []
    }
    
    total.value = res.total || 0
    
    if (topics.value.length >= total.value) {
      noMore.value = true
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadHotTopics() {
  try {
    const res = await api.get('/topics', {
      params: {
        page: 1,
        page_size: 5,
        order_by: 'view_count'
      }
    })
    hotTopics.value = res.list || []
  } catch (e) {
    console.error(e)
  }
}

async function loadCreditUsers() {
  try {
    const res = await api.get('/users/credit')
    creditUsers.value = res || []
  } catch (e) {
    console.error(e)
  }
}

function loadMore() {
  if (!loading.value && !noMore.value) {
    page.value++
    loadTopics(true)
  }
}

useIntersectionObserver(
  loadMoreTrigger,
  ([{ isIntersecting }]) => {
    if (isIntersecting) {
      loadMore()
    }
  },
  { threshold: 0.1 }
)

watch([() => route.query.forum, () => route.query.tag], () => {
  page.value = 1
  noMore.value = false
  topics.value = []
  loadTopics()
})

async function toggleLike(topic) {
  if (!userStore.isLoggedIn) {
    Swal.fire('提示', '请先登录', 'warning')
    return
  }

  try {
    if (topic.liked) {
      await api.delete(`/likes?target_type=topic&target_id=${topic.id}`)
      topic.like_count--
    } else {
      await api.post('/likes', {
        target_type: 'topic',
        target_id: topic.id
      })
      topic.like_count++
    }
    topic.liked = !topic.liked
  } catch (e) {
    console.error(e)
    Swal.fire('错误', '操作失败', 'error')
  }
}

async function loadSignInStatus() {
  if (!userStore.isLoggedIn) return
  try {
    const res = await api.get('/user/signin/status')
    signInStatus.value.signed_today = res.signed_today || false
    signInStatus.value.last_sign_at = res.last_sign_at || null
    signInStatus.value.credits = res.credits || 0
  } catch (e) {
    console.error(e)
  }
}

async function handleSignIn() {
  if (!userStore.isLoggedIn) return
  signInLoading.value = true
  try {
    const res = await api.post('/user/signin')
    Swal.fire({
      icon: 'success',
      title: '签到成功',
      text: `获得 ${res.credits} 积分，总积分 ${res.total_credits}`
    })
    signInStatus.value.signed_today = true
    signInStatus.value.last_sign_at = new Date().toISOString()
    signInStatus.value.credits = res.total_credits
    if (userStore.user) {
      userStore.user.credits = res.total_credits
    }
  } catch (e) {
    console.error(e)
    Swal.fire('错误', e.response?.data?.message || '签到失败', 'error')
  } finally {
    signInLoading.value = false
  }
}

function getStreakDays() {
  if (!signInStatus.value.last_sign_at) return 0
  const today = new Date().toISOString().split('T')[0]
  const lastSign = new Date(signInStatus.value.last_sign_at).toISOString().split('T')[0]
  
  if (today === lastSign || 
      (new Date(today).getTime() - new Date(lastSign).getTime()) <= 86400000) {
    return 2
  }
  return 1
}

async function checkTopicLikes() {
  if (!userStore.isLoggedIn || topics.value.length === 0) return
  
  for (const topic of topics.value) {
    try {
      const res = await api.post('/likes/check', {
        target_type: 'topic',
        target_id: topic.id
      })
      topic.liked = res.liked
    } catch (e) {
      console.error(e)
    }
  }
}

onMounted(() => {
  loadTags()
  loadTopics()
  loadHotTopics()
  loadCreditUsers()
  if (userStore.isLoggedIn) {
    loadSignInStatus()
  }
})

watch(() => topics.value, () => {
  checkTopicLikes()
}, { deep: true })

watch(() => userStore.isLoggedIn, () => {
  if (userStore.isLoggedIn) {
    loadSignInStatus()
    checkTopicLikes()
  }
})
</script>
