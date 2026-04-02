<template>
  <div class="bg-white rounded-lg shadow-sm p-6">
    <h1 class="text-2xl font-bold text-gray-900 mb-6">发表新帖</h1>
    <form @submit.prevent="handleSubmit">
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-medium mb-2">标题</label>
        <input type="text" v-model="form.title"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500" placeholder="请输入标题"
          required>
      </div>
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-medium mb-2">版块 <span class="text-red-500">*</span></label>
        <select v-model="form.forum_id"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500" required>
          <option value="">请选择版块</option>
          <option v-for="forum in forums" :key="forum.id" :value="forum.id">{{ forum.name }}</option>
        </select>
      </div>
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-medium mb-2">话题（可选，最多3个）</label>
        <div class="relative">
          <div class="flex flex-wrap gap-2 mb-2">
            <span v-for="(tag, index) in selectedTags" :key="index"
              class="inline-flex items-center px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm">
              {{ tag }}
              <button type="button" @click="removeTag(index)" class="ml-2 text-blue-500 hover:text-blue-700">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
              </button>
            </span>
          </div>
          <input type="text" v-model="tagInput" @input="searchTags" @keydown.enter.prevent="addTag"
            @keydown.down="navigateSuggestion(1)" @keydown.up="navigateSuggestion(-1)" @keydown.escape="showSuggestions = false"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
            placeholder="输入话题名称，按回车添加"
            :disabled="selectedTags.length >= 3">
          <div v-if="showSuggestions && suggestions.length > 0"
            class="absolute z-10 w-full mt-1 bg-white border rounded-lg shadow-lg max-h-48 overflow-y-auto">
            <div v-for="(suggestion, index) in suggestions" :key="suggestion.id"
              @click="selectSuggestion(suggestion.name)"
              :class="['px-4 py-2 cursor-pointer', index === suggestionIndex ? 'bg-blue-50' : 'hover:bg-gray-50']">
              <span class="font-medium">{{ suggestion.name }}</span>
              <span class="text-xs text-gray-400 ml-2">{{ suggestion.usage_count }}次使用</span>
            </div>
          </div>
        </div>
        <p class="text-xs text-gray-500 mt-1">话题由用户自由创建，2-20个字符</p>
      </div>
      <div class="mb-6">
        <label class="block text-gray-700 text-sm font-medium mb-2">内容 <span class="text-red-500">*</span></label>
        <textarea v-model="form.content" rows="12"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500" placeholder="请输入内容"
          required></textarea>
      </div>
      <div class="flex justify-end space-x-4">
        <button type="button" @click="$router.back()" class="px-6 py-2 border rounded-lg hover:bg-gray-50">取消</button>
        <button type="submit" class="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600">发布</button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'

const router = useRouter()
const form = ref({
  title: '',
  content: '',
  forum_id: null,
  tag_names: []
})
const forums = ref([])
const selectedTags = ref([])
const tagInput = ref('')
const suggestions = ref([])
const showSuggestions = ref(false)
const suggestionIndex = ref(-1)
let searchTimeout = null

async function loadForums() {
  try {
    const res = await api.get('/forums')
    forums.value = res || []
  } catch (e) {
    console.error(e)
  }
}

function searchTags() {
  if (searchTimeout) clearTimeout(searchTimeout)
  
  if (!tagInput.value.trim()) {
    suggestions.value = []
    showSuggestions.value = false
    return
  }
  
  searchTimeout = setTimeout(async () => {
    try {
      const res = await api.get('/tags/search', { params: { q: tagInput.value.trim() } })
      suggestions.value = res || []
      showSuggestions.value = suggestions.value.length > 0
      suggestionIndex.value = -1
    } catch (e) {
      console.error(e)
    }
  }, 300)
}

function addTag() {
  const tagName = tagInput.value.trim()
  if (!tagName) return
  
  if (tagName.length < 2 || tagName.length > 20) {
    return
  }
  
  if (selectedTags.value.includes(tagName)) {
    tagInput.value = ''
    showSuggestions.value = false
    return
  }
  
  if (selectedTags.value.length >= 3) {
    return
  }
  
  selectedTags.value.push(tagName)
  form.value.tag_names = selectedTags.value
  tagInput.value = ''
  showSuggestions.value = false
}

function selectSuggestion(name) {
  if (selectedTags.value.includes(name)) {
    tagInput.value = ''
    showSuggestions.value = false
    return
  }
  
  if (selectedTags.value.length >= 3) return
  
  selectedTags.value.push(name)
  form.value.tag_names = selectedTags.value
  tagInput.value = ''
  showSuggestions.value = false
}

function removeTag(index) {
  selectedTags.value.splice(index, 1)
  form.value.tag_names = selectedTags.value
}

function navigateSuggestion(direction) {
  if (!showSuggestions.value || suggestions.value.length === 0) return
  
  const newIndex = suggestionIndex.value + direction
  if (newIndex >= 0 && newIndex < suggestions.value.length) {
    suggestionIndex.value = newIndex
  }
}

async function handleSubmit() {
  try {
    const res = await api.post('/topics', form.value)
    router.push(`/topic/${res.id}`)
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadForums()
})
</script>
