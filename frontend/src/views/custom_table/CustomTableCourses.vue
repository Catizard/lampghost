<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.customTableCourses') }}</n-text>
    </n-h1>
  </n-flex>
  <n-flex justify="space-between">
    <SelectCustomTable v-model:value="currentCustomTableID" style="width: 200px;" ignoreDefaultTable />
    <n-flex justify="end">
      <n-button :disabled="currentCustomTableID == null" type="primary" @click="showAddModel = true">
        {{ t('button.addCustomCourse') }}
      </n-button>
    </n-flex>
  </n-flex>
  <n-spin :show="loading">
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false"
      :row-key="(row: dto.CustomCourseDto) => row.ID"></n-data-table>
  </n-spin>
  <CustomCourseAddForm v-model:show="showAddModel" :customTableId="currentCustomTableID" @refresh="loadData" />
  <SelectSongFromFolder :customTableId="currentCustomTableID" @select="handleSelectSong"
    v-model:show="showSelectSongFromFolder" />
  <SortTableModal v-model:show="sortTableSettings.show" :query-func="sortTableSettings.queryFunc"
    @select="sortTableSettings.handleUpdateSort" :title="sortTableSettings.title"
    :labelField="sortTableSettings.labelField" :keyField="sortTableSettings.keyField" />
</template>

<script setup lang="ts">
import { h, reactive, ref, Ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import SelectCustomTable from '@/components/custom_table/SelectCustomTable.vue';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NButton, NDropdown } from 'naive-ui';
import CustomCourseDetail from './CustomCourseDetail.vue';
import { BindFolderContentToCustomCourse, FindCustomCourseList, QueryCourseSongListWithRival, QueryCustomCourseSongListWithRival, UpdateCustomCourseDataOrder } from '@wailsjs/go/main/App';
import CustomCourseAddForm from './CustomCourseAddForm.vue';
import SelectSongFromFolder from '@/components/folder/SelectSongFromFolder.vue';
import SortTableModal from '@/components/SortTableModal.vue';
import { CurrentCustomCourse, useCurrentCourseStore } from '@/stores/customCourse';

const { t } = useI18n();
const loading = ref(false);
const showAddModel = ref(false);
const currentCustomTableID: Ref<number | null> = ref(null);
const currentCustomCourseID: Ref<number | null> = ref(null);
const showSelectSongFromFolder = ref(false);
const currentCustomCourseStore = useCurrentCourseStore();

let data: Ref<dto.CustomCourseDto[]> = ref([]);
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
});
const columns: DataTableColumns<dto.CustomCourseDto> = [
  {
    type: "expand",
    renderExpand(row: dto.CustomCourseDto) {
      return h(
        CustomCourseDetail,
        { customCourseId: row.ID, customTableId: currentCustomTableID.value },
      );
    }
  },
  { title: t('column.name'), key: "Name" },
  { title: t('column.constraints'), key: "Constraints" },
  {
    title: t('column.actions'), key: "Actions", width: "75px",
    render(row: dto.CustomCourseDto) {
      return h(
        NDropdown,
        {
          trigger: "hover",
          options: [
            { label: t('button.addToCustomCourse'), key: "Bind" },
            { label: t('button.sort'), key: "Sort" },
          ],
          onSelect: (key: "Bind" | "Sort") => {
            currentCustomCourseID.value = row.ID;
            switch (key) {
              case 'Bind': showSelectSongFromFolder.value = true; break;
              case 'Sort': sortTableSettings.show = true; break;
            }
          }
        },
        { default: () => h(NButton, null, { default: () => "..." }) }
      )
    },
  }
];

function loadData() {
  loading.value = true;
  FindCustomCourseList({} as any).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    data.value = [...result.Rows];
  }).catch(err => {
    window.$notifyError(err);
  }).finally(() => {
    loading.value = false;
  });
}

function handleSelectSong(checkedRowKeys: number[]) {
  loading.value = true;
  BindFolderContentToCustomCourse(checkedRowKeys[0], currentCustomCourseID.value)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
    }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

const sortTableSettings = reactive({
  show: false,
  queryFunc: () => {
    return QueryCustomCourseSongListWithRival({
      ID: currentCustomCourseID.value,
      RivalID: 1,
    } as any);
  },
  handleUpdateSort: (ids: number[]) => {
    loading.value = true;
    UpdateCustomCourseDataOrder(ids).then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
    }).catch(err => window.$notifyError(err))
      .finally(() => loading.value = false);
  },
  title: t('title.refactorCustomCourseDataOrder'),
  labelField: "Title",
  keyField: "ID"
});

watch(currentCustomTableID, () => {
  loadData();
});
</script>
