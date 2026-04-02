<template>
  <router-link :to="`/topic/${topic.id}`" class="block p-4 hover:bg-gray-50">
    <div class="flex items-start">
      <img :src="topic.user?.avatar || 'https://via.placeholder.com/40'" class="w-10 h-10 rounded-full flex-shrink-0">
      <div class="ml-3 flex-1 min-w-0">
        <div class="flex items-center">
          <span v-if="topic.is_pinned" class="text-red-500 text-xs mr-1">[置顶]</span>
          <span v-if="topic.is_essence" class="text-yellow-500 text-xs mr-1">[精华]</span>
          <span v-if="topic.is_locked" class="text-gray-500 text-xs mr-1">[锁定]</span>
          <h3 class="font-medium text-gray-900 truncate">{{ topic.title }}</h3>
        </div>
        <div class="mt-1 flex items-center text-sm text-gray-500">
          <span>{{ topic.user?.username }}</span>
          <span class="mx-1">·</span>
          <span>{{ formatDate(topic.created_at) }}</span>
          <span class="mx-1">·</span>
          <span>{{ topic.forum?.name }}</span>
        </div>
        <div class="mt-2 flex items-center text-xs text-gray-400">
          <span class="flex items-center mr-4">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z">
              </path>
            </svg>
            {{ topic.view_count }}
          </span>
          <span class="flex items-center mr-4">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
              </path>
            </svg>
            {{ topic.like_count }}
          </span>
          <span class="flex items-center">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
              </path>
            </svg>
            {{ topic.reply_count }}
          </span>
        </div>
      </div>
    </div>
  </router-link>
</template>

<script setup>
defineProps({
  topic: {
    type: Object,
    required: true
  }
})

function formatDate(date) {
  const d = new Date(date)
  const now = new Date()
  const diff = now - d

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} 天前`
  return d.toLocaleDateString('zh-CN')
}
</script>
