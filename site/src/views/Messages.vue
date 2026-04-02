<template>
  <div class="max-w-4xl mx-auto px-4 py-6">
    <div class="bg-white rounded-lg shadow-sm">
      <div class="p-4 border-b">
        <h2 class="text-xl font-bold">私信</h2>
      </div>

      <div class="flex h-[600px]">
        <div class="w-1/3 border-r overflow-y-auto">
          <div class="p-2">
            <div v-if="conversations.length === 0" class="text-center text-gray-500 py-8">
              暂无私信
            </div>
            <div v-for="conv in conversations" :key="conv.user_id" @click="selectConversation(conv)"
              :class="['p-3 cursor-pointer hover:bg-gray-50 rounded-lg', selectedUser?.id === conv.user_id ? 'bg-blue-50' : '']">
              <div class="flex items-center">
                <img :src="conv.user?.avatar || 'https://via.placeholder.com/40'" class="w-10 h-10 rounded-full">
                <div class="ml-3 flex-1 min-w-0">
                  <div class="flex justify-between items-center">
                    <span class="font-medium truncate">{{ conv.user?.username }}</span>
                    <span class="text-xs text-gray-400">{{ formatTime(conv.last_message?.created_at) }}</span>
                  </div>
                  <p class="text-sm text-gray-500 truncate">{{ conv.last_message?.content }}</p>
                </div>
                <span v-if="conv.unread_count > 0" class="ml-2 bg-red-500 text-white text-xs rounded-full px-2 py-0.5">
                  {{ conv.unread_count }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="flex-1 flex flex-col">
          <div v-if="selectedUser" class="flex-1 flex flex-col">
            <div class="p-4 border-b flex items-center">
              <img :src="selectedUser.avatar || 'https://via.placeholder.com/40'" class="w-10 h-10 rounded-full">
              <span class="ml-3 font-medium">{{ selectedUser.username }}</span>
            </div>

            <div ref="messageList" class="flex-1 overflow-y-auto p-4 space-y-4">
              <div v-for="msg in messages" :key="msg.id"
                :class="['flex', msg.from_user_id === currentUserId ? 'justify-end' : 'justify-start']">
                <div :class="['max-w-[70%] rounded-lg px-4 py-2',
                  msg.from_user_id === currentUserId ? 'bg-blue-500 text-white' : 'bg-gray-100']">
                  <p>{{ msg.content }}</p>
                  <span class="text-xs opacity-70">{{ formatTime(msg.created_at) }}</span>
                </div>
              </div>
            </div>

            <div class="p-4 border-t">
              <form @submit.prevent="sendMessage" class="flex space-x-2">
                <input type="text" v-model="newMessage" placeholder="输入消息..."
                  class="flex-1 px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500">
                <button type="submit" class="bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-600">发送</button>
              </form>
            </div>
          </div>

          <div v-else class="flex-1 flex items-center justify-center text-gray-400">
            选择一个对话开始聊天
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import api from '@/api'

const conversations = ref([])
const messages = ref([])
const selectedUser = ref(null)
const newMessage = ref('')
const currentUserId = ref(0)
const messageList = ref(null)

function formatTime(date) {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  if (d.toDateString() === now.toDateString()) {
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  return d.toLocaleDateString('zh-CN')
}

async function loadConversations() {
  try {
    const res = await api.get('/messages')
    const userMap = new Map()

    res.list.forEach(msg => {
      const otherUser = msg.from_user_id === currentUserId.value ? msg.to_user : msg.from_user
      const otherId = msg.from_user_id === currentUserId.value ? msg.to_user_id : msg.from_user_id

      if (!userMap.has(otherId)) {
        userMap.set(otherId, {
          user_id: otherId,
          user: otherUser,
          last_message: msg,
          unread_count: 0
        })
      }

      if (msg.to_user_id === currentUserId.value && !msg.is_read) {
        userMap.get(otherId).unread_count++
      }
    })

    conversations.value = Array.from(userMap.values())
  } catch (e) {
    console.error(e)
  }
}

async function selectConversation(conv) {
  selectedUser.value = conv.user
  await loadMessages(conv.user_id)
  conv.unread_count = 0
}

async function loadMessages(userId) {
  try {
    messages.value = await api.get(`/messages/with/${userId}`)
    await nextTick()
    scrollToBottom()
  } catch (e) {
    console.error(e)
  }
}

async function sendMessage() {
  if (!newMessage.value.trim() || !selectedUser.value) return

  try {
    const msg = await api.post('/messages', {
      to_user_id: selectedUser.value.id,
      content: newMessage.value
    })
    messages.value.push(msg)
    newMessage.value = ''
    await nextTick()
    scrollToBottom()
  } catch (e) {
    console.error(e)
  }
}

function scrollToBottom() {
  if (messageList.value) {
    messageList.value.scrollTop = messageList.value.scrollHeight
  }
}

onMounted(async () => {
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  currentUserId.value = user.id
  await loadConversations()
})
</script>
