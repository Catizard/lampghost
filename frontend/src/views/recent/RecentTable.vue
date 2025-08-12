<template>
  <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" :loading="loading"
    :row-key="row => row.ID" />
  <select-folder v-model:show="showFolderSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <select-difficult v-model:show="showDifficultSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <ChartPreview ref="chartPreviewRef" />
</template>

<script setup lang="ts">
import ClearTag from '@/components/ClearTag.vue';
import TableTags from '@/components/TableTags.vue';
import { useUserStore } from '@/stores/user';
import { BindSongToFolder, QueryRivalScoreDataLogPageList, QueryRivalScoreLogPageList, ReadConfig } from '@wailsjs/go/main/App';
import { config, dto } from '@wailsjs/go/models';
import dayjs from 'dayjs';
import { DataTableColumns, NDropdown, NButton, NText, NFlex } from 'naive-ui';
import { h, Ref, ref, reactive, onMounted, VNode, VNodeChild } from 'vue';
import { useI18n } from 'vue-i18n';
import SelectDifficult from '../custom_table/SelectDifficult.vue';
import SelectFolder from '../folder/SelectFolder.vue';
import ChartPreview from '@/components/ChartPreview.vue';
import SongTitleParagraph from '@/components/SongTitleParagraph.vue';

type PlayLog = dto.RivalScoreLogDto | dto.RivalScoreDataLogDto;
type SongInfo = {
  Sha256: string,
  Title: string
};

const props = defineProps<{
  songNameLike?: string
}>();
defineExpose({ loadData });

const { t } = useI18n();
const loading = ref(false);
const logType: Ref<"scorelog" | "datalog" | null> = ref(null);
const userStore = useUserStore();
const showFolderSelection = ref<boolean>(false);
const showDifficultSelection = ref<boolean>(false);
const chartPreviewRef = ref<InstanceType<typeof ChartPreview>>(null);
let candidateSongInfo = reactive<SongInfo>({
  Sha256: null,
  Title: null,
});

function createColumns(useScorelog: boolean): DataTableColumns<PlayLog> {
  const ret: DataTableColumns<PlayLog> = [
    {
      title: t('column.title'), key: "Title", resizable: true, align: 'left',
      render(row: PlayLog) {
        return h(SongTitleParagraph, { data: row });
      }
    },
    {
      title: t('column.tag'), key: "Tag", minWidth: "100px", resizable: true,
      render(row: PlayLog) {
        return h(TableTags, { tableTags: row.TableTags })
      }
    },
    {
      title: t('column.clear'), key: "Clear", minWidth: "100px", resizable: true,
      render(row: PlayLog) {
        return h(ClearTag, { clear: row.Clear }, {});
      }
    },
    {
      title: t('column.recordTime'), key: "RecordTime", minWidth: "100px", resizable: true,
      render(row: PlayLog) {
        return dayjs(row.RecordTime).format('YYYY-MM-DD HH:mm:ss');
      }
    },
  ];

  if (!useScorelog) {
    ret.push(...[
      {
        title: t('column.accuracy'), key: "Accuracy", minWidth: "100px", resizable: true,
        render(row: dto.RivalScoreDataLogDto) {
          if (row.Notes == 0) {
            return "/";
          }
          return `${(((row.Epg + row.Lpg) * 2 + row.Egr + row.Lgr) / (row.Notes * 2) * 100).toFixed(2)}%`;
        }
      },
      {
        title: t('column.perfectGreat'), key: "PerfectGreat", minWidth: "100px", resizable: true,
        render(row: dto.RivalScoreDataLogDto) {
          return row.Epg + row.Lpg;
        }
      },
      {
        title: t('column.great'), key: "Great", minWidth: "100px", resizable: true,
        render(row: dto.RivalScoreDataLogDto) {
          return row.Egr + row.Lgr;
        }
      },
      {
        title: t('column.good'), key: "Good", minWidth: "100px", resizable: true,
        render(row: dto.RivalScoreDataLogDto) {
          return row.Egd + row.Lgd;
        }
      },
      {
        title: t('column.bad'), key: "Bad", minWidth: "100px", resizable: true,
        render(row: dto.RivalScoreDataLogDto) {
          return row.Ebd + row.Lbd;
        }
      },
      {
        title: t('column.miss'), key: "Miss", minWidth: "100px", resizable: true,
        render(row: dto.RivalScoreDataLogDto) {
          return row.Ems + row.Lms;
        }
      },
    ] as DataTableColumns<PlayLog>)
  } else {
    ret.push({
      title: t('column.minbp'), key: "MinBP", minWidth: "100px", resizable: true,
      render(row: PlayLog) {
        return row.Minbp;
      }
    });
  }

  ret.push({
    title: t('column.actions'), key: "actions", resizable: true, minWidth: "90px",
    render(row: PlayLog) {
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
  });

  return ret;
}

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

async function queryConfig() {
  console.log("query config");
  try {
    const result = await ReadConfig();
    if (result.Code != 200) {
      throw result.Msg;
    }
    const config: config.ApplicationConfig = result.Data;
    logType.value = config.UseScoredatalog ? "datalog" : "scorelog";
  } catch (err) {
    throw err;
  }
}

async function loadData() {
  loading.value = true;
  try {
    if (logType.value == null) {
      queryConfig();
    }
    let arg: any = {
      RivalID: userStore.id,
      Pagination: pagination,
      SongNameLike: props.songNameLike,
      NoCourseLog: true,
    };
    const queryFunction = logType.value == "scorelog" ? QueryRivalScoreLogPageList : QueryRivalScoreDataLogPageList;
    const result = await queryFunction(arg);
    if (result.Code != 200) {
      throw Promise.reject(result.Msg);
    }
    data.value = [...result.Rows];
    pagination.pageCount = result.Pagination.pageCount;
    columns = createColumns(logType.value == "scorelog");
  } catch (err) {
    window.$notifyError(err);
  } finally {
    loading.value = false;
  }
}

let columns: DataTableColumns<PlayLog> = [];
let data: Ref<PlayLog[]> = ref([]);
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
});
</script>
