<!-- Display one custom course's song list -->
<template>
  <n-data-table :loading="loading" :columns="columns" :data="data" :bordered="false"
    :row-key="(row: dto.RivalSongDataDto) => row.ID" />
  <SelectSongFromFolder :customTableId="customTableId" @select="handleSelectSong"
    v-model:show="showSelectSongFromFolder" />
</template>

<script lang="ts" setup>
import ClearTag from '@/components/ClearTag.vue';
import { BindFolderContentToCustomCourse, QueryCustomCourseSongListWithRival } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NDropdown, NxButton } from 'naive-ui';
import { h, onMounted, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import SelectSongFromFolder from '@/components/folder/SelectSongFromFolder.vue';

const loading = ref(false);
const { t } = useI18n();
const showSelectSongFromFolder = ref(false);

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
  {
    title: t('column.actions'), key: "Actions", width: "75px",
    render(row: dto.RivalSongDataDto) {
      return h(
        NDropdown,
        {
          trigger: "hover",
          options: [
            { label: t('button.addToCustomCourse'), key: "Bind" },
            { label: t('button.sort'), key: "Sort" },
          ],
          onSelect: (key: "Bind" | "Sort") => {
            switch (key) {
              case 'Bind': showSelectSongFromFolder.value = true; break;
              case 'Sort': handleSortCustomCourseData(row); break;
            }
          }
        },
        { default: () => h(NxButton, null, { default: () => "..." }) }
      )
    },
  }
];


function handleSortCustomCourseData(row: dto.RivalSongDataDto) {
  window.$notifyWarning("TODO");
}

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

function handleSelectSong(checkedRowKeys: number[]) {
  loading.value = true;
  BindFolderContentToCustomCourse(checkedRowKeys[0], props.customCourseId)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
    }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

onMounted(() => {
  loadData();
});

watch(props, () => {
  loadData();
});
</script>
