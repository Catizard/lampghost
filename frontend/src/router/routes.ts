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
    component: () => import("../views/difficult_table/DifficultTable.vue"),
  },
  {
    name: "folder",
    path: "/folder",
    component: () => import("../views/folder/folder.vue"),
  },
  {
    name: "recent",
    path: "/recent",
    component: () => import("../views/recent/Recent.vue"),
  },
  {
    path: "/",
    redirect: "/home",
  },
];

const routes: RouteRecordRaw[] = [
  {
    name: "initialize",
    path: "/initialize",
    component: () => import("../views/initialize.vue"),
  },
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
  },
];

export default routes;
