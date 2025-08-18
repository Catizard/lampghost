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
      :row-key="(row: dto.CourseInfoDto) => row.ID" :rowClassName="rowClassName" />
  </n-spin>
</template>

<script setup lang="ts">
import { DataTableColumns } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { ref, h, reactive, Ref, watch, inject } from 'vue';
import { FindCourseInfoListWithRival } from '@wailsjs/go/main/App';
import { useI18n } from 'vue-i18n';
import CourseTableDetail from './CourseTableDetail.vue';
import SelectDifficultTable from '@/components/difficult_table/SelectDifficultTable.vue';
import SelectRivalTag from '@/components/rivals/SelectRivalTag.vue';
import SelectRival from '@/components/rivals/SelectRival.vue';
import { useUserStore } from '@/stores/user';
import SongClearParagraph from '@/components/SongClearParagraph.vue';
import { queryClearTypeColorStyle } from '@/constants/cleartype';

const i18n = useI18n();
const { t } = i18n;
const resetScrollBar = inject("resetScrollBar", () => { });

const loading = ref(false);
const userStore = useUserStore();
const columns = createColumns();
const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
    resetScrollBar();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
  }
});
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
      title: t('column.clear'), key: "Clear", align: "center", className: "clearColumn",
      render(row: dto.CourseInfoDto) {
        return h(SongClearParagraph, {
          clearType: row.Clear,
          scoreOption: null,
          bestRecordTimestamp: row.FirstClearTimestamp != 0 ? row.FirstClearTimestamp : null
        });
      }
    },
  ]
}

function rowClassName(row: dto.CourseInfoDto): string {
  let clearText = queryClearTypeColorStyle(row.Clear).text;
  let ghostClearText = queryClearTypeColorStyle(row.GhostClear).text;
  return `${clearText} ghost-${ghostClearText}`
}

const currentDiffTableID: Ref<number | null> = ref(null);
const currentRivalID: Ref<number | null> = ref(null);
const currentRivalTagID: Ref<number | null> = ref(null);

function loadData() {
  loading.value = true;
  FindCourseInfoListWithRival({
    HeaderID: currentDiffTableID.value,
    RivalID: userStore.id,
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

<style lang="css" scoped>
@import "@/assets/css/clearBackground.css";
</style>
