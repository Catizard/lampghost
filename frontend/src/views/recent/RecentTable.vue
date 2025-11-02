<template>
  <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" :loading="loading"
    :row-key="row => row.ID" :rowClassName="rowClassName" />
  <select-folder v-model:show="showFolderSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <select-difficult v-model:show="showDifficultSelection" :sha256="candidateSongInfo?.Sha256" @submit="handleSubmit" />
  <ChartPreview ref="chartPreviewRef" />
</template>

<script setup lang="ts">
import TableTags from '@/components/TableTags.vue';
import { useUserStore } from '@/stores/user';
import { BindSongToFolder, QueryRivalScoreDataLogPageList, QueryRivalScoreLogPageList } from '@wailsjs/go/main/App';
import { config, dto } from '@wailsjs/go/models';
import dayjs from 'dayjs';
import { DataTableColumns, NDropdown, NButton, NTooltip, NText, NFlex } from 'naive-ui';
import { h, Ref, ref, reactive, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import SelectDifficult from '../custom_table/SelectDifficult.vue';
import SelectFolder from '../folder/SelectFolder.vue';
import ChartPreview from '@/components/ChartPreview.vue';
import SongTitleParagraph from '@/components/SongTitleParagraph.vue';
import SongScoreParagraph from '@/components/SongScoreParagraph.vue';
import { ClearType, ClearTypeDef, DefaultClearTypeColorStyle, queryClearTypeColorStyle } from '@/constants/cleartype';
import SongClearParagraph from '@/components/SongClearParagraph.vue';
import RecordTimeParagraph from '@/components/RecordTimeParagraph.vue';
import MinBPParagraph from '@/components/MinBPParagraph.vue';
import { useConfigStore } from '@/stores/config';
import { ScoreOption } from '@/constants/scoreOption';

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
      title: t('column.tag'), key: "Tag", minWidth: "100px", resizable: true, align: "center",
      render(row: PlayLog) {
        return h(TableTags, { tableTags: row.TableTags });
      }
    },
    {
      title: t('column.clear'), key: "Clear", width: "140px", resizable: true, className: "clearColumn", align: "center",
      render(row: PlayLog) {
        const p = { clearType: row.Clear, scoreOption: null, disableTimestamp: true };
        if (!useScorelog) {
          p.scoreOption = (row as dto.RivalScoreDataLogDto).Option;
        }
        return h(SongClearParagraph, p);
      }
    },
  ];

  if (!useScorelog) {
    ret.push(...[
      {
        title: t('column.score'), key: "Score", width: "80px", resizable: true,
        render(row: dto.RivalScoreDataLogDto) {
          return h(SongScoreParagraph, { data: row });
        }
      },
      {
        title: t('column.minbp'), key: "MinBP", width: "100px", resizable: true, align: "center",
        render(row: dto.RivalScoreDataLogDto) {
          return h(MinBPParagraph, {
            noplay: false,
            optionValue: row.Option,
            minbp: row.Minbp
          });
        }
      },
      // NOTE: It's impossible to calculate bp here since beatoraja doesn't provide
      // complete data: we don't know there're how many 'passnotes'
    ] as DataTableColumns<PlayLog>)
  } else {
    ret.push(...[
      {
        title: t('column.minbp'), key: "MinBP", width: "100px", resizable: true, align: "center",
        render(row: dto.RivalScoreDataLogDto) {
          return h(MinBPParagraph, { noplay: false, minbp: row.Minbp, noOption: true });
        }
      }
    ] as DataTableColumns<PlayLog>)
  }

  ret.push(...[
    {
      title: t('column.recordTime'), key: "RecordTime", width: "120px", resizable: true, align: "center",
      render(row: PlayLog) {
        return h(RecordTimeParagraph, {
          recordTimestamp: row.RecordTimestamp * 1000,
        })
      }
    },
    {
      title: t('column.actions'), key: "actions", resizable: true, width: "90px",
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
                case "GotoPreview": {
                  let option = null;
                  if (!useScorelog) {
                    const dataLog = row as dto.RivalScoreDataLogDto
                    if (dataLog.Option == ScoreOption.MIRROR) {
                      option = 1;
                    } else if (dataLog.Option == ScoreOption.RANDOM) {
                      option = dataLog.RandomPattern;
                      console.log("md5: %s, option: %s", row.Md5, option);
                    }
                  }
                  chartPreviewRef.value.open(row.Md5, option, configStore.config.PreviewSide as ("P1" | "P2"));
                  break;
                }
              }
            }
          },
          { default: () => h(NButton, null, { default: () => '...' }) }
        );
      }
    }] as DataTableColumns<PlayLog>);

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

const configStore = useConfigStore();
async function loadData() {
  loading.value = true;
  try {
    console.log('no logType, querying');
    if (logType.value == null) {
      logType.value = configStore.config.UseScoredatalog != 0 ? "datalog" : "scorelog";
    }
    console.log('logType: ', logType.value);
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

function rowClassName(row: PlayLog): string {
  for (const [k, v] of Object.entries(ClearType).reverse()) {
    if (row.Clear == parseInt(k)) {
      const def: ClearTypeDef = DefaultClearTypeColorStyle[k];
      return def.text;
    }
  }
  return "";
}

onMounted(() => {
  loadData();
});
</script>

<style lang="css" scoped>
@import "@/assets/css/clearBackground.css"
</style>
