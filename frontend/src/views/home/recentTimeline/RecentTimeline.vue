<template>
  <n-infinite-scroll :loading="loading" style="height: 90%;" :distance="10" @load="handleLoad">
    <TimelineContent :nodes="timelineNodes" />
  </n-infinite-scroll>
</template>

<script lang="ts" setup>
import { QueryPrevDayScoreLogList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import dayjs from 'dayjs';
import { ref, Ref, watch, nextTick } from 'vue';
import { parseOneDayLogs, type Node } from './node';
import TimelineContent from './TimelineContent.vue'

const props = defineProps<{
  rivalId?: number;
}>();

const loading = ref(false);
const timelineNodes: Ref<Array<Node>> = ref([]);
const noMore = ref(false);
let queryDate = dayjs().add(1, 'day');

async function handleLoad() {
  if (loading.value || noMore.value) {
    return;
  }
  loading.value = true;
  try {
    const result = await QueryPrevDayScoreLogList({
      EndRecordTimestamp: queryDate.unix(),
      RivalId: props.rivalId,
    } as any);
    if (result.Code != 200) {
      throw new Error(result.Msg);
    }
    const nextLogs: Array<dto.RivalScoreLogDto> = result.Rows;
    if (nextLogs.length == 0) {
      console.log('no more values')
      noMore.value = true;
    } else {
      console.log('next logs: ', nextLogs);
      const [nextQueryDate, nextNodes] = parseOneDayLogs(nextLogs);
      timelineNodes.value.push(...nextNodes);
      queryDate = nextQueryDate;
      // Hack: To prevent having a very long white space, keep load data if log count < 15
      // There might exist a rare race condition...but I hardly think this is damaging
      if (timelineNodes.value.length < 15) {
        nextTick(() => {
          handleLoad();
        });
      }
    }
  } catch (err) {
    window.$notifyError(err);
  } finally {
    loading.value = false;
  }
}

watch(() => props.rivalId, (_, newValue) => {
  handleLoad();
}, { once: true });

</script>
