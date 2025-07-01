<template>
  <perfect-scrollbar>
    <n-h1 prefix="bar" style="text-align: left;">
      <n-text type="primary">{{ t('title.courses') }}</n-text>
    </n-h1>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" />
  </perfect-scrollbar>
</template>

<script setup lang="ts">
import { DataTableColumns } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { ref, h } from 'vue';
import { FindCourseInfoListWithRival } from '@wailsjs/go/main/App';
import ClearTag from "@/components/ClearTag.vue"
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const { t } = i18n;

const columns = createColumns();
const pagination = false as const;
let data = ref<Array<dto.CourseInfoDto>>([]);

function createColumns(): DataTableColumns<dto.CourseInfoDto> {
  return [
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

function loadData() {
  // TODO: remove magical 1
  FindCourseInfoListWithRival(1)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      data.value = [...result.Rows]
    }).catch(err => window.$notifyError(err));
}

loadData();
</script>
