<template>
  <template v-if="recordTimestamp == null">/</template>
  <template v-else>
    <n-flex vertical :size="4">
      <n-text style="font-size: 1em;">{{ dayjs(recordTimestamp).format("YYYY-MM-DD") }}</n-text>
      <n-text style="font-size: 0.9em;">{{ dayjs(recordTimestamp).format("HH:mm:ss") }}</n-text>
    </n-flex>
  </template>
</template>

<script setup lang="ts">
import dayjs from 'dayjs';

const { recordTimestamp } = defineProps<{
  recordTimestamp?: number
}>();
</script>
