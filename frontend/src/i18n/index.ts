import { createI18n } from "vue-i18n";
import en from './en.json';
import zhCN from './zh-CN.json';

const i18n = createI18n({
  legacy: false,
  locale: "en",
  fallbackLocale: "en",
  globalInjection: true,
  messages: {
    en,
    zhCN,
  }
});

export default i18n;
