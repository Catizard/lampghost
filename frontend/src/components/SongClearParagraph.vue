<template>
  <n-flex vertical align="center" :size="6">
    <strong style="font-size: 1.1em;">{{ clearType.clearText }}</strong>
    <!-- <small v-if="scoreOptionValue != null && clearTypeValue != ClearType.NO_PLAY">{{ scoreOption.text }}</small> -->
    <template v-if="!disableTimestamp">
      <small v-if="bestRecordTimestamp != null && clearTypeValue != ClearType.NO_PLAY">{{
        dayjs(bestRecordTimestamp * 1000).format('YYYY-MM-DD') }}
      </small>
      <small v-else>
        -
      </small>
    </template>
  </n-flex>
</template>

<script setup lang="ts">
import { ClearType, ClearTypeDef, queryClearTypeColorStyle } from '@/constants/cleartype';
import { queryScoreOptionColorStyle, ScoreOptionDef } from '@/constants/scoreOption';
import dayjs from 'dayjs';
import { computed } from 'vue';

const { clearType: clearTypeValue, scoreOption: scoreOptionValue, bestRecordTimestamp, disableTimestamp = false } = defineProps<{
  clearType: number,
  scoreOption?: number,
  bestRecordTimestamp?: number
  disableTimestamp?: boolean
}>();

const clearType = computed<ClearTypeDef>(() => {
  return queryClearTypeColorStyle(clearTypeValue);
});

const scoreOption = computed<ScoreOptionDef>(() => {
  return queryScoreOptionColorStyle(scoreOptionValue);
});
</script>

<style lang="css" scoped>
small {
  font-size: 0.85em;
}
</style>
