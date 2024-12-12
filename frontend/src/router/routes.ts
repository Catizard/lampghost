import { type RouteRecordRaw } from "vue-router";

const mainRoutes: RouteRecordRaw[] = [
    {
        name: "home",
        path: "/home",
        component: () => import("../views/Home.vue"),
    },
    {
        name: "rivals",
        path: "/rivals",
        component: () => import("../views/rivals.vue"),
    },
    {
        name: "difftable",
        path: "/difftable",
        component: () => import("../views/DifficultTable.vue")
    }
];

const routes: RouteRecordRaw[] = [
    {
        name: "not-found",
        path: "/:path*",
        component: () => import("../views/error.vue"),
    },
    {
        name: "layout",
        path: "/",
        component: () => import("../layout/index.vue"),
        children: [...mainRoutes],
    }
]

export default routes;
