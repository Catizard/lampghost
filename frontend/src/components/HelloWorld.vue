<template>
  <n-layout has-sider style="height: 100%;">
    <n-layout-sider bordered show-trigger collapse-mode="width" :collapsed-width="64" :width="240"
      :native-scrollbar="false" :inverted="inverted">
      <n-menu :inverted="inverted" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions" />
    </n-layout-sider>
    <n-layout>
      1234
    </n-layout>
  </n-layout>
</template>

<script lang="ts">
import type { Component } from 'vue'
import {
  BookOutline as BookIcon,
  PersonOutline as PersonIcon,
  WineOutline as WineIcon,
} from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'
import { defineComponent, h, ref } from 'vue'
import { RouterView } from 'vue-router'

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

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

export default defineComponent({
  setup() {
    return {
      inverted: ref(false),
      menuOptions
    }
  }
})
</script>