<template>
  <template v-if="noplay">/</template>
  <template v-else>
    <n-flex vertical :size="4">
      <strong style="font-size: 1.1em;">{{ normalizeMinbp }}</strong>
      <n-text v-if="!noOption" depth="2" style="font-size: 0.8em;">{{ optionName }}</n-text>
    </n-flex>
  </template>

</template>

<script setup lang="ts">
import { queryScoreOptionColorStyle } from '@/constants/scoreOption';
import { computed } from 'vue';

const { noplay = false, optionValue, noOption = false, minbp } = defineProps<{
  noplay?: boolean,
  optionValue?: number
  noOption?: boolean
  minbp: number
}>();

const optionName = computed<string>(() => {
  return queryScoreOptionColorStyle(optionValue).text;
});

const normalizeMinbp = computed<number>(() => {
  return minbp == 2147483647 ? 0 : minbp
});

</script>
