<template>
  <calendar-heatmap :values="data" :endDate="dayjs()" :tooltip-formatter="tooltipFormatter" :no-data-text="'No Play'"
    style="width: 95%;" />
</template>

<script setup lang="ts">
import { QueryUserKeyCountInYear } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { useNotification } from 'naive-ui';
import { ref, Ref } from 'vue';
import dayjs from 'dayjs';
import 'vue3-calendar-heatmap/dist/style.css';

const props = defineProps<{
  rivalId?: number
}>();

const notification = useNotification();

const data: Ref<Array<any>> = ref([]);

function loadKeyCountData() {
  // TODO: Should we use the "year" from PlayCountChart component?
  QueryUserKeyCountInYear({
    SpecifyYear: dayjs().get('year').toString(),
    RivalId: props.rivalId,
  } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      console.log(result);
      data.value = [...result.Rows.map((row: dto.KeyCountDto) => {
        return {
          count: row.KeyCount,
          date: row.RecordDate,
        }
      })];
      console.log(data.value);
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true
      });
    })
}

function tooltipFormatter(value: { date: Date | string; count: number; }): string {
  const { date, count } = value;
  const mmdd = dayjs(date).format("MM-DD")
  return `${count} ${mmdd}`;
}

loadKeyCountData();

</script>
