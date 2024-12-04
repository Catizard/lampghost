import { createRouter, createWebHashHistory } from 'vue-router';
import HelloWorld from '../views/HelloWorld.vue';
import routes from './routes.ts'

export default createRouter({
    history: createWebHashHistory(),
    routes
});
