import { type RouteRecordRaw } from "vue-router";

const mainRoutes: RouteRecordRaw[] = [
  {
    name: "home",
    path: "/home",
    component: () => import("../views/home/Home.vue"),
  },
  {
    path: "/rivals/management",
    component: () => import("../views/rivals/RivalManagement.vue"),
  },
  {
    path: "/rivals/tags",
    component: () => import("../views/rivals/RivalTags.vue"),
  },
  {
    path: "/difftable/management",
    component: () => import("../views/difficult_table/DifficultTable.vue"),
  },
  {
    path: "/difftable/scores",
    component: () => import("../views/difficult_table/DifficultScores.vue"),
  },
  {
    path: "/customtable/management",
    component: () => import("../views/custom_table/CustomTableManagement.vue"),
  },
  {
    path: "/customtable/design",
    component: () => import("../views/custom_table/CustomTableDesign.vue"),
  },
  {
    name: "folder",
    path: "/folder",
    component: () => import("../views/folder/FolderManagement.vue"),
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
    name: "download",
    path: "/download",
    component: () => import("../views/download/Download.vue"), 
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
    component: () => import("../views/initialize/initialize.vue"),
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
