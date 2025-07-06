<template>
  <n-data-table remote :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
    :row-key="(row: dto.FolderDto) => row.ID" v-model:checked-row-keys="checkedRowKeys" />
</template>

<script lang="ts" setup>
import { dto } from '@wailsjs/go/models';
import { create, DataTableColumns, NButton } from 'naive-ui';
import { h, reactive, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import ClearTag from '@/components/ClearTag.vue';
import { DelFolderContent, QueryFolderContentWithRival } from '@wailsjs/go/main/App';
import TableTags from '@/components/TableTags.vue';

const loading = ref(false);
const { t } = useI18n();

const { folderId, rivalId = 1, type, selectSong, noActions = false } = defineProps<{
  folderId: number,
  rivalId?: number,
  type: "table" | "folder",
  selectSong?: "single" | "multiple",
  noActions?: boolean
}>();
const checkedRowKeys = defineModel<number[]>("checkedRowKeys");

let data: Ref<dto.FolderContentDto[]> = ref([]);
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
function createColumns(): DataTableColumns<dto.FolderContentDto> {
  let columns: DataTableColumns<dto.FolderContentDto> = [];
  if (selectSong) {
    columns.push({
      type: "selection",
      multiple: selectSong == "multiple"
    });
  }
  columns.push(
    { title: t('column.name'), key: "Title" },
    {
      title: (): string => {
        if (type == "folder") {
          return t('column.tag');
        } else if (type == "table") {
          return t('column.externalTag');
        }
      },
      key: "Tag",
      width: "200px",
      render(row: dto.FolderContentDto) {
        return h(TableTags, { tableTags: row.TableTags });
      }
    },
    {
      title: t('column.clear'), key: "Clear",
      width: "125px",
      render: (row: dto.FolderContentDto) => {
        return h(ClearTag, { clear: row.Lamp },)
      }
    },
  );
  if (!noActions) {
    columns.push({
      title: t('column.actions'),
      key: "actions",
      width: "150px",
      render(row: dto.FolderContentDto) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: "small",
            type: "error",
            onClick: () => deleteFolderContent(row.ID),
          },
          { default: () => t('button.delete') },
        );
      },
    });
  }
  return columns;
}
const columns: DataTableColumns<dto.FolderContentDto> = createColumns();
function loadData() {
  loading.value = true;
  QueryFolderContentWithRival({
    RivalID: rivalId,
    FolderID: folderId,
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

function deleteFolderContent(id: number) {
  DelFolderContent(id)
    .then((result) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      window.$notifySuccess(t('message.deleteSuccess'));
    }).catch(err => window.$notifyError(err))
    .finally(() => loadData());
}

watch([() => folderId, () => rivalId], () => {
  loadData();
});
loadData();
</script>
