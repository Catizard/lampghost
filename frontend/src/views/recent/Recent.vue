<template>
  <perfect-scrollbar>
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">Recent Records</n-text>
      </n-h1>
      <n-flex justify="end">
        <n-button>choose clear type</n-button>
        <n-button>minimum clear type</n-button>
      </n-flex>

      <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" :loading="loading" />
    </n-flex>
  </perfect-scrollbar>
</template>

<script lang="ts" setup>
import { QueryRivalScoreLogPageList } from '@wailsjs/go/controller/RivalScoreLogController';
import { vo } from '@wailsjs/go/models';
import { DataTableColumn, DataTableColumns, NButton, useNotification } from 'naive-ui';
import { getDefaultPageSize } from 'naive-ui/es/pagination/src/utils';
import { h, onMounted, reactive, Ref, ref } from 'vue';

const notification = useNotification();
const loading: Ref<boolean> = ref(false);

function createColumns(): DataTableColumns<any> {
  return [
    { title: "Song Name", key: "Title", resizable: true },
    { title: "Tag", key: "Tag", minWidth: "100px", resizable: true },
    { title: "Clear", key: "Lamp", minWidth: "100px", resizable: true },
    { title: "Time", key: "RecordTime", minWidth: "100px", resizable: true },
    {
      title: "Action",
      key: "actions",
      resizable: true,
      minWidth: "150px",
      render(row) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: "small",
          },
          { default: () => "Add to Folder" }
        )
      }
    }
  ]
}

function loadData() {
  loading.value = true;
  let arg: any = {
    Page: pagination.page,
    PageSize: pagination.pageSize,
  };
  QueryRivalScoreLogPageList(arg)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows];
    })
    .catch(err => {
      notification.error({
        content: "读取最近游玩记录出错:" + err,
        duration: 3000,
        keepAliveOnHover: true
      });
    }).finally(() => {
      loading.value = false;
    });
}

const columns = createColumns();
let data: Ref<Array<any>> = ref([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    loadData();
  }
});

onMounted(() => {
  loadData();
})
</script>
