import { createApp } from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import router from './router';
import VueApexCharts from 'vue3-apexcharts';
import { PerfectScrollbarPlugin } from 'vue3-perfect-scrollbar';
import 'vue3-perfect-scrollbar/style.css';
import VueCalendarHeatmap from 'vue3-calendar-heatmap';
import { createPinia } from 'pinia';
import i18n from './i18n';

const app = createApp(App)

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
