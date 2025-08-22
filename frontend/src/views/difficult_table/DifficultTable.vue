<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start">
      <n-text type="primary">{{ t('title.tableManagement') }}</n-text>
    </n-h1>
    <n-flex justify="end">
      <n-button :loading="loading" type="info" @click="sortTableSettings.show = true">
        {{ t('button.sort') }}
      </n-button>
      <n-button :loading="loading" type="primary" @click="showAddModal = true">
        {{ t('button.addDifficultTable') }}
      </n-button>
    </n-flex>
  </n-flex>
  <n-data-table :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
    :row-key="(row: dto.DiffTableHeaderDto) => row.ID" />

  <DifficultTableAddForm v-model:show="showAddModal" @refresh="loadDiffTableData()" />
  <DifficultTableEditForm ref="editFormRef" @refresh="loadDiffTableData()" />
  <SortTableModal v-model:show="sortTableSettings.show" :query-func="sortTableSettings.queryFunc"
    @select="sortTableSettings.handleUpdateSort" :title="sortTableSettings.title"
    :labelField="sortTableSettings.labelField" :keyField="sortTableSettings.keyField" />
  <SortTableModal v-model:show="sortLevelSettings.show" :query-func="sortLevelSettings.queryFunc"
    @select="sortLevelSettings.handleUpdateSort" :title="sortLevelSettings.title"
    :labelField="sortLevelSettings.labelField" :keyField="sortLevelSettings.keyField" />
  <DifficultTableSupply ref="supplyFormRef" />
</template>

<script lang="ts" setup>
import type { DropdownOption } from "naive-ui";
import { NButton, NDataTable, NDropdown, NTag, useDialog } from "naive-ui";
import { h, reactive, Ref, ref } from "vue";
import {
  DelDiffTableHeader,
  FindDiffTableHeaderTree,
  FindDiffTableLevelList,
  ReloadDiffTableHeader,
  UpdateHeaderLevelOrders,
  UpdateHeaderOrder
} from "@wailsjs/go/main/App";
import { dto } from "@wailsjs/go/models";
import { useI18n } from "vue-i18n";
import DifficultTableAddForm from "./DifficultTableAddForm.vue";
import DifficultTableEditForm from "./DifficultTableEditForm.vue";
import { TagColor } from "naive-ui/es/tag/src/common-props";
import SortTableModal from "@/components/SortTableModal.vue";
import DifficultTableSupply from "./DifficultTableSupply.vue";

const i18n = useI18n();
const { t } = i18n;
const showAddModal = ref(false);
const editFormRef = ref<InstanceType<typeof DifficultTableEditForm>>(null);
const supplyFormRef = ref<InstanceType<typeof DifficultTableSupply>>(null);

const dialog = useDialog();
const loading = ref(false);
loadDiffTableData();

const sortTableSettings = reactive({
  show: false,
  queryFunc: () => {
    return FindDiffTableHeaderTree(null);
  },
  handleUpdateSort: (ids: number[]) => {
    loading.value = true;
    UpdateHeaderOrder(ids).then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      loadDiffTableData();
    }).catch(err => window.$notifyError(err))
      .finally(() => loading.value = false);
  },
  title: t('title.refactorDifficultTableOrder'),
  labelField: "Name",
  keyField: "ID"
});

const sortLevelSettings = reactive({
  tableId: 0,
  show: false,
  queryFunc: () => {
    return FindDiffTableLevelList(sortLevelSettings.tableId).then(result => {
      result.Rows = result.Rows.map(row => {
        return {
          key: row,
          label: row,
        }
      });
      return result;
    });
  },
  handleUpdateSort: (levels: string[]) => {
    loading.value = true;
    UpdateHeaderLevelOrders({
      ID: sortLevelSettings.tableId,
      LevelOrders: levels.join(",")
    } as any).then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      loadDiffTableData();
    }).catch(err => window.$notifyError(err))
      .finally(() => loading.value = false);
  },
  title: t('title.refactorDifficultTableLevelOrder'),
  keyField: "key",
  labelField: "label"
});

const columns = [
  { title: t('column.name'), key: "Name", },
  {
    title: t('column.tag'), key: "Tag",
    render(row: dto.DiffTableHeaderDto) {
      let tagColorProp: TagColor = {};
      if (row.TagColor != '') {
        tagColorProp.color = row.TagColor;
      }
      if (row.TagTextColor != '') {
        tagColorProp.textColor = row.TagTextColor;
      }
      return h(
        NTag,
        { color: tagColorProp },
        { default: () => row.Symbol == '' ? '/' : row.Symbol },
      )
    }
  },
  { title: t('column.url'), key: "HeaderUrl", },
  {
    title: t('column.actions'),
    key: "actions",
    render(row) {
      return h(
        NDropdown,
        {
          trigger: "hover",
          options: otherActionOptions,
          onSelect: (key: string) => handleSelectOtherAction(row, key)
        },
        {
          default: () => h(
            NButton,
            null,
            { default: () => '...' }
          )
        },
      );
    },
  },
];

let data: Ref<Array<any>> = ref([]);
const pagination = false as const;

const otherActionOptions: Array<DropdownOption> = [
  { label: t('button.reload'), key: "Reload" },
  { label: t('button.edit'), key: "Edit", },
  { label: t('button.supplyMissingBMS'), key: "Supply" },
  { label: t('button.sortLevels'), key: "SortLevels", },
  {
    label: t('button.delete'),
    key: "Delete",
    props: {
      style: "color: red"
    }
  }
];

function handleSelectOtherAction(row: dto.DiffTableHeaderDto, key: string) {
  if ("Reload" === key) {
    reloadTableHeader(row.ID);
  }
  if ("Supply" === key) {
    supplyFormRef.value.open(row.ID, row.Symbol);
  }
  if ("Delete" === key) {
    dialog.warning({
      title: t('deleteDialog.title'),
      positiveText: t('deleteDialog.positiveText'),
      negativeText: t('deleteDialog.negativeText'),
      onPositiveClick: () => {
        delDiffTableHeader(row.ID)
      }
    })
  }
  if ("Edit" === key) {
    editFormRef.value.open(row.ID);
  }
  if ("SortLevels" === key) {
    sortLevelSettings.tableId = row.ID;
    sortLevelSettings.show = true;
  }
}

function reloadTableHeader(id: number) {
  loading.value = true;
  ReloadDiffTableHeader(id)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      loadDiffTableData();
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function delDiffTableHeader(id: number) {
  loading.value = true;
  DelDiffTableHeader(id)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      window.$notifySuccess(t('message.deleteSuccess'));
      loadDiffTableData();
    }).catch(err => {
      window.$notifyError(err);
      loadDiffTableData();
    }).finally(() => loading.value = false);
}

function loadDiffTableData() {
  loading.value = true;
  FindDiffTableHeaderTree(null)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows]
    })
    .catch((err) => window.$notifyError(t('message.loadTableDataError', { msg: err })))
    .finally(() => loading.value = false);
}

</script>
