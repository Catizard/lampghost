import {createApp} from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import router from './router';
import VueApexCharts from 'vue3-apexcharts';


const app = createApp(App)

app.use(naive)
app.use(router)
app.use(VueApexCharts);

router.replace("/home");

app.mount('#app')