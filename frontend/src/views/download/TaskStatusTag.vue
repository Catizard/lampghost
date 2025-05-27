<template>
  <!-- TODO: I simply cannot find a better solution -->
  <template v-if="status == DownloadTaskStatus.ERROR">
    <n-tooltip trigger="hover">
      <template #trigger>
        <n-tag :type="intoType()">
          {{ status }}
        </n-tag>
      </template>
      {{ errorMsg }}
    </n-tooltip>
  </template>
  <template v-else>
    <n-tag :type="intoType()">
      {{ status }}
    </n-tag>
  </template>
</template>

<script setup lang="ts">
import { DownloadTaskStatus } from '@/constants/downloadTaskStatus';

const props = defineProps<{
  status: DownloadTaskStatus,
  errorMsg?: string,
}>();

function intoType(): "info" | "error" {
  if (props.status == DownloadTaskStatus.ERROR) {
    return "error";
  }
  return "info"
}
</script>