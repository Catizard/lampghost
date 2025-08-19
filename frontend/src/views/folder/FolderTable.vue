<template>
  <n-data-table :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
    :row-key="(row: dto.FolderDto) => row.ID" :max-height="maxHeight" />
</template>

<script setup lang="ts">
import { h, Ref, ref, watch } from "vue";
import { dto } from "@wailsjs/go/models";
import {
  DataTableColumns,
  NButton,
  NDataTable,
  useDialog,
} from "naive-ui";
import { useI18n } from "vue-i18n";
import FolderDetail from "./FolderDetail.vue";
import { DelFolder, FindFolderList } from "@wailsjs/go/main/App";

const checkedRowKeys = defineModel<number[]>("checkedRowKeys");
const props = defineProps<{
  customTableId: number
  type: "table" | "folder",
  selectSong?: "single" | "multiple", // need to be passed to FolderDetail
  noActions?: boolean // also being passed to FolderDetail
  maxHeight?: string
}>();
defineExpose({ loadData })

const i18n = useI18n();
const { t } = i18n;
const dialog = useDialog();
const emit = defineEmits<{
  (e: 'update:checkedRowKeys', value: number[])
}>()

const loading = ref(false);
const pagination = false as const;
let data: Ref<Array<any>> = ref([]);
const columns: DataTableColumns<dto.FolderDto> = createColumns({
  deleteFolder(row: any) {
    dialog.warning({
      title: t('deleteDialog.title'),
      positiveText: t('deleteDialog.positiveText'),
      negativeText: t('deleteDialog.negativeText'),
      onPositiveClick: () => {
        deleteFolder(row.ID);
      },
    });
  },
});

function createColumns({
  deleteFolder,
}: {
  deleteFolder: (row: dto.FolderDto) => void;
}): DataTableColumns<dto.FolderDto> {
  let ret: DataTableColumns<dto.FolderDto> = [
    {
      type: "expand",
      renderExpand: (row: dto.FolderDto) => {
        return h(FolderDetail, {
          folderId: row.ID,
          type: props.type,
          selectSong: props.selectSong,
          checkedRowKeys: checkedRowKeys.value,
          noActions: props.noActions,
          'onUpdate:checkedRowKeys': value => emit('update:checkedRowKeys', value)
        });
      },
    },
    {
      title: (): string => {
        if (props.type == "table") {
          return (t('column.difficultName'));
        } else if (props.type == "folder") {
          return t('column.folderName');
        }
      },
      key: "FolderName",
      minWidth: "150px"
    },

  ];
  if (props.noActions == null || props.noActions == false) {
    ret.push({
      title: t('column.actions'),
      key: "actions",
      render(row: dto.FolderDto) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: "small",
            type: "error",
            onClick: () => deleteFolder(row),
          },
          { default: () => t('button.delete') },
        );
      }
    })
  }
  return ret;
}

function deleteFolder(id: number) {
  DelFolder(id)
    .then((result) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      window.$notifySuccess(t('message.deleteSuccess'));
      loadData();
    })
    .catch((err) => {
      window.$notifyError(err);
      loadData();
    });
}

function loadData() {
  // Do nothing if custom table is empty
  if (props.customTableId == null || props.customTableId == 0) {
    return;
  }
  loading.value = true;
  FindFolderList({
    RivalID: 1,
    CustomTableID: props.customTableId
  } as any)
    .then((result) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows];
    })
    .catch((err) => window.$notifyError(t('message.loadFolderDataError', { msg: err })))
    .finally(() => loading.value = false);
}

watch(() => props.customTableId, () => {
  loadData();
}, { immediate: true });
</script>
