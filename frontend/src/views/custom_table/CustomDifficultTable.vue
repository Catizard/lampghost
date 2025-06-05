<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title') }}</n-text>
    </n-h1>
    <n-flex justify="flex-end">
      <!-- <n-button :loading="loading" type="primary" @click="showAddModal = true">
        {{ t('button.add') }}
      </n-button> -->
    </n-flex>
  </n-flex>

  <n-data-table remote :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
    :row-key="(row: dto.CustomDiffTableDto) => row.ID" />

  <CustomDifficultTableEditForm ref="editFormRef" @refresh="loadData()" />
</template>

<script lang="ts" setup>
import { DataTableColumns, NButton } from 'naive-ui';
import { h, reactive, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { dto } from '@wailsjs/go/models';
import { QueryCustomDiffTablePageList } from '@wailsjs/go/main/App';
import CustomDifficultTableEditForm from './CustomDifficultTableEditForm.vue';

const { t } = useI18n();
const editFormRef = ref<InstanceType<typeof CustomDifficultTableEditForm>>(null);

const loading = ref(false);

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
const columns: DataTableColumns<dto.CustomDiffTableDto> = [
  { title: t('column.name'), key: "Name" },
  {
    title: t('column.actions'), key: "actions", width: "75px",
    render(row: dto.CustomDiffTableDto) {
      return h(NButton, { type: "primary", size: "small", onClick: () => editFormRef.value.open(row.ID) }, { default: () => t('button.update') })
    }
  }
];
const data: Ref<dto.CustomDiffTableDto[]> = ref([]);

function loadData() {
  loading.value = true;
  QueryCustomDiffTablePageList({
    Pagination: pagination,
  } as any).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    data.value = [...result.Rows];
    pagination.pageCount = result.Pagination.pageCount;
  })
  .catch(err => window.$notifyError(err))
  .finally(() => loading.value = false)
}

loadData();
</script>

<i18n lang="json">{
  "en": {
    "title": "Custom Difficult Table",
    "column": {
      "name": "Name",
      "actions": "Actions"
    },
    "button": {
      "update": "Edit"
    }
  },
  "zh-CN": {
    "title": "自定义难度表",
    "column": {
      "name": "名称",
      "actions": "操作"
    },
    "button": {
      "update": "修改"
    }
  }
}</i18n>