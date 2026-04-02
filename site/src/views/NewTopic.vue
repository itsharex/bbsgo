<template>
  <div class="max-w-4xl mx-auto bg-white rounded-lg shadow-sm p-6">
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
            @keydown.down="navigateSuggestion(1)" @keydown.up="navigateSuggestion(-1)"
            @keydown.escape="showSuggestions = false"
            class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
            placeholder="输入话题名称，按回车添加" :disabled="selectedTags.length >= 3">
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

        <div class="border rounded-lg overflow-hidden">
          <div class="toolbar bg-gray-50 border-b p-2 flex flex-wrap gap-1 items-center">
            <button type="button" @click="editorMode = 'wysiwyg'"
              :class="['px-3 py-1.5 text-sm rounded-lg transition-colors', editorMode === 'wysiwyg' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200']">
              富文本
            </button>
            <button type="button" @click="editorMode = 'markdown'"
              :class="['px-3 py-1.5 text-sm rounded-lg transition-colors', editorMode === 'markdown' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200']">
              Markdown
            </button>
            <button type="button" @click="editorMode = 'preview'"
              :class="['px-3 py-1.5 text-sm rounded-lg transition-colors', editorMode === 'preview' ? 'bg-blue-500 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200']">
              预览
            </button>
            
            <template v-if="editorMode === 'wysiwyg'">
              <div class="w-px h-8 bg-gray-300 mx-1"></div>
              <button type="button" @click="toggleBold"
                :class="['p-2 rounded hover:bg-gray-200', isActive('bold') && 'bg-blue-100 text-blue-600']" title="粗体">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M6 4h8a4 4 0 014 4 4 4 0 01-4 4H6z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M6 12h9a4 4 0 014 4 4 4 0 01-4 4H6z"></path>
                </svg>
              </button>
              <button type="button" @click="toggleItalic"
                :class="['p-2 rounded hover:bg-gray-200', isActive('italic') && 'bg-blue-100 text-blue-600']" title="斜体">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 4h4m-2 0v16m-4 0h8"></path>
                </svg>
              </button>
              <button type="button" @click="toggleUnderline"
                :class="['p-2 rounded hover:bg-gray-200', isActive('underline') && 'bg-blue-100 text-blue-600']"
                title="下划线">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8v8a5 5 0 0010 0V8M5 21h14">
                  </path>
                </svg>
              </button>
              <div class="w-px h-8 bg-gray-300 mx-1"></div>
              <button type="button" @click="toggleHeading(1)"
                :class="['p-2 rounded hover:bg-gray-200', isActiveHeading(1) && 'bg-blue-100 text-blue-600']" title="标题1">
                <span class="font-bold text-lg">H1</span>
              </button>
              <button type="button" @click="toggleHeading(2)"
                :class="['p-2 rounded hover:bg-gray-200', isActiveHeading(2) && 'bg-blue-100 text-blue-600']" title="标题2">
                <span class="font-bold">H2</span>
              </button>
              <button type="button" @click="toggleHeading(3)"
                :class="['p-2 rounded hover:bg-gray-200', isActiveHeading(3) && 'bg-blue-100 text-blue-600']" title="标题3">
                <span class="font-bold text-sm">H3</span>
              </button>
              <div class="w-px h-8 bg-gray-300 mx-1"></div>
              <button type="button" @click="toggleBulletList"
                :class="['p-2 rounded hover:bg-gray-200', isActive('bulletList') && 'bg-blue-100 text-blue-600']"
                title="无序列表">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16">
                  </path>
                </svg>
              </button>
              <button type="button" @click="toggleOrderedList"
                :class="['p-2 rounded hover:bg-gray-200', isActive('orderedList') && 'bg-blue-100 text-blue-600']"
                title="有序列表">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14"></path>
                </svg>
              </button>
              <div class="w-px h-8 bg-gray-300 mx-1"></div>
              <button type="button" @click="toggleBlockquote"
                :class="['p-2 rounded hover:bg-gray-200', isActive('blockquote') && 'bg-blue-100 text-blue-600']"
                title="引用">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                  </path>
                </svg>
              </button>
              <button type="button" @click="toggleCodeBlock"
                :class="['p-2 rounded hover:bg-gray-200', isActive('codeBlock') && 'bg-blue-100 text-blue-600']"
                title="代码块">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
                </svg>
              </button>
              <div class="w-px h-8 bg-gray-300 mx-1"></div>
              <button type="button" @click="addImage" :class="['p-2 rounded hover:bg-gray-200']" title="插入图片">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z">
                  </path>
                </svg>
              </button>
              <button type="button" @click="addLink" :class="['p-2 rounded hover:bg-gray-200']" title="插入链接">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l1.1-1.1">
                  </path>
                </svg>
              </button>
            </template>
          </div>
          <editor-content v-if="editorMode === 'wysiwyg'" :editor="editor" class="tiptap-editor" />
          <textarea v-else-if="editorMode === 'markdown'" v-model="markdownContent"
            class="w-full min-h-[400px] p-4 font-mono text-sm border-0 focus:outline-none resize-none"
            placeholder="在此输入 Markdown 内容..."></textarea>
          <div v-else-if="editorMode === 'preview'" class="prose prose-sm p-4 min-h-[400px]" v-html="previewHtml"></div>
        </div>
        <p class="text-xs text-gray-500 mt-2">💡 提示：可以直接 Ctrl+V 粘贴图片</p>
      </div>

      <div class="flex justify-end space-x-4">
        <button type="button" @click="$router.back()" class="px-6 py-2 border rounded-lg hover:bg-gray-50">取消</button>
        <button type="submit" :disabled="submitting"
          class="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed">
          {{ submitting ? '发布中...' : '发布' }}
        </button>
      </div>
    </form>

    <input type="file" ref="imageInput" accept="image/*" class="hidden" @change="handleImageSelect">
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import ImageResize from 'tiptap-extension-resize-image'
import Link from '@tiptap/extension-link'
import TextAlign from '@tiptap/extension-text-align'
import Underline from '@tiptap/extension-underline'
import api from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import MarkdownIt from 'markdown-it'
import TurndownService from 'turndown'

