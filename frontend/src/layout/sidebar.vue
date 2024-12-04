<template>
    <n-layout-sider bordered show-trigger collapse-mode="width" :collapsed-width="64" :width="240"
        :native-scrollbar="false" :inverted="inverted">
        <n-menu :inverted="inverted" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions" />
    </n-layout-sider>
</template>

<script lang="ts" setup>
import type { Component } from 'vue'
import {
    BookOutline as BookIcon,
    PersonOutline as PersonIcon,
    WineOutline as WineIcon,
} from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'
import { computed, h, ref, watchEffect } from 'vue'

function renderIcon(icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) })
}

const inverted = ref(false);
const menuOptions = [
    {
        label: 'home',
        key: 'home',
        icon: renderIcon(PersonIcon)
    },
    {
        label: "rivals",
        key: 'rivals',
        icon: renderIcon(WineIcon)
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
</style>
