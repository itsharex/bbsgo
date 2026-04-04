<template>
  <router-link :to="`/topic/${topic.id}`" class="block p-3 sm:p-4 hover:bg-gray-50">
    <div class="flex items-start">
      <img :src="topic.user?.avatar || 'https://via.placeholder.com/40'" class="w-8 h-8 sm:w-10 sm:h-10 rounded-full flex-shrink-0">
      <div class="ml-2 sm:ml-3 flex-1 min-w-0">
        <div class="flex items-center flex-wrap gap-1">
          <span v-if="topic.is_pinned" class="text-red-500 text-xs">[{{ t('topic.pinned') }}]</span>
          <span v-if="topic.is_essence" class="text-yellow-500 text-xs">[{{ t('topic.essence') }}]</span>
          <span v-if="topic.is_locked" class="text-gray-500 text-xs">[{{ t('topic.lock') }}]</span>
          <h3 class="font-medium text-gray-900 truncate text-sm sm:text-base">{{ topic.title }}</h3>
        </div>
        <div class="mt-1 flex items-center flex-wrap gap-1 text-xs sm:text-sm text-gray-500">
          <span>{{ topic.user?.username }}</span>
          <span class="mx-1">·</span>
          <span>{{ formatDate(topic.created_at) }}</span>
          <span class="mx-1">·</span>
          <span>{{ topic.forum?.name }}</span>
        </div>
        <div class="mt-2 flex items-center flex-wrap gap-2 text-xs text-gray-400">
          <span class="flex items-center">
            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z">
              </path>
            </svg>
            {{ topic.view_count }}
          </span>
          <span class="flex items-center">
            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
              </path>
            </svg>
            {{ topic.like_count }}
          </span>
          <span class="flex items-center">
            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

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

  if (diff < 60000) return t('notifications.justNow')
  if (diff < 3600000) return t('notifications.minutesAgo', { 0: Math.floor(diff / 60000) })
  if (diff < 86400000) return t('notifications.hoursAgo', { 0: Math.floor(diff / 3600000) })
  if (diff < 604800000) return t('notifications.daysAgo', { 0: Math.floor(diff / 86400000) })
  return d.toLocaleDateString('zh-CN')
}
</script>
