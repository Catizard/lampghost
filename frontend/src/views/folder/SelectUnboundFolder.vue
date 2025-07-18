<!-- A selectable data table represents unbound folders under one custom table  -->
<!-- Used both for custom table and favrioute folder module -->
<template>
  <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false"
    :row-key="(row: dto.FolderDto) => row.ID" @update:checked-row-keys="handleCheck" :loading="loading"
    :checked-row-keys="checkedFolerIds" />
</template>

<script lang="ts" setup>
import { FindFolderList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns } from 'naive-ui';
import { reactive, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps<{
  type: "table" | "folder",
  sha256?: string,
  customTableId?: number
}>();
defineExpose({ reload });
let checkedFolerIds = defineModel<number[]>("checkedFolderIds");

const loading = ref(false);

const data: Ref<dto.FolderDto[]> = ref([]);
const columns: DataTableColumns<dto.FolderDto> = [
  { type: "selection" },
  {
    title: () => {
      if (props.type == "folder") {
        return t('column.folderName');
      } else if (props.type == "table") {
        return t('column.difficultName');
      }
    },
    key: "FolderName"
  }
];
const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
  }
});

function reload() {
  console.log('SelectUnboundFolder component reload: ', props.customTableId)
  loading.value = true;
  checkedFolerIds.value = [];
  FindFolderList({
    IgnoreSha256: props.sha256,
    CustomTableID: props.customTableId ?? 1
  } as any).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    data.value = [...result.Rows];
  })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function handleCheck(rowKeys: number[]) {
  checkedFolerIds.value = [...rowKeys];
}

watch(() => props.customTableId, () => {
  reload();
});
</script>
