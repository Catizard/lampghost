<template>
  <n-spin :show="loading">
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">{{ t('title') }}</n-text>
      </n-h1>
    </n-flex>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false"
      :row-key="(row: dto.RivalInfoDto) => row.ID" />
  </n-spin>
</template>

<script lang="ts" setup>
import YesNotTag from '@/components/YesNotTag.vue';
import { QueryRivalInfoPageList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NFlex } from 'naive-ui';
import { h, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const loading = ref<boolean>(false);
const { t } = useI18n();

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
const columns = createColumns();
const data = ref<Array<dto.RivalInfoDto>>([]);

function createColumns(): DataTableColumns<dto.RivalInfoDto> {
  return [
    { title: t('column.name'), key: "Name", ellipsis: { tooltip: true }, },
    { title: t('column.version'), key: "TagName", width: "300px" },
    {
      title: t('column.reverseImport'), key: "ReverseImport", width: "100px",
      render(row: dto.RivalInfoDto) {
        return h(YesNotTag, {
          state: row.ReverseImport == 1,
          onClick: () => {
            console.log('clicked: rival id: ', row.ID);
          }
        });
      }
    },
  ]
}

function loadData() {
  loading.value = true;
  QueryRivalInfoPageList({
    Pagination: pagination,
  } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      data.value = [...result.Rows];
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

loadData();
</script>

<i18n lang="json">{
  "en": {
    "title": "Lock Player Version",
    "column": {
      "name": "Name",
      "version": "Version",
      "reverseImport": "Reverse Import",
      "actions": "Actions"
    }
  },
  "zh-CN": {
    "title": "锁定玩家版本",
    "column": {
      "name": "名称",
      "version": "版本",
      "reverseImport": "逆向导入",
      "actions": "操作"
    }
  }
}</i18n>
