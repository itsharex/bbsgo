import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import DOMPurify from 'isomorphic-dompurify'
import 'highlight.js/styles/github.css'

// 创建 markdown-it 实例
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return '<pre class="hljs"><code>' +
          hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
          '</code></pre>'
      } catch (__) {}
    }
    return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>'
  }
})

// 渲染 markdown 为 HTML
export function renderMarkdown(content) {
  if (!content) return ''
  const html = md.render(content)
  // 使用 DOMPurify 消毒，防止 XSS 攻击
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: [
      'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'p', 'br', 'hr',
      'ul', 'ol', 'li', 'blockquote', 'pre', 'code', 'span',
      'strong', 'em', 'del', 'sup', 'sub',
      'a', 'img', 'video', 'audio', 'table', 'thead', 'tbody', 'tr', 'th', 'td',
      'div', 'section', 'article'
    ],
    ALLOWED_ATTR: ['href', 'src', 'alt', 'title', 'class', 'target', 'rel', 'controls', 'loop', 'muted'],
    ALLOW_DATA_ATTR: false,
    // 强制所有链接添加 rel="noopener noreferrer" 防止钓鱼
    FORCE_BODY: true,
    ADD_ATTR: ['target']
  })
}

// 去除 markdown 语法，转换为纯文本（用于预览）
export function stripMarkdown(content) {
  if (!content) return ''
  return content
    // 移除 HTML 标签
    .replace(/<[^>]+>/g, '')
    // 移除代码块
    .replace(/```[\s\S]*?```/g, '')
    // 移除行内代码
    .replace(/`[^`]+`/g, '')
    // 移除图片
    .replace(/!\[.*?\]\(.*?\)/g, '')
    // 移除链接，保留文字
    .replace(/\[([^\]]+)\]\(.*?\)/g, '$1')
    // 移除标题标记
    .replace(/^#{1,6}\s+/gm, '')
    // 移除加粗标记
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/__([^_]+)__/g, '$1')
    // 移除斜体标记
    .replace(/\*([^*]+)\*/g, '$1')
    .replace(/_([^_]+)_/g, '$1')
    // 移除引用标记
    .replace(/^>\s+/gm, '')
    // 移除列表标记
    .replace(/^[-*+]\s+/gm, '')
    .replace(/^\d+\.\s+/gm, '')
    // 移除水平线
    .replace(/^[-*_]{3,}$/gm, '')
    // 移除多余空白
    .replace(/\n{3,}/g, '\n\n')
    .trim()
}

// 提取内容中的第一张图片URL
export function extractFirstImage(content) {
  if (!content) return null
  // 匹配 HTML img 标签
  const imgMatch = content.match(/<img[^>]+src=["']([^"']+)["']/i)
  if (imgMatch) return imgMatch[1]
  // 匹配 markdown 图片
  const mdImgMatch = content.match(/!\[.*?\]\((.*?)\)/)
  if (mdImgMatch) return mdImgMatch[1]
  return null
}

// 检查内容是否包含视频
export function hasVideo(content) {
  if (!content) return false
  return /<video[^>]*>|!\[.*?\]\(.*?\.(mp4|webm|ogg)\)/i.test(content)
}
