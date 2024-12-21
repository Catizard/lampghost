import { createApp } from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import router from './router';
import VueApexCharts from 'vue3-apexcharts';
import { PerfectScrollbarPlugin } from 'vue3-perfect-scrollbar';
import 'vue3-perfect-scrollbar/style.css';

const app = createApp(App)

app.use(naive)
app.use(router)
app.use(VueApexCharts);
app.use(PerfectScrollbarPlugin, {
    componentName: "PerfectScrollbar"
})

app.mount('#app')