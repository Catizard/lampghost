import { createApp } from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import router from './router';
import VueApexCharts from 'vue3-apexcharts';
import { createI18n } from 'vue-i18n';
import { PerfectScrollbarPlugin } from 'vue3-perfect-scrollbar';
import 'vue3-perfect-scrollbar/style.css';

const app = createApp(App)

app.use(naive)
app.use(router)
app.use(VueApexCharts);
app.use(PerfectScrollbarPlugin, {
	componentName: "PerfectScrollbar"
})

const i18n = createI18n({
	legacy: false,
	locale: "en",
	fallbackLocale: "en",
});

app.use(i18n);

app.mount('#app')
