<template>
  <n-modal v-model:show="show" preset="dialog" :content-style="{ 'height': '100%' }" style="width: 95vw; height: 95vh;">
    <template #header>
      <n-flex>
        Preview
        <n-button>
          <template #icon>
            <n-icon height="30px" @click="copy">
              <CopyIcon />
            </n-icon>
          </template>
        </n-button>
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
import { Refresh as RefreshIcon, Copy as CopyIcon } from '@vicons/ionicons5';
import { QueryPreviewURLByMd5 } from '@wailsjs/go/main/App';
import { ClipboardSetText } from '@wailsjs/runtime/runtime';
import { nextTick, ref, Ref } from 'vue';

defineExpose({ open });

const show = ref(false);
const url: Ref<string | null> = ref(null);

async function open(md5: string, option: number | null = null) {
  try {
    const result = await QueryPreviewURLByMd5(md5)
    url.value = result.Data;
    if (option != null) {
      url.value = `${url.value}&o=${option}`
    }
    show.value = true;
  } catch (err) {
    window.$notifyError(err)
  }
}

function refresh() {
  const prev = url.value;
  url.value = "";
  nextTick(() => {
    url.value = prev;
  });
}

async function copy() {
  try {
    const result = await ClipboardSetText(url.value);
    if (!result) {
      throw "Failed to copy";
    }
  } catch (err) {
    window.$notifyError(err);
  }
}
</script>