const md = new MarkdownIt()
const turndown = new TurndownService()

const router = useRouter()
const imageInput = ref(null)

const editorMode = ref('wysiwyg')
const markdownContent = ref('')

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
const submitting = ref(false)
let searchTimeout = null

const previewHtml = computed(() => {
  return md.render(markdownContent.value)
})

const editor = useEditor({
  extensions: [
    StarterKit.configure({
      heading: {
        levels: [1, 2, 3]
      }
    }),
    ImageResize.configure({
      inline: false,
      minWidth: 100,
      maxWidth: 800,
      HTMLAttributes: {
        class: 'rounded-lg'
      }
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: {
        rel: 'noopener noreferrer',
        target: '_blank',
        class: 'text-blue-600 underline hover:text-blue-800'
      }
    }),
    TextAlign.configure({
      types: ['heading', 'paragraph']
    }),
    Underline
  ],
  content: '',
  editorProps: {
    attributes: {
      class: 'prose prose-sm focus:outline-none min-h-[400px] p-4 max-w-none'
    },
    handlePaste: (view, event) => {
      const items = event.clipboardData?.items
      if (items) {
        for (let item of items) {
          if (item.type.indexOf('image') !== -1) {
            const file = item.getAsFile()
            if (file) {
              event.preventDefault()
              uploadImage(file)
              return true
            }
          }
        }
      }
      return false
    }
  },
  onUpdate: ({ editor }) => {
    if (editorMode.value === 'wysiwyg') {
      form.value.content = editor.getHTML()
    }
  }
})

watch(editorMode, (newMode, oldMode) => {
  if (oldMode === 'wysiwyg' && newMode !== 'wysiwyg') {
    if (editor.value) {
      const html = editor.value.getHTML()
      markdownContent.value = turndown.turndown(html)
    }
  } else if (oldMode !== 'wysiwyg' && newMode === 'wysiwyg') {
    if (editor.value) {
      const html = md.render(markdownContent.value)
      editor.value.commands.setContent(html)
    }
  }
})

watch(markdownContent, () => {
  if (editorMode.value === 'markdown' || editorMode.value === 'preview') {
    form.value.content = md.render(markdownContent.value)
  }
})

function isActive(format) {
  if (!editor.value) return false
  return editor.value.isActive(format)
}

function isActiveHeading(level) {
  if (!editor.value) return false
  return editor.value.isActive('heading', { level })
}

function toggleBold() {
  editor.value?.chain().focus().toggleBold().run()
}

function toggleItalic() {
  editor.value?.chain().focus().toggleItalic().run()
}

function toggleUnderline() {
  editor.value?.chain().focus().toggleUnderline().run()
}

function toggleHeading(level) {
  editor.value?.chain().focus().toggleHeading({ level }).run()
}

function toggleBulletList() {
  editor.value?.chain().focus().toggleBulletList().run()
}

function toggleOrderedList() {
  editor.value?.chain().focus().toggleOrderedList().run()
}

function toggleBlockquote() {
  editor.value?.chain().focus().toggleBlockquote().run()
}

function toggleCodeBlock() {
  editor.value?.chain().focus().toggleCodeBlock().run()
}

function addImage() {
  imageInput.value?.click()
}

async function addLink() {
  try {
    const result = await ElMessageBox.prompt('请输入链接地址', '插入链接', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /^https?:\/\/.+/,
      inputErrorMessage: '请输入有效的 URL'
    })
    if (result.value) {
      editor.value?.chain().focus().setLink({ href: result.value }).run()
    }
  } catch {
  }
}

async function handleImageSelect(event) {
  const file = event.target.files[0]
  if (file) {
    await uploadImage(file)
  }
  event.target.value = ''
}

async function uploadImage(file) {
  try {
    const formData = new FormData()
    formData.append('file', file)

    const response = await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    if (response && response.url) {
      if (editorMode.value === 'wysiwyg') {
        editor.value?.chain().focus().setImage({ src: response.url }).run()
      } else {
        markdownContent.value += `\n![图片](${response.url})\n`
      }
    } else {
      throw new Error('上传失败')
    }
  } catch (error) {
    console.error('Image upload error:', error)
    ElMessage.error('图片上传失败')
  }
}

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
  if (!form.value.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }

  if (!form.value.forum_id) {
    ElMessage.warning('请选择版块')
    return
  }

  const contentToCheck = editorMode.value === 'wysiwyg' 
    ? form.value.content 
    : markdownContent.value

  if (!contentToCheck || contentToCheck.trim() === '') {
    ElMessage.warning('请输入内容')
    return
  }

  submitting.value = true

  try {
    const res = await api.post('/topics', form.value)
    ElMessage.success('发布成功')
    setTimeout(() => {
      router.push(`/topic/${res.id}`)
    }, 1500)
  } catch (error) {
    console.error(error)
    ElMessage.error(error.response?.data?.message || '发布失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadForums()
})

onBeforeUnmount(() => {
  if (editor.value) {
    editor.value.destroy()
  }
})
</script>

<style scoped>
.tiptap-editor {
  min-height: 400px;
}

.tiptap-editor :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin: 1rem 0;
}

.tiptap-editor :deep(a) {
  color: #2563eb;
  text-decoration: underline;
}

.tiptap-editor :deep(a:hover) {
  color: #1e40af;
}
</style>
