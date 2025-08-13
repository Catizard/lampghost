<template>
  <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" min-height="500px"
    :loading="loading" :row-key="(row: dto.DiffTableDataDto) => row.ID" @update-sorter="handleUpdateSorter" />
  <select-folder v-model:show="showFolderSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <select-difficult v-model:show="showTableSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <ChartPreview ref="chartPreviewRef" />
</template>

<script setup lang="ts">
import { DataTableColumns, DataTableSortState, NButton, NDropdown, NFlex, NText, NTooltip } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { h, reactive, ref, Ref, watch } from 'vue';
import { BindSongToFolder, QueryDiffTableDataWithRival, SubmitSingleMD5DownloadTask } from '@wailsjs/go/main/App';
import SelectFolder from '@/views/folder/SelectFolder.vue';
import ClearTag from '@/components/ClearTag.vue';
import ChartPreview from '@/components/ChartPreview.vue';
import { useI18n } from 'vue-i18n';
import { BrowserOpenURL } from '@wailsjs/runtime/runtime';
import SelectDifficult from '../custom_table/SelectDifficult.vue';
import dayjs from 'dayjs';
import { useUserStore } from '@/stores/user';
import SongTitleParagraph from '@/components/SongTitleParagraph.vue';
import SongScoreParagraph from '@/components/SongScoreParagraph.vue';

const i18n = useI18n();
const { t } = i18n;
const loading = ref<boolean>(false);
const userStore = useUserStore();

const props = defineProps<{
  headerId?: number
  level?: string
  ghostRivalId?: number
  ghostRivalTagId?: number
}>();

const chartPreviewRef = ref<InstanceType<typeof ChartPreview>>(null);

const sorter: Ref<Sorter> = ref({
  SortBy: null,
  SortOrder: null,
});
const columns: DataTableColumns<dto.DiffTableDataDto> = [
  {
    title: t('column.title'), key: "Title", resizable: true, sorter: true,
    render: (row: dto.DiffTableDataDto) => {
      return h(SongTitleParagraph, { data: row, lost: row.DataLost });
    }
  },
  {
    title: t('column.score'), key: "Score", width: "110px", resizable: true, align: "center",
    render(row: dto.DiffTableDataDto) {
      return h(SongScoreParagraph, { data: row });
    }
  },
  {
    title: t('column.clear'), key: "Lamp", width: "100px", resizable: true, sorter: true,
    render(row: dto.DiffTableDataDto) {
      return h(ClearTag, { clear: row.Lamp },)
    }
  },
  {
    // TODO: Change sorter from false to true when user choosed ghost rival
    title: t('column.ghost'), key: "GhostLamp", width: "100px", resizable: true, sorter: true,
    render(row: dto.DiffTableDataDto) {
      return h(ClearTag, { clear: row.GhostLamp, },)
    }
  },
  {
    title: t('column.lastPlayed'), key: "LastPlayedTimestamp", width: "100px", resizable: true, sorter: true,
    render(row: dto.DiffTableDataDto) {
      if (row.LastPlayedTimestamp == 0) {
        return "NO PLAY";
      }
      // return dayjs(row.LastPlayedTimestamp * 1000).format("YYYY-MM-DD HH:mm:ss");
      return h(
        NTooltip,
        { trigger: "hover" },
        {
          trigger: () => dayjs(row.LastPlayedTimestamp * 1000).format("YYYY-MM-DD"),
          default: () => dayjs(row.LastPlayedTimestamp * 1000).format("YYYY-MM-DD HH:mm:ss")
        }
      );
    }
  },
  { title: t('column.playCount'), key: "PlayCount", width: "75px" },
  {
    title: t('column.actions'),
    key: "actions",
    resizable: true,
    width: "90px",
    render(row: dto.DiffTableDataDto) {
      return h(
        NDropdown,
        {
          trigger: "hover",
          options: [
            { label: t('button.addToFavoriteFolder'), key: "AddToFolder" },
            { label: t('button.addToCustomTable'), key: "AddToTable" },
            { label: t('button.download'), key: "Download" },
            { label: t('button.gotoSongURL'), key: "GotoURL", disabled: row.Url == "" },
            { label: t('button.gotoURLDiff'), key: "GotoURLDiff", disabled: row.UrlDiff == "" },
            { label: t('button.gotoLR2IR'), key: "GotoLR2IR", disabled: row.Md5 == "" },
            { label: t('button.gotoMochaIR'), key: "GotoMochaIR", disabled: row.Sha256 == "" },
            { label: t('button.gotoPreview'), key: "GotoPreview", disabled: row.Md5 == "" },
          ],
          onSelect: (key: string) => {
            const md5 = row.Md5;
            const sha256 = row.Sha256;
            switch (key) {
              case 'AddToFolder': handleAddToFolder(row); break;
              case 'AddToTable': handleAddToTable(row); break;
              case 'Download': handleSubmitSingleMD5DownloadTask(row); break;
              case "GotoURL": BrowserOpenURL(row.Url); break;
              case "GotoURLDiff": BrowserOpenURL(row.UrlDiff); break;
              case "GotoLR2IR": BrowserOpenURL(`https://www.dream-pro.info/~lavalse/LR2IR/search.cgi?mode=ranking&bmsmd5=${md5}`); break;
              case "GotoMochaIR": BrowserOpenURL(`https://mocha-repository.info/song.php?sha256=${sha256}`); break;
              // case "GotoPreview": BrowserOpenURL(`https://bms-score-viewer.pages.dev/view?md5=${md5}`); break;
              case "GotoPreview": chartPreviewRef.value.open(row.Md5); break;
            }
          }
        },
        { default: () => h(NButton, null, { default: () => '...' }) }
      );
    }
  }
];

