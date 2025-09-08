<!-- Display one course's song list -->
<template>
  <n-data-table :loading="loading" :columns="columns" :data="data" :boarded="false"
    :row-key="(row: dto.RivalSongDataDto) => row.ID" :rowClassName="rowClassName" />
</template>

<script setup lang="ts">
import SongClearParagraph from '@/components/SongClearParagraph.vue';
import SongTitleParagraph from '@/components/SongTitleParagraph.vue';
import { queryClearTypeColorStyle } from '@/constants/cleartype';
import { useUserStore } from '@/stores/user';
import { QueryCourseSongListWithRival } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns } from 'naive-ui';
import { h, onMounted, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const loading = ref(false);
const userStore = useUserStore();
const { t } = useI18n();

const props = defineProps<{
  courseId: number,
  ghostRivalId?: number,
  ghostRivalTagId?: number
}>();

let data: Ref<dto.RivalSongDataDto[]> = ref([]);
const columns: DataTableColumns<dto.RivalSongDataDto> = [
  {
    title: t('column.name'), key: "Title",
    render(row: dto.RivalSongDataDto) {
      return h(SongTitleParagraph, { data: row });
    }
  },
  {
    title: t('column.clear'), key: "Clear", width: "140px", className: "clearColumn", align: "center",
    render(row: dto.RivalSongDataDto) {
      return h(SongClearParagraph, {
        clearType: row.Lamp,
        scoreOption: null,
        bestRecordTimestamp: row.BestRecordTimestamp != 0 ? row.BestRecordTimestamp : null
      });
    }
  },
  { title: t('column.minbp'), key: "MinBP", width: "75px", align: "center" },
  {
    title: t('column.ghost'), key: "GhostLamp", width: "140px", className: "ghostClearColumn", align: "center",
    render(row: dto.RivalSongDataDto) {
      return h(SongClearParagraph, {
        clearType: row.GhostLamp,
        scoreOption: null,
        bestRecordTimestamp: null,
      });
    }
  },
  { title: t('column.minbp'), key: "GhostMinBP", width: "75px", align: "center" },
];

function rowClassName(row: dto.RivalSongDataDto): string {
  let clearText = queryClearTypeColorStyle(row.Lamp).text;
  let ghostClearText = queryClearTypeColorStyle(row.GhostLamp).text;
  return `${clearText} ghost-${ghostClearText}`;
}

function loadData() {
  loading.value = true;
  QueryCourseSongListWithRival({
    ID: props.courseId,
    RivalID: userStore.id,
    GhostRivalID: props.ghostRivalId ?? 0,
    GhostRivalTagID: props.ghostRivalTagId ?? 0
  } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows];
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false)
}

onMounted(() => {
  loadData();
});

watch(props, () => {
  loadData();
});
</script>

<style lang="css" scoped>
@import "@/assets/css/clearBackground.css";
</style>
