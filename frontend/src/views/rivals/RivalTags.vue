<template>
  <perfect-scrollbar>
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">{{ t('title') }}</n-text>
      </n-h1>
    </n-flex>
    <n-flex justify="flex-start">
      <n-select :loading="tableLoading" v-model:value="currentRivalID" :options="rivalOptions" style="width: 200px;" />
    </n-flex>
    <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :loading="tableLoading"
      :row-key="(row: dto.RivalTagDto) => row.ID" />
  </perfect-scrollbar>
</template>

<script setup lang="ts">
import { FindRivalInfoList } from '@wailsjs/go/controller/RivalInfoController';
import { FindRivalTagList, QueryRivalTagPageList } from '@wailsjs/go/controller/RivalTagController';
import { dto, vo } from '@wailsjs/go/models';
import { DataTableColumns, SelectOption, useNotification } from 'naive-ui';
import { reactive, Ref, ref, watch } from 'vue';
import * as dayjs from 'dayjs';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const { t } = i18n;
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
        return Promise.reject(t('message.noRivalError'));
      }
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
  { title: t('column.tagName'), key: "TagName", },
  { title: t('column.generated'), key: "Generated", },
  {
    title: t('column.tagTime'),
    key: "RecordTime",
    render: (row: dto.RivalTagDto) => dayjs(row.RecordTime).format('YYYY-MM-DD HH:mm:ss')
  }
];
const data: Ref<Array<dto.RivalTagDto>> = ref([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
    loadData();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    loadData();
  }
});
function loadData() {
  tableLoading.value = true;
  QueryRivalTagPageList({
    RivalId: currentRivalID.value,
    Pagination: pagination,
  } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows];
      pagination.pageCount = result.Pagination.pageCount;
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
watch(currentRivalID, () => {
  loadData();
})

</script>

<i18n lang="json">{
  "en": {
    "title": "Rival Tags",
    "column": {
      "tagName": "Tag Name",
      "generated": "Generated",
      "tagTime": "Tag Time"
    },
    "message": {
      "noRivalError": "FATAL ERROR: no rival data found"
    }
  },
  "zh-CN": {
    "title": "玩家标签",
    "column": {
      "tagName": "标签名称",
      "generated": "自动生成",
      "tagTime": "标签时间"
    },
    "message": {
      "noRivalError": "未知错误: 找不到任何玩家信息?"
    }
  }
}</i18n>