let data: Ref<Array<dto.DiffTableDataDto>> = ref([]);
function loadData() {
  loading.value = true;
  QueryDiffTableDataWithRival({
    ID: props.headerId,
    Level: props.level,
    RivalID: userStore.id,
    GhostRivalID: props.ghostRivalId ?? 0,
    GhostRivalTagID: props.ghostRivalTagId ?? 0,
    Pagination: pagination,
    SortBy: sorter.value.SortBy,
    SortOrder: sorter.value.SortOrder,
  } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows];
      pagination.pageCount = result.Pagination.pageCount;
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

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

watch(props, () => {
  loadData()
});

type SongInfo = {
  Sha256: string,
  Title: string
};
const candidateSongInfo = ref<SongInfo | null>(null);
const showFolderSelection = ref<boolean>(false);
function handleAddToFolder(row: dto.DiffTableDataDto) {
  candidateSongInfo.value = {
    Sha256: row.Sha256,
    Title: row.Title
  };
  showFolderSelection.value = true;
}

const showTableSelection = ref<boolean>(false);
function handleAddToTable(row: dto.DiffTableDataDto) {
  candidateSongInfo.value = {
    Sha256: row.Sha256,
    Title: row.Title,
  };
  showTableSelection.value = true;
}

function handleSubmit(selected: [number]) {
  BindSongToFolder({
    Md5: "",
    FolderIDs: selected,
    ...candidateSongInfo.value
  } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      window.$notifySuccess(t('message.bindSuccess'));
    }).catch(err => window.$notifyError(err));
}

function handleUpdateSorter(option: DataTableSortState | null) {
  // TODO: This is a pain in the a**
  switch (option.columnKey) {
    case "Title": sorter.value.SortBy = "title"; break;
    case "Artist": sorter.value.SortBy = "artist"; break;
    default: sorter.value.SortBy = option.columnKey as string; break;
  }
  if (option.order != false) {
    sorter.value.SortOrder = option.order;
  }
  loadData();
}

function handleSubmitSingleMD5DownloadTask(row: dto.DiffTableDataDto) {
  loading.value = true;
  SubmitSingleMD5DownloadTask(row.Md5, row.Title)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      window.$notifySuccess(t('message.submitSuccess'));
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

loadData();
</script>
