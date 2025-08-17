<template>
  <n-flex vertical align="center">
    <n-text strong>{{ clearType.clearText }}</n-text>
    <!-- <small v-if="scoreOptionValue != null && clearTypeValue != ClearType.NO_PLAY">{{ scoreOption.text }}</small> -->
    <small v-if="bestRecordTimestamp != null && clearTypeValue != ClearType.NO_PLAY">{{
      dayjs(bestRecordTimestamp * 1000).format('YYYY-MM-DD') }}
    </small>
    <small v-else>
      -
    </small>
  </n-flex>
</template>

<script setup lang="ts">
import { ClearType, ClearTypeDef, queryClearTypeColorStyle } from '@/constants/cleartype';
import { queryScoreOptionColorStyle, ScoreOptionDef } from '@/constants/scoreOption';
import dayjs from 'dayjs';
import { computed } from 'vue';

const { clearType: clearTypeValue, scoreOption: scoreOptionValue } = defineProps<{
  clearType: number,
  scoreOption?: number,
  bestRecordTimestamp?: number
}>();

const clearType = computed<ClearTypeDef>(() => {
  return queryClearTypeColorStyle(clearTypeValue);
});

const scoreOption = computed<ScoreOptionDef>(() => {
  return queryScoreOptionColorStyle(scoreOptionValue);
});
</script>
