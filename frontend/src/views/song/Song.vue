<template>
  <n-spin :show="loading">
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">{{ t('title.song') }}</n-text>
      </n-h1>
    </n-flex>
    <n-flex justify="space-between">
      <n-input :placeholder="t('form.placeholderSearchSongOrSabunName')" v-model:value="searchNameLike"
        @keyup.enter="loadData()" style="width: 350px;" />
      <n-button type="info" @click="reloadSongData">{{ t('button.reload') }}</n-button>
    </n-flex>
    <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false"
      :row-key="row => row.ID" />
  </n-spin>
  <select-folder v-model:show="showFolderSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <select-difficult v-model:show="showDifficultSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <ChartPreview ref="chartPreviewRef" />
</template>

<script setup lang="ts">
import { BindSongToFolder, QuerySongDataPageList, ReloadRivalSongData } from '@wailsjs/go/main/App';
import { dto, vo } from '@wailsjs/go/models';
import { DataTableColumns, NButton, NDropdown } from 'naive-ui';
import { h, reactive, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import ChartPreview from '@/components/ChartPreview.vue';
import SelectFolder from '../folder/SelectFolder.vue';
import SelectDifficult from '../custom_table/SelectDifficult.vue';

const { t } = useI18n();

const loading = ref(false);
const chartPreviewRef = ref<InstanceType<typeof ChartPreview>>(null);

const searchNameLike: Ref<string | null> = ref(null);
const columns: DataTableColumns<dto.RivalSongDataDto> = [
  { title: t('column.title'), key: "Title", resizable: true },
  { title: t('column.subTitle'), key: "SubTitle", resizable: true, width: "125px", ellipsis: { tooltip: true } },
  { title: t('column.artist'), key: "Artist", resizable: true, width: "125px", ellipsis: { tooltip: true } },
  { title: t('column.genre'), key: "Genre", resizable: true, width: "125px", ellipsis: { tooltip: true } },
  {
    title: t('column.actions'), key: "actions", resizable: true, width: "90px",
    render(row: dto.RivalSongDataDto) {
      return h(
        NDropdown,
        {
          trigger: "hover",
          options: [
            { label: t('button.addToFavoriteFolder'), key: "AddToFolder" },
            { label: t('button.addToCustomTable'), key: "AddToTable" },
            { label: t('button.gotoPreview'), key: "GotoPreview" },
          ],
          onSelect: (key: string) => {
            switch (key) {
              case 'AddToFolder': handleAddToFolder(row.Sha256, row.Title); break;
              case 'AddToTable': handleAddToTable(row.Sha256, row.Title); break;
              case "GotoPreview": chartPreviewRef.value.open(row.Md5); break;
            }
          }
        },
        { default: () => h(NButton, null, { default: () => '...' }) }
      )
    }
  }
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

const showFolderSelection = ref<boolean>(false);
const showDifficultSelection = ref<boolean>(false);
type SongInfo = {
  Sha256: string,
  Title: string
};
const candidateSongInfo = ref<SongInfo | null>(null);

function handleAddToFolder(sha256: string, title: string) {
  candidateSongInfo.value = {
    Sha256: sha256,
    Title: title
  };
  showFolderSelection.value = true;
}

function handleAddToTable(sha256: string, title: string) {
  candidateSongInfo.value = {
    Sha256: sha256,
    Title: title
  };
  showDifficultSelection.value = true;
}

function handleSubmit(folderIds: number[]) {
  BindSongToFolder({
    Md5: "",
    FolderIDs: folderIds,
    ...candidateSongInfo.value,
  })
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      window.$notifySuccess(t('message.bindSuccess'));
    }).catch(err => window.$notifyError(err));
}

loadData();
</script>
