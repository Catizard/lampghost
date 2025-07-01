<template>
  <n-modal v-model:show="show" preset="dialog" :content-style="{ 'height': '100%' }" style="width: 95vw; height: 95vh;">
    <template #header>
      <n-flex>
        Preview
        <n-button>
          <template #icon>
            <n-icon height="30px" @click="refresh">
              <RefreshIcon />
            </n-icon>
          </template>
        </n-button>
      </n-flex>
    </template>
    <iframe :src="url" width="100%" height="95%" />
  </n-modal>
</template>

<script setup lang="ts">
import { Refresh as RefreshIcon } from '@vicons/ionicons5';
import { nextTick, ref, Ref } from 'vue';

defineExpose({ open });

const show = ref(false);
const url: Ref<string | null> = ref(null);

function open(md5: string) {
  url.value = `https://bms-score-viewer.pages.dev/view?md5=${md5}`
  show.value = true;
}

function refresh() {
  const prev = url.value;
  url.value = "";
  nextTick(() => {
    url.value = prev;
  });
}
</script>
