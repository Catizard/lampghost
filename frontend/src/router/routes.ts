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
    path: "/difftable/management",
    component: () => import("../views/difficult_table/DifficultTable.vue"),
  },
  {
    path: "/difftable/scores",
    component: () => import("../views/difficult_table/DifficultScores.vue")
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
    name: "settings",
    path: "/settings",
    component: () => import("../views/settings/Settings.vue"),
  },
  {
    name: "courses",
    path: "/courses",
    component: () => import("../views/courses/Courses.vue"),
  },
  {
    name: "goals",
    path: "/goals",
    component: () => import("../views/goals/Goals.vue"),
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
