<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.recentPlay') }}</n-text>
    </n-h1>
  </n-flex>
  <n-flex justify="start">
    <n-input :placeholder="t('form.placeholderSearchSongOrSabunName')" v-model:value="searchNameLike"
      @keyup.enter="loadData()" style="width: 350px;" />
    <!-- <n-button>{{ t('button.chooseClearType') }}</n-button>
      <n-button>{{ t('button.minimumClearType') }}</n-button> -->
  </n-flex>
  <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" :loading="loading"
    :row-key="row => row.ID" />
  <select-folder v-model:show="showFolderSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <select-difficult v-model:show="showDifficultSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <ChartPreview ref="chartPreviewRef" />
</template>

<script lang="ts" setup>
import ClearTag from '@/components/ClearTag.vue';
import { BindSongToFolder, QueryRivalScoreLogPageList } from '@wailsjs/go/main/App';
import { DataTableColumns, NButton, NDropdown, useModal } from 'naive-ui';
import { h, onMounted, reactive, Ref, ref } from 'vue';
import SelectFolder from '../folder/SelectFolder.vue';
import { dto, vo } from '@wailsjs/go/models';
import { useI18n } from 'vue-i18n';
import dayjs from 'dayjs';
import TableTags from '@/components/TableTags.vue';
import SelectDifficult from '../custom_table/SelectDifficult.vue';
import ChartPreview from '@/components/ChartPreview.vue';
import { useUserStore } from '@/stores/user';

const i18n = useI18n();
const { t } = i18n;

const loading = ref<boolean>(false);
const userStore = useUserStore();
const showFolderSelection = ref<boolean>(false);
const showDifficultSelection = ref<boolean>(false);
const searchNameLike: Ref<string | null> = ref(null);
const chartPreviewRef = ref<InstanceType<typeof ChartPreview>>(null);

type SongInfo = {
  Sha256: string,
  Title: string
};
let candidateSongInfo = reactive<SongInfo>({
  Sha256: null,
  Title: null,
});

function handleAddToFolder(sha256: string, title: string) {
  candidateSongInfo = {
    Sha256: sha256,
    Title: title
  };
  showFolderSelection.value = true;
}

function handleAddToTable(sha256: string, title: string) {
  candidateSongInfo = {
    Sha256: sha256,
    Title: title
  };
  showDifficultSelection.value = true;
}

function handleSubmit(folderIds: number[]) {
  BindSongToFolder({
    Md5: "",
    FolderIDs: folderIds,
    ...candidateSongInfo,
  })
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      window.$notifySuccess(t('message.bindSuccess'));
    }).catch(err => window.$notifyError(err));
}

function createColumns(): DataTableColumns<dto.RivalScoreLogDto> {
  return [
    { title: t('column.title'), key: "Title", resizable: true },
    {
      title: t('column.tag'), key: "Tag", minWidth: "100px", resizable: true,
      render(row: dto.RivalScoreLogDto) {
        return h(TableTags, { tableTags: row.TableTags })
      }
    },
    {
      title: t('column.clear'), key: "Clear", minWidth: "100px", resizable: true,
      render(row: dto.RivalScoreLogDto) {
        return h(ClearTag, { clear: row.Clear }, {});
      }
    },
    {
      title: t('column.recordTime'), key: "RecordTime", minWidth: "100px", resizable: true,
      render(row: dto.RivalScoreLogDto) {
        return dayjs(row.RecordTime).format('YYYY-MM-DD HH:mm:ss');
      }
    },
    {
      title: t('column.minbp'), key: "MinBP", minWidth: "100px", resizable: true,
      render(row: dto.RivalScoreLogDto) {
        return row.Minbp;
      }
    },
    {
      title: t('column.actions'), key: "actions", resizable: true, minWidth: "90px",
      render(row: dto.RivalScoreLogDto) {
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
        );
      }
    }
  ]
}

function loadData() {
  loading.value = true;
  let arg: vo.RivalScoreLogVo = {
    RivalID: userStore.id,
    Pagination: pagination,
    SongNameLike: searchNameLike.value,
    NoCourseLog: true,
  } as any;
  QueryRivalScoreLogPageList(arg)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      console.log(result);
      data.value = [...result.Rows];
      pagination.pageCount = result.Pagination.pageCount;
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
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
