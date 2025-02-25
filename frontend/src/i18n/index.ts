import { createI18n } from "vue-i18n";

const i18n = createI18n({
  legacy: false,
  locale: "en",
  fallbackLocale: "en",
  en: {
    menuName: {
      home: "home",
      rivals: {
        self: "rivals",
        management: "management",
        tags: "tags",
      },
      difftable: {
        self: "table",
        management: "management",
        scores: "scores",
      },
      folder: "folder",
      recent: "recent",
      courses: "courses",
      goals: "goals",
    },
  },
  "zh-CN": {
    menuName: {
      home: "个人主页",
      rivals: {
        self: "好友列表",
        management: "管理",
        tags: "标记",
      },
      difftable: {
        self: "难度表",
        management: "管理",
        scores: "得分",
      },
      folder: "收藏夹",
      recent: "最近游玩",
      courses: "段位列表",
      goals: "目标列表",
    },
  },
});

export default i18n;
