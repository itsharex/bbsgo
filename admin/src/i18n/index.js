// i18n 配置文件
import { createI18n } from 'vue-i18n'
import zh from './zh'
import en from './en'

// 获取浏览器语言
function getBrowserLocale() {
  const nav = window.navigator
  const browserLocale = nav.language || nav.userLanguage || 'zh'
  // 只取前两位，如 zh、en
  const shortLocale = browserLocale.substring(0, 2).toLowerCase()
  return shortLocale === 'en' ? 'en' : 'zh'
}

// 从 localStorage 获取保存的语言设置
function getSavedLocale() {
  return localStorage.getItem('admin_locale')
}

const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: getSavedLocale() || getBrowserLocale(),
  fallbackLocale: 'zh',
  messages: {
    zh,
    en
  }
})

export default i18n
