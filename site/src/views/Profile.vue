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
                <div class="text-2xl font-bold text-gray-700">643</div>
                <div class="text-xs text-gray-400">帖子</div>
              </div>
              <div>
                <div class="text-2xl font-bold text-gray-700">1167</div>
                <div class="text-xs text-gray-400">评论</div>
              </div>
              <div>
                <div class="text-2xl font-bold text-gray-700">1</div>
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
                <span class="text-gray-900 text-sm">{{ user?.username }}</span>
              </div>
              <div class="flex">
                <span class="w-20 text-gray-500 text-sm">签名</span>
                <span class="text-gray-900 text-sm">{{ user?.signature || '-' }}</span>
              </div>
              <div class="flex">
                <span class="w-20 text-gray-500 text-sm">主页</span>
                <span class="text-blue-500 text-sm">{{ user?.intro ? user.intro : 'https://mlog.club/user/' + (user?.id
                  || '') }}</span>
              </div>
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-4">
            <div class="flex justify-between items-center mb-4">
              <h3 class="font-medium text-gray-900">粉丝 34</h3>
              <button class="text-blue-500 text-sm hover:underline">更多</button>
            </div>
            <div class="space-y-3">
              <div v-for="follower in followers" :key="follower.id" class="flex items-center space-x-3">
                <img :src="follower.avatar || 'https://via.placeholder.com/40'" class="w-10 h-10 rounded-full">
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-medium text-gray-900 truncate">{{ follower.username }}</div>
                  <div class="text-xs text-gray-400 truncate">{{ follower.signature || '这家伙很懒，什么都没留下' }}</div>
                </div>
                <button class="bg-blue-500 text-white text-xs px-3 py-1 rounded hover:bg-blue-600">+ 关注</button>
              </div>
            </div>
          </div>
        </div>
      </aside>
      <div class="flex-1 min-w-0">
        <div class="bg-white rounded-lg shadow-sm">
          <div class="p-4">
            <div v-if="currentTab === 'topics'" class="space-y-4">
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
                  <span>心赞 4</span>
                  <span>评论 3</span>
                  <span>浏览 149</span>
                  <span v-if="topic.forum" class="bg-gray-100 px-2 py-0.5 rounded">{{ topic.forum.name }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'

const route = useRoute()
const user = ref(null)
const currentTab = ref('topics')
const userTopics = ref([])

const followers = ref([
  { id: 1, username: 'user/profile', avatar: '', signature: '这家伙很懒，什么都没留下' },
  { id: 2, username: 'abc123', avatar: '', signature: '这家伙很懒，什么都没留下' },
  { id: 3, username: 'geraint', avatar: '', signature: 'boyboyboy' },
  { id: 4, username: 'Ame', avatar: '', signature: '这家伙很懒，什么都没留下' },
  { id: 5, username: 'wolfCoder', avatar: '', signature: '这家伙很懒，什么都没留下' },
  { id: 6, username: '哈哈', avatar: '', signature: '这家伙很懒，什么都没留下' },
  { id: 7, username: '糖果屋里', avatar: '', signature: '糖果屋里' }
])

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
    user.value = {
      id: 1,
      username: '小码哥',
      avatar: 'https://via.placeholder.com/128',
      signature: '一个爱折腾的老码农。',
      credits: 3098,
      level: 10
    }
    userTopics.value = [
      {
        id: 1,
        title: '热心的朋友问："bbs-go最近有没有计划更新新的功能"，我的回答："重构！"',
        content: '为什么要不断的进行重构呢？我认为主要有以下几点：作为一个开源程序，大家关注他最主要的目的是为了交流和学习，所以要确保自己的代码与时俱进（至少代码不能太丑）。之前的架构并不一定是最优的，随着技术的进步会觉得之前的代码可能换一种方式实现更加美观、合理...',
        user: { username: '小码哥' },
        forum: { name: '交流' },
        created_at: new Date(Date.now() - 6 * 3600000).toISOString()
      },
      {
        id: 2,
        title: 'bbs-go v3.5.0 发布，升级go1.18',
        content: '文档地址: 帮助文档: https://docs.bbs-go.com/ 官网交流: https://mlog.club 问题反馈: https://mlog.club/topic/node/3 功能建议收集: https://mlog.club/topic/60...',
        user: { username: '小码哥' },
        forum: { name: '开源' },
        created_at: new Date(Date.now() - 7 * 3600000).toISOString()
      },
      {
        id: 3,
        title: '构建 Go 应用 docker 镜像的十八种姿势',
        content: '转载自: https://mp.weixin.qq.com/s/cJcOsCDL_XHG4QpcRWIqYg 夜以继日，加班加点开发了一个最简单的 Go Hello world 应用 通宵熬夜，夜以继日，加班加点开发了一个最简单的 Go Hello ...',
        user: { username: '小码哥' },
        forum: { name: '分享' },
        created_at: new Date(Date.now() - 13 * 24 * 3600000).toISOString()
      },
      {
        id: 4,
        title: '我计划将bbs-go的服务端接口修改为由Java实现，各位觉得怎么样。',
        content: '',
        user: { username: '小码哥' },
        forum: { name: '交流' },
        created_at: new Date(Date.now() - 15 * 24 * 3600000).toISOString()
      }
    ]
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadUser()
})
</script>
