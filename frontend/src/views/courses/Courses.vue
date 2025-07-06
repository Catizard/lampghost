<template>
  <n-h1 prefix="bar" style="text-align: left;">
    <n-text type="primary">{{ t('title.courses') }}</n-text>
  </n-h1>
  <n-flex justify="start">
    <SelectDifficultTable v-model:value="currentDiffTableID" style="width: 200px" />
    <n-select :loading="loading" v-model:value="currentRivalID" :options="rivalOptions" style="width: 200px;"
      :placeholder="t('form.placeholderRival')" />
    <SelectRivalTag v-model:value="currentRivalTagID" :rivalId="currentRivalID" width="200px"
      :placeholder="t('form.placeholderRivalTag')" />
  </n-flex>
  <n-spin :show="loading">
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false"
      :row-key="(row: dto.CourseInfoDto) => row.ID" />
  </n-spin>
</template>

<script setup lang="ts">
import { DataTableColumns, SelectOption } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { ref, h, reactive, Ref, watch } from 'vue';
import { FindCourseInfoListWithRival, FindRivalInfoList } from '@wailsjs/go/main/App';
import ClearTag from "@/components/ClearTag.vue"
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n';
import CourseTableDetail from './CourseTableDetail.vue';
import SelectDifficultTable from '@/components/difficult_table/SelectDifficultTable.vue';
import SelectRivalTag from '@/components/SelectRivalTag.vue';

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
const rivalOptions: Ref<Array<SelectOption>> = ref([]);

function loadRivalOptions() {
  loading.value = true;
  FindRivalInfoList()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Rows.length == 0) {
        return Promise.reject(t('message.noRivalError'));
      }
      rivalOptions.value = result.Rows.map((row: dto.RivalInfoDto) => {
        return {
          label: row.Name,
          value: row.ID,
        } as SelectOption
      });
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}
loadRivalOptions();


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
