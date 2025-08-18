<template>
  <n-layout-content ref="contentRef" content-style="padding: 20px 40px; min-height: 85vh">
    <router-view />
  </n-layout-content>
</template>

<script lang="ts" setup>
import { RouterView } from 'vue-router';
import router from "@/router";
import { provide, ref, watch } from 'vue';
import { NLayoutContent } from 'naive-ui';

provide("resetScrollBar", resetScrollBar);

const contentRef = ref<InstanceType<typeof NLayoutContent>>(null);
function resetScrollBar() {
  contentRef?.value?.scrollTo({ top: 0, left: 0 });
}

watch(router.currentRoute, () => {
  resetScrollBar();
});
</script>
