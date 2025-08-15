<template>
  <n-flex vertical align="center">
    <strong>{{ clearType.clearText }}</strong>
    <small v-if="scoreOptionValue != null && clearTypeValue != ClearType.NO_PLAY">{{ scoreOption.text }}</small>
  </n-flex>
</template>

<script setup lang="ts">
import { ClearType, ClearTypeDef, queryClearTypeColorStyle } from '@/constants/cleartype';
import { queryScoreOptionColorStyle, ScoreOptionDef } from '@/constants/scoreOption';
import { computed } from 'vue';

const { clearType: clearTypeValue, scoreOption: scoreOptionValue } = defineProps<{
  clearType: number,
  scoreOption?: number
}>();

const clearType = computed<ClearTypeDef>(() => {
  return queryClearTypeColorStyle(clearTypeValue);
});

const scoreOption = computed<ScoreOptionDef>(() => {
  return queryScoreOptionColorStyle(scoreOptionValue);
});
</script>
