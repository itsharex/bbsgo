// i18n 配置
import { createI18n } from 'vue-i18n'
import zh from './zh'
import en from './en'

// 获取浏览器语言
function getBrowserLocale() {
  const nav = window.navigator
  const browserLocale = nav.language || nav.userLanguage || 'zh'
  const shortLocale = browserLocale.substring(0, 2).toLowerCase()
  return shortLocale === 'en' ? 'en' : 'zh'
}

// 获取保存的语言设置
function getSavedLocale() {
  return localStorage.getItem('site_locale')
}

const i18n = createI18n({
  legacy: false,
  locale: getSavedLocale() || getBrowserLocale(),
  fallbackLocale: 'zh',
  messages: { zh, en }
})

export default i18n
