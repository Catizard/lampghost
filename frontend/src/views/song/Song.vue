<template>
  <n-spin :show="loading">
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">{{ t('title') }}</n-text>
      </n-h1>
    </n-flex>
    <n-flex justify="space-between">
      <n-input :placeholder="t('searchNamePlaceholder')" v-model:value="searchNameLike" @keyup.enter="loadData()"
        style="width: 350px;" />
      <n-button type="primary" @click="reloadSongData">{{ t('button.reload') }}</n-button>
    </n-flex>
    <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false"
      :row-key="row => row.ID" />
  </n-spin>
</template>

<script setup lang="ts">
import { QuerySongDataPageList, ReloadRivalSongData } from '@wailsjs/go/main/App';
import { dto, vo } from '@wailsjs/go/models';
import { DataTableColumns } from 'naive-ui';
import { reactive, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const loading = ref(false);

const searchNameLike: Ref<string | null> = ref(null);
const columns: DataTableColumns<dto.RivalSongDataDto> = [
  { title: t('column.name'), key: "Title", resizable: true },
  { title: t('column.artist'), key: "Artist", resizable: true },
];
const data: Ref<dto.RivalSongDataDto[]> = ref([]);
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

async function loadData() {
  loading.value = true;
  try {
    // TODO: Remove magical 1 here
    let arg: vo.RivalSongDataVo = {
      RivalID: 1,
      TitleLike: searchNameLike.value,
      Pagination: pagination
    } as any;
    const result = await QuerySongDataPageList(arg);
    if (result.Code != 200) {
      throw result.Msg;
    }
    data.value = [...result.Rows];
    pagination.pageCount = result.Pagination.pageCount;
  } catch (e) {
    window.$notifyError(e);
  } finally {
    loading.value = false;
  }
}

async function reloadSongData() {
  loading.value = true;
  try {
    const result = await ReloadRivalSongData();
    if (result.Code != 200) {
      throw result.Msg;
    }
  } catch (e) {
    window.$notifyError(e);
  } finally {
    loading.value = false;
  }
  loadData();
}

loadData();
</script>

<i18n lang="json">{
  "en": {
    "title": "Song",
    "button": {
      "reload": "Reload"
    },
    "column": {
      "name": "Title",
      "artist": "Artist"
    },
    "searchNamePlaceholder": "Search Song/Sabun Name"
  },
  "zh-CN": {
    "title": "歌曲列表",
    "button": {
      "reload": "同步文件"
    },
    "column": {
      "name": "谱面名称",
      "artist": "谱师"
    },
    "searchNamePlaceholder": "搜索歌曲/差分名称"
  }
}</i18n>
