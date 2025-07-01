import { createI18n } from "vue-i18n";
import en from './locales/en.json';
import zhCN from './locales/zh-CN.json';

const i18n = createI18n({
  legacy: false,
  locale: "en",
  fallbackLocale: "en",
  globalInjection: true,
  silentFallbackWarn: true, // seems not working
  messages: {
    en,
    zhCN,
  }
});

export default i18n;
