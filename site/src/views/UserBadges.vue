<template>
  <div class="max-w-6xl mx-auto px-4 py-6">
    <div class="bg-white rounded-2xl shadow-lg p-8 mb-6">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 flex items-center gap-2">
            <svg class="w-7 h-7 text-yellow-500" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
            {{ t('userBadges.badgeGallery') }}
          </h1>
          <p class="text-gray-500 mt-2">{{ user?.username ? t('userBadges.badgeAchievementOf', { name: user.username }) : '' }}</p>
        </div>
        <router-link :to="`/user/${userId}`" class="flex items-center gap-2 text-blue-500 hover:text-blue-600 transition-colors px-4 py-2 rounded-lg hover:bg-blue-50">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
          </svg>
          {{ t('userBadges.returnToProfile') }}
        </router-link>
      </div>

      <div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6 xl:grid-cols-8 gap-2">
        <div v-for="item in badgeProgress" :key="item.badge.id"
          class="group relative rounded-lg p-2 transition-all duration-200 hover:-translate-y-0.5 cursor-pointer"
          :class="[
            item.awarded
              ? 'bg-gradient-to-br from-white to-yellow-50 border border-yellow-200 shadow-sm hover:shadow-md'
              : 'bg-gradient-to-br from-gray-50 to-gray-100 border border-gray-200 opacity-60 hover:opacity-80'
          ]">
          <!-- 状态标识 -->
          <div v-if="item.awarded" class="absolute -top-0.5 -right-0.5 flex items-center justify-center w-4 h-4 bg-green-500 text-white rounded-full text-xs shadow-sm">
            ✓
          </div>
          <div v-else class="absolute -top-0.5 -right-0.5 flex items-center justify-center w-4 h-4 bg-gray-400 text-white rounded-full shadow-sm text-xs">
            🔒
          </div>

          <el-tooltip :content="item.badge.description || item.badge.name" placement="top" :show-after="300">
            <div class="flex flex-col items-center text-center">
              <!-- 勋章图标 -->
              <div :class="['relative transition-all', item.awarded ? '' : 'grayscale']">
                <div v-if="item.awarded" class="absolute inset-0 bg-yellow-200 rounded-full blur-md opacity-30"></div>
                <SvgBadge :type="item.badge.icon" :size="item.awarded ? 32 : 24" class="relative" />
              </div>

              <!-- 名称 -->
              <h3 :class="['font-medium text-xs w-full truncate mt-1', item.awarded ? 'text-gray-800' : 'text-gray-500']">
                {{ item.badge.name }}
              </h3>

              <!-- 类型标签 -->
              <el-tag :type="getTypeColor(item.badge.type)" size="small" class="text-xs px-1 mt-0.5">
                {{ getTypeName(item.badge.type) }}
              </el-tag>

              <!-- 已获得时间 -->
              <div v-if="item.awarded_at" class="text-xs text-green-600 mt-0.5">
                {{ formatDateTime(item.awarded_at) }}
              </div>
            </div>
          </el-tooltip>
        </div>
      </div>

      <div v-if="loading" class="text-center py-16">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-blue-50 mb-4">
          <el-icon class="is-loading text-blue-500" :size="32"><Loading /></el-icon>
        </div>
        <p class="text-gray-500">{{ t('userBadges.loading') }}</p>
      </div>

      <div v-if="!loading && badgeProgress.length === 0" class="text-center py-16">
        <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-gray-100 mb-4">
          <svg class="w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
          </svg>
        </div>
        <p class="text-gray-500 text-lg">{{ t('userBadges.noBadges') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { getErrorI18nKey } from '@/utils/error'
import api from '@/api'
import SvgBadge from '@/components/SvgBadge.vue'

const { t } = useI18n()
const route = useRoute()
const userId = route.params.id
const user = ref(null)
const badgeProgress = ref([])
const loading = ref(false)

function getTypeName(type) {
  const types = {
    basic: t('userBadges.basic'),
    advanced: t('userBadges.advanced'),
    top: t('userBadges.top')
  }
  return types[type] || type
}

function getTypeColor(type) {
  const colors = {
    basic: 'info',
    advanced: 'warning',
    top: 'danger'
  }
  return colors[type] || 'info'
}

function getProgressLabel(key) {
  const labels = {
    register_days: t('userBadges.registerDays'),
    topic_count: t('userBadges.topicCount'),
    like_count: t('userBadges.likeCount'),
    best_comment: t('userBadges.bestComment')
  }
  return labels[key] || key
}

function formatDateTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleDateString()
}

async function loadUser() {
  try {
    user.value = await api.get(`/users/${userId}`)
  } catch (e) {
    console.error(t('userBadges.loadFailed'), e)
    if (e.code) {
      ElMessage.error(t(getErrorI18nKey(e.code)))
    }
  }
}

async function loadBadgeProgress() {
  loading.value = true
  try {
    // 获取指定用户的勋章列表
    const res = await api.get(`/users/${userId}/badges`)
    // 转换数据格式：UserBadge -> badgeProgress 格式
    // 只显示未撤销的勋章，并且添加 awarded 字段
    badgeProgress.value = (res || []).map(ub => ({
      badge: ub.badge,
      awarded: !!ub.awarded_at,
      awarded_at: ub.awarded_at
    }))
  } catch (e) {
    console.error(t('userBadges.loadFailed'), e)
    if (e.code) {
      ElMessage.error(t(getErrorI18nKey(e.code)))
    } else {
      ElMessage.error(t('userBadges.loadFailed'))
    }
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await Promise.all([
    loadUser(),
    loadBadgeProgress()
  ])
})
</script>

<style scoped>
.grayscale {
  filter: grayscale(100%);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
