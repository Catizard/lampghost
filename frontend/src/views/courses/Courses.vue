<template>
  <perfect-scrollbar>
    <n-h1 prefix="bar" style="text-align: left;">
      <n-text type="primary">courses</n-text>
    </n-h1>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" />
  </perfect-scrollbar>
</template>

<script setup lang="ts">
import { DataTableColumns, useNotification } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { ref } from 'vue';
import { FindCourseInfoList } from '@wailsjs/go/controller/CourseInfoController';

const notification = useNotification();

const columns = createColumns();
const pagination = false as const;
let data = ref<Array<dto.CourseInfoDto>>([]);

function createColumns(): DataTableColumns<dto.CourseInfoDto> {
  return [
    { title: "Name", key: "Name", },
    { title: "Constraints", key: "Constraints" },
    { title: "Clear", key: "Clear" },
    { title: "First Clear", key: "FirstClearTime" },
  ]
}

function loadData() {
  FindCourseInfoList()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      data.value = [...result.Rows]
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      })
    })
}

loadData();
</script>