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
	FishOutline as RivalIcon,
	SchoolOutline as CoursesIcon,
	ListOutline as GoalsIcon,
	Search as SearchIcon,
	BalloonOutline as BalloonIcon,
	MedalOutline as MedalIcon,
} from '@vicons/ionicons5'
import { MenuOption, NIcon } from 'naive-ui'
import { computed, h, ref, watchEffect } from 'vue'
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
				icon: renderIcon(SearchIcon),
			},
			{
				label: renderOption("/rivals/tags", t('menuName.rivals.tags')),
				key: "rival_tags",
				icon: renderIcon(MedalIcon),
			}
		],
	},
	{
		label: t('menuName.difftable.self'),
		key: "difftable",
		icon: renderIcon(BookIcon),
		children: [
			{
				label: renderOption("/difftable/management", t('menuName.difftable.management')),
				key: "difftable_management",
				icon: renderIcon(SearchIcon),
			},
			{
				label: renderOption("/difftable/scores", t('menuName.difftable.scores')),
				key: "difftable_scores",
				icon: renderIcon(BalloonIcon)
			}
		]
	},
	{
		label: renderOption("/folder", t('menuName.folder')),
		key: "folder",
		icon: renderIcon(BookmarksIcon)
	},
	{
		label: renderOption("/recent", t('menuName.recent')),
		key: "recent",
		icon: renderIcon(TimeIcon)
	},
	{
		label: renderOption("/courses", t('menuName.courses')),
		key: "courses",
		icon: renderIcon(CoursesIcon)
	},
	{
		label: renderOption("/goals", t('menuName.goals')),
		key: "goals",
		icon: renderIcon(GoalsIcon),
	},
	{
		label: renderOption("/settings", t('menuName.settings')),
		key: "settings",
		icon: renderIcon(SettingsIcon),
	}
]
</script>

<i18n>
{
	"en": {
    "menuName": {
      "home": "home",
      "rivals": {
        "self": "rivals",
        "management": "management",
        "tags": "tags",
      },
      "difftable": {
        "self": "table",
        "management": "management",
        "scores": "scores",
      },
      "folder": "folder",
      "recent": "recent",
      "courses": "courses",
      "goals": "goals",
			"settings": "settings",
    },
  },
  "zh-CN": {
    "menuName": {
      "home": "个人主页",
      "rivals": {
        "self": "好友列表",
        "management": "管理",
        "tags": "标签",
      },
      "difftable": {
        "self": "难度表",
        "management": "管理",
        "scores": "得分",
      },
      "folder": "收藏夹",
      "recent": "最近游玩",
      "courses": "段位列表",
      "goals": "目标列表",
			"settings": "设置",
    },
  },
}
</i18n>

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
