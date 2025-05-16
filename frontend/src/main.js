import { createApp } from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import router from './router';
import VueApexCharts from 'vue3-apexcharts';
import { createI18n } from 'vue-i18n';
import { PerfectScrollbarPlugin } from 'vue3-perfect-scrollbar';
import 'vue3-perfect-scrollbar/style.css';
import VueCalendarHeatmap from 'vue3-calendar-heatmap';
import { createPinia } from 'pinia';

const app = createApp(App)

const i18n = createI18n({
  legacy: false,
  locale: "en",
  fallbackLocale: "en",
});

const pinia = createPinia();

app
  .use(naive)
  .use(router)
  .use(VueApexCharts)
  .use(PerfectScrollbarPlugin, {
    componentName: "PerfectScrollbar"
  })
  .use(VueCalendarHeatmap)
  .use(i18n)
  .use(pinia)
  .mount('#app')
