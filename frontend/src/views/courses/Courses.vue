<template>
  <n-h1 prefix="bar" style="text-align: left;">
    <n-text type="primary">{{ t('title.courses') }}</n-text>
  </n-h1>
  <n-flex justify="start">
    <SelectDifficultTable v-model:value="currentDiffTableID" style="width: 200px" />
    <SelectRival v-model:value="currentRivalID" width="200px" :placeholder="t('form.placeholderRival')" clearable />
    <SelectRivalTag v-model:value="currentRivalTagID" :rivalId="currentRivalID" width="200px"
      :placeholder="t('form.placeholderRivalTag')" clearable />
  </n-flex>
  <n-spin :show="loading">
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false"
      :row-key="(row: dto.CourseInfoDto) => row.ID" />
  </n-spin>
</template>

<script setup lang="ts">
import { DataTableColumns } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { ref, h, reactive, Ref, watch } from 'vue';
import { FindCourseInfoListWithRival } from '@wailsjs/go/main/App';
import ClearTag from "@/components/ClearTag.vue"
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n';
import CourseTableDetail from './CourseTableDetail.vue';
import SelectDifficultTable from '@/components/difficult_table/SelectDifficultTable.vue';
import SelectRivalTag from '@/components/rivals/SelectRivalTag.vue';
import SelectRival from '@/components/rivals/SelectRival.vue';

const i18n = useI18n();
const { t } = i18n;

const loading = ref(false);
const columns = createColumns();
const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
  }
});;
let data = ref<Array<dto.CourseInfoDto>>([]);

function createColumns(): DataTableColumns<dto.CourseInfoDto> {
  return [
    {
      type: "expand",
      renderExpand(row: dto.CourseInfoDto) {
        return h(
          CourseTableDetail,
          {
            courseId: row.ID,
            ghostRivalId: currentRivalID.value,
            ghostRivalTagId: currentRivalTagID.value
          }
        );
      }
    },
    { title: t('column.name'), key: "Name", },
    { title: t('column.constraints'), key: "Constraints" },
    {
      title: t('column.clear'),
      key: "Clear",
      render(row) {
        return h(
          ClearTag,
          {
            clear: row.Clear
          },
        )
      }
    },
    {
      title: t('column.firstClearTime'),
      key: "FirstClearTime",
      render(row: dto.CourseInfoDto) {
        if (row.Clear > 1) {
          return dayjs(row.FirstClearTimestamp).format('YYYY-MM-DD HH:mm:ss');
        } else {
          return "/";
        }
      }
    },
  ]
}

const currentDiffTableID: Ref<number | null> = ref(null);
const currentRivalID: Ref<number | null> = ref(null);
const currentRivalTagID: Ref<number | null> = ref(null);

function loadData() {
  loading.value = true;
  // TODO: remove magical 1
  FindCourseInfoListWithRival({
    HeaderID: currentDiffTableID.value,
    RivalID: 1
  } as any).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg)
    }
    data.value = [...result.Rows]
  })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false)
}

watch(currentDiffTableID, () => {
  loadData();
});
</script>
