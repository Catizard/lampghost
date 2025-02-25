<template>
	<n-layout-sider bordered show-trigger collapse-mode="width" :collapsed-width="64" :width="190"
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
		label: renderOption("/home", "home"),
		key: 'home',
		icon: renderIcon(PersonIcon)
	},
	{
		label: "rivals",
		key: "rivals",
		icon: renderIcon(RivalIcon),
		children: [
			{
				label: renderOption("/rivals/management", "management"),
				key: "rival_management",
				icon: renderIcon(SearchIcon),
			},
			{
				label: renderOption("/rivals/tags", "tags"),
				key: "rival_tags",
				icon: renderIcon(MedalIcon),
			}
		],
	},
	{
		label: "difftable",
		key: "difftable",
		icon: renderIcon(BookIcon),
		children: [
			{
				label: renderOption("/difftable/management", "management"),
				key: "difftable_management",
				icon: renderIcon(SearchIcon),
			},
			{
				label: renderOption("/difftable/scores", "scores"),
				key: "difftable_scores",
				icon: renderIcon(BalloonIcon)
			}
		]
	},
	{
		label: renderOption("/folder", "folder"),
		key: "folder",
		icon: renderIcon(BookmarksIcon)
	},
	{
		label: renderOption("/recent", "recent"),
		key: "recent",
		icon: renderIcon(TimeIcon)
	},
	{
		label: renderOption("/courses", "courses"),
		key: "courses",
		icon: renderIcon(CoursesIcon)
	},
	{
		label: renderOption("/goals", "goals"),
		key: "goals",
		icon: renderIcon(GoalsIcon),
	},
	{
		label: renderOption("/settings", "settings"),
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
