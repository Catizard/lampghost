<template>
  <n-flex justify="start">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.playerTags') }}</n-text>
    </n-h1>
  </n-flex>
  <n-flex justify="start" style="margin-bottom: 15px;">
    <SelectRival v-model:value="currentRivalID" width="200px" defaultSelect />
    <n-button style="margin-left: auto;" type="primary" @click="showAddModal = true">
      {{ t('button.addCustomTag') }}
    </n-button>
  </n-flex>
  <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :loading="tableLoading"
    :row-key="(row: dto.RivalTagDto) => row.ID" />
  <RivalTagAddForm :rivalId="currentRivalID" v-model:show="showAddModal" @refresh="loadData" />
  <RivalTagEditForm ref="editFormRef" @refresh="loadData" />
</template>

<script setup lang="ts">
import { DeleteRivalTagByID, QueryRivalTagPageList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NButton, NDropdown, useDialog } from 'naive-ui';
import { h, reactive, Ref, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { useI18n } from 'vue-i18n';
import YesNotTag from '@/components/YesNotTag.vue';
import SelectRival from '@/components/rivals/SelectRival.vue';
import RivalTagAddForm from './RivalTagAddForm.vue';
import RivalTagEditForm from './RivalTagEditForm.vue';

const i18n = useI18n();
const { t } = i18n;
const dialog = useDialog();

const tableLoading = ref(false);
const currentRivalID: Ref<number | null> = ref(null);
const editFormRef = ref<InstanceType<typeof RivalTagEditForm>>(null);

const columns: DataTableColumns<dto.RivalTagDto> = [
  { title: t('column.tagName'), key: "TagName", width: "200px", ellipsis: { tooltip: true } },
  { title: t('column.symbol'), key: "Symbol", width: "100px", ellipsis: { tooltip: true } },
  {
    title: t('column.tagTime'),
    width: "125px",
    key: "RecordTime",
    render: (row: dto.RivalTagDto) => dayjs(row.RecordTime).format('YYYY-MM-DD HH:mm:ss')
  },
  {
    title: t('column.enabled'), key: "Enabled",
    width: "75px",
    render: (row: dto.RivalTagDto) => {
      return h(
        YesNotTag,
        { state: row.Enabled, onClick: () => { } }
      );
    }
  },
  {
    title: t('column.actions'), key: "Actions",
    width: "100px",
    render: (row: dto.RivalTagDto) => {
      return h(NDropdown, {
        trigger: "hover",
        options: [
          { label: t('button.edit'), key: "Edit", },
          { label: t('button.delete'), key: "Delete", disabled: !row.Generated }
        ],
        onSelect(key: string) {
          switch (key) {
            case "Edit": editFormRef.value.open(row.ID); break;
            case "Delete": deleteTag(row); break;
          }
        },
      }, { default: () => h(NButton, null, { default: () => '...' }) });
    }
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
  },
  reset: () => {
    pagination.page = 1;
    loadData();
  }
});
function loadData() {
  tableLoading.value = true;
  QueryRivalTagPageList({
    RivalId: currentRivalID.value,
    Pagination: pagination,
    NoIgnoreEnabled: true,
  } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows];
      pagination.pageCount = result.Pagination.pageCount;
    })
    .catch(err => window.$notifyError(err))
    .finally(() => tableLoading.value = false);
}

const showAddModal = ref(false);

function deleteTag(tag: dto.RivalTagDto) {
  dialog.warning({
    title: t('deleteDialog.title'),
    positiveText: t('deleteDialog.positiveText'),
    negativeText: t('deleteDialog.negativeText'),
    onPositiveClick: () => {
      DeleteRivalTagByID(tag.ID)
        .then(result => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          loadData();
        }).catch(err => window.$notifyError(err));
    }
  })
}

// Watch: Whenever changing current rival, reset current page to the first one
// NOTE: There is no need to call loadData(), reset method would call it
watch(currentRivalID, () => {
  pagination.reset();
})
</script>
