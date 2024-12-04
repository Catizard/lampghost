import { createRouter, createWebHashHistory } from 'vue-router';
import HelloWorld from '../views/HelloWorld.vue';

export default createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: "/",
            component: HelloWorld
        }
    ],
});