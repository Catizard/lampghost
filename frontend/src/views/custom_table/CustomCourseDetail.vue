<!-- Display one custom course's song list -->
<template>
  <n-data-table :loading="loading" :columns="columns" :data="data" :bordered="false"
    :row-key="(row: dto.RivalSongDataDto) => row.ID" />
</template>

<script lang="ts" setup>
import ClearTag from '@/components/ClearTag.vue';
import { useUserStore } from '@/stores/user';
import { DeleteCustomCourseData, QueryCustomCourseSongListWithRival } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NButton, useDialog } from 'naive-ui';
import { h, onMounted, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const loading = ref(false);
const userStore = useUserStore();
const { t } = useI18n();
const dialog = useDialog();

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
    title: t('column.actions'), key: "Actions", width: "100px",
    render(row: dto.RivalSongDataDto) {
      return h(NButton, {
        type: "error",
        size: "small",
        tertiary: true,
        onClick: () => {
          dialog.warning({
            title: t('deleteDialog.title'),
            positiveText: t('deleteDialog.positiveText'),
            negativeText: t('deleteDialog.negativeText'),
            onPositiveClick: () => {
              loading.value = true;
              DeleteCustomCourseData(row.ID)
                .then(result => {
                  if (result.Code != 200) {
                    return Promise.reject(result.Msg);
                  }
                  loadData();
                })
                .catch(err => window.$notifyError(err))
                .finally(() => loading.value = false);
            }
          });
        }
      }, { default: () => t('button.delete') })
    }
  }
];

function loadData() {
  loading.value = true;
  QueryCustomCourseSongListWithRival({
    ID: props.customCourseId,
    RivalID: userStore.id,
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

<style lang="css" scoped></style>
