<!-- Display one custom course's song list -->
<template>
  <n-data-table :loading="loading" :columns="columns" :data="data" :bordered="false"
    :row-key="(row: dto.RivalSongDataDto) => row.ID" />
</template>

<script lang="ts" setup>
import ClearTag from '@/components/ClearTag.vue';
import { QueryCustomCourseSongListWithRival } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns } from 'naive-ui';
import { h, onMounted, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const loading = ref(false);
const { t } = useI18n();

const props = defineProps<{
  customCourseId: number,
  customTableId: number,
}>();

let data: Ref<dto.RivalSongDataDto[]> = ref([]);
const columns: DataTableColumns<dto.RivalSongDataDto> = [
  { title: t('column.name'), key: 'Title' },
  {
    title: t('column.clear'), key: "Clear",
    width: "125px",
    render(row: dto.RivalSongDataDto) {
      return h(ClearTag, { clear: row.Lamp })
    }
  },
  { title: t('column.minbp'), key: "MinBP", width: "75px", },
];

function loadData() {
  loading.value = true;
  // TODO: Remove magical 1
  QueryCustomCourseSongListWithRival({
    ID: props.customCourseId,
    RivalID: 1
  } as any).then(result => {
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
