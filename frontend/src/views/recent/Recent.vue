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

      <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" :loading="loading"
        :row-key="row => row.ID" />
    </n-flex>
  </perfect-scrollbar>

  <select-folder v-model:show="showFolderSelection" @submit="handleSubmit"></select-folder>
</template>

<script lang="ts" setup>
import ClearTag from '@/components/ClearTag.vue';
import { QueryRivalScoreLogPageList } from '@wailsjs/go/controller/RivalScoreLogController';
import { DataTableColumns, NButton, useNotification } from 'naive-ui';
import { h, onMounted, reactive, Ref, ref } from 'vue';
import SelectFolder from '../folder/SelectFolder.vue';
import { BindRivalSongDataToFolder } from '@wailsjs/go/controller/RivalSongDataController';
import { dto } from '@wailsjs/go/models';

const notification = useNotification();
const loading = ref<boolean>(false);
const showFolderSelection = ref<boolean>(false);
const candidateSongDataID = ref<number>(null);

function handleSubmit(selected: [any]) {
  console.log('ID=', candidateSongDataID.value)
  BindRivalSongDataToFolder(candidateSongDataID.value, selected)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      notification.success({
        content: "添加成功",
        duration: 5000,
        keepAliveOnHover: true
      });
    }).catch((err) => {
      notification.error({
        content: "添加至收藏夹失败:" + err,
        duration: 5000,
        keepAliveOnHover: true
      });
    });
}

function handleAddToFolder(ID: number) {
  candidateSongDataID.value = ID;
  showFolderSelection.value = true;
}

function createColumns(): DataTableColumns<dto.RivalScoreLogDto> {
  return [
    { title: "Song Name", key: "Title", resizable: true },
    { title: "Tag", key: "Tag", minWidth: "100px", resizable: true },
    {
      title: "Clear", key: "Clear", minWidth: "100px", resizable: true,
      render(row) {
        return h(
          ClearTag,
          {
            clear: row.Clear
          },
          {}
        )
      }
    },
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
            onClick: () => handleAddToFolder(row.RivalSongDataID),
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
      if (result.Rows.length > 0) {
        pagination.pageCount = result.Rows[0].PageCount;
      } else {
        pagination.pageCount = 0;
      }
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

onMounted(() => {
  loadData();
})
</script>
