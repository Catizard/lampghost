<template>
  <perfect-scrollbar>
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">玩家tag管理</n-text>
      </n-h1>
    </n-flex>
    <n-flex justify="flex-start">
      <n-select :loading="tableLoading" v-model:value="currentRivalID" :options="rivalOptions" style="width: 200px;" />
    </n-flex>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :loading="tableLoading"
      :row-key="(row: dto.RivalTagDto) => row.ID" />
  </perfect-scrollbar>
</template>

<script setup lang="ts">
import { FindRivalInfoList } from '@wailsjs/go/controller/RivalInfoController';
import { FindRivalTagList } from '@wailsjs/go/controller/RivalTagController';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, SelectOption, useNotification } from 'naive-ui';
import { Ref, ref, watch } from 'vue';
import * as dayjs from 'dayjs';

const notification = useNotification();

const tableLoading = ref(false);
const currentRivalID: Ref<number | null> = ref(null);
const rivalOptions: Ref<Array<SelectOption>> = ref([]);
function loadRivalOptions() {
  tableLoading.value = true;
  FindRivalInfoList()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Rows.length == 0) {
        return Promise.reject("一个rival都没有, 发生什么事了?");
      }
      console.log(result.Rows)
      rivalOptions.value = result.Rows.map((row: dto.RivalInfoDto) => {
        return {
          label: row.Name,
          value: row.ID,
        } as SelectOption
      });
      currentRivalID.value = rivalOptions.value[0].value as number;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      })
    }).finally(() => {
      tableLoading.value = false;
    });
}
loadRivalOptions();

const columns: DataTableColumns<dto.RivalTagDto> = [
  { title: "Tag Name", key: "TagName", },
  { title: "Generated", key: "Generated", },
  {
    title: "Tag Time",
    key: "TagTime",
    render: (row: dto.RivalTagDto) => dayjs(row.TagTime).format('YYYY-MM-DD HH:mm:ss')
  }
];
const data: Ref<Array<dto.RivalTagDto>> = ref([]);
const pagination = false as const;
function loadData(rivalID: number) {
  tableLoading.value = true;
  FindRivalTagList({ RivalId: rivalID } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows];
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      });
    }).finally(() => {
      tableLoading.value = false;
    });
}

// Watch: Whenever changing current rival, reload the tag table
watch(currentRivalID, (newID: number) => {
  loadData(newID);
})

</script>
