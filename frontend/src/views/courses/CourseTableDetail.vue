<!-- Display one course's song list -->
<template>
  <n-data-table :loading="loading" :columns="columns" :data="data" :boarded="false"
    :row-key="(row: dto.RivalSongDataDto) => row.ID" />
</template>

<script setup lang="ts">
import ClearTag from '@/components/ClearTag.vue';
import { useUserStore } from '@/stores/user';
import { QueryCourseSongListWithRival } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NButton, useDialog } from 'naive-ui';
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
  { title: t('column.name'), key: "Title" },
  {
    title: t('column.clear'), key: "Clear",
    width: "125px",
    render(row: dto.RivalSongDataDto) {
      return h(ClearTag, { clear: row.Lamp })
    }
  },
  { title: t('column.minbp'), key: "MinBP", width: "75px", },
  {
    title: t('column.ghost'), key: "GhostLamp", width: "125px",
    render(row: dto.RivalSongDataDto) {
      return h(ClearTag, { clear: row.GhostLamp, },)
    }
  },
  { title: t('column.minbp'), key: "GhostMinBP", width: "75px", },
];

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
