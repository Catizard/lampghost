<template>
  <n-spin :show="loading">
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">{{ t('title.playerIRImportManagement') }}</n-text>
      </n-h1>
    </n-flex>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false"
      :row-key="(row: dto.RivalInfoDto) => row.ID" />
  </n-spin>
  <RivalEditReverseImportForm ref="editFormRef" @refresh="loadData()" />
</template>

<script lang="ts" setup>
import YesNotTag from '@/components/YesNotTag.vue';
import { QueryRivalInfoPageList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NButton, NFlex } from 'naive-ui';
import { h, reactive, ref, VNode } from 'vue';
import { useI18n } from 'vue-i18n';
import RivalEditReverseImportForm from './RivalEditReverseImportForm.vue';

const loading = ref<boolean>(false);
const { t } = useI18n();
const editFormRef = ref<InstanceType<typeof RivalEditReverseImportForm>>(null);

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
    {
      title: t('column.actions'), key: "action", width: "100px",
      render(row: dto.RivalInfoDto): VNode {
        return h(NButton, {
          type: "primary",
          size: "small",
          onClick: () => editFormRef.value.open(row.ID),
        }, { default: () => t('button.edit') });
      }
    }
  ]
}

function loadData() {
  loading.value = true;
  QueryRivalInfoPageList({
    Pagination: pagination,
    IgnoreMainUser: true
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
