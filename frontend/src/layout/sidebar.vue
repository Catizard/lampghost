<template>
  <n-layout-sider bordered show-trigger collapse-mode="width" :collapsed-width="64" :width="210"
    :native-scrollbar="false" :inverted="inverted">
    <n-menu :inverted="inverted" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions"
      defaultValue="home" />
  </n-layout-sider>
</template>

<script lang="ts" setup>
import type { Component } from 'vue'
import {
  BookOutline as BookIcon,
  PersonOutline as PersonIcon,
  Bookmarks as BookmarksIcon,
  TimeOutline as TimeIcon,
  SettingsOutline as SettingsIcon,
  SchoolOutline as CoursesIcon,
  StatsChart as StatisticIcon,
  MedalOutline as MedalIcon,
  People as RivalIcon,
  List as ListIcon,
  DownloadOutline as DownloadIcon,
  ColorWand as CustomTableIcon,
  Pencil as PencilIcon,
  SnowOutline as LockIcon,
  MusicalNotes as SongIcon,
  Flag
} from '@vicons/ionicons5'
import { MenuOption, NIcon } from 'naive-ui'
import { h, ref } from 'vue'
import { RouterLink } from 'vue-router';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const { t } = i18n;

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

function renderOption(path: string, name: string) {
  return () => h(
    RouterLink,
    { to: path },
    { default: () => name }
  );
}

const inverted = ref(false);
const menuOptions: MenuOption[] = [
  {
    label: renderOption("/home", t('menuName.home')),
    key: 'home',
    icon: renderIcon(PersonIcon)
  },
  {
    label: t('menuName.rivals.self'),
    key: "rivals",
    icon: renderIcon(RivalIcon),
    children: [
      {
        label: renderOption("/rivals/management", t('menuName.rivals.management')),
        key: "rival_management",
        icon: renderIcon(ListIcon),
      },
      {
        label: renderOption("/rivals/tags", t('menuName.rivals.tags')),
        key: "rival_tags",
        icon: renderIcon(MedalIcon),
      },
      {
        label: renderOption("/rivals/lock", t('menuName.rivals.lock')),
        key: "rival_lock",
        icon: renderIcon(LockIcon)
      }
    ],
  },
  {
    label: renderOption("/song", t('menuName.song')),
    key: "song",
    icon: renderIcon(SongIcon),
  },
  {
    label: t('menuName.difftable.self'),
    key: "difftable",
    icon: renderIcon(BookIcon),
    children: [
      {
        label: renderOption("/difftable/management", t('menuName.difftable.management')),
        key: "difftable_management",
        icon: renderIcon(ListIcon),
      },
      {
        label: renderOption("/difftable/scores", t('menuName.difftable.scores')),
        key: "difftable_scores",
        icon: renderIcon(StatisticIcon)
      },
      {
        label: renderOption("/difftable/courses", t('menuName.difftable.courses')),
        key: "difftable_courses",
        icon: renderIcon(CoursesIcon)
      }
    ]
  },
  {
    label: renderOption("/folder", t('menuName.folder')),
    key: "folder",
    icon: renderIcon(BookmarksIcon)
  },
  {
    label: t('menuName.customtable.self'),
    key: "customtable",
    icon: renderIcon(CustomTableIcon),
    children: [
      {
        label: renderOption("/customtable/management", t('menuName.customtable.management')),
        key: "customtable_management",
        icon: renderIcon(ListIcon)
      },
      {
        label: renderOption("/customtable/design", t('menuName.customtable.design')),
        key: "customtable_design",
        icon: renderIcon(PencilIcon)
      },
      {
        label: renderOption("/customtable/courses", t('menuName.customtable.courses')),
        key: "customtable_courses",
        icon: renderIcon(Flag)
      }
    ]
  },
  {
    label: renderOption("/recent", t('menuName.recent')),
    key: "recent",
    icon: renderIcon(TimeIcon)
  },
  // {
  // 	label: renderOption("/goals", t('menuName.goals')),
  // 	key: "goals",
  // 	icon: renderIcon(GoalsIcon),
  // },
  {
    label: renderOption("/download", t('menuName.download')),
    key: "download",
    icon: renderIcon(DownloadIcon)
  },

  {
    label: renderOption("/settings", t('menuName.settings')),
    key: "settings",
    icon: renderIcon(SettingsIcon),
  }
]
</script>

<style scoped>
.logo {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  padding: 12px 20px;
  background: var(--n-color);
  text-decoration: none;
}

.n-layout-sider--collapsed .logo {
  padding: 9px;
}

.logo svg {
  height: 32px;
}

:deep(.n-menu > .n-menu-item:last-child) {
  position: absolute;
  bottom: 10px;
  width: 100%;
}
</style>
