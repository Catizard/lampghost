<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start">
      <n-text type="primary">{{ t('title.tableManagement') }}</n-text>
    </n-h1>
    <n-flex justify="end">
      <n-button :loading="loading" type="info" @click="showSortModal = true">
        {{ t('button.sort') }}
      </n-button>
      <n-button :loading="loading" type="primary" @click="showAddModal = true">
        {{ t('button.add') }}
      </n-button>
    </n-flex>
  </n-flex>
  <n-data-table :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
    :row-key="(row: dto.DiffTableHeaderDto) => row.ID" />

  <DifficultTableSortModal v-model:show="showSortModal" @refresh="loadDiffTableData()" />
  <DifficultTableAddForm v-model:show="showAddModal" @refresh="loadDiffTableData()" />
  <DifficultTableEditForm ref="editFormRef" @refresh="loadDiffTableData()" />
  <SortTableModal v-model:show="sortTableSettings.show" :query-func="sortTableSettings.queryFunc"
    @select="sortTableSettings.handleUpdateSort" :title="sortTableSettings.title"
    :labelField="sortTableSettings.labelField" :keyField="sortTableSettings.keyField" />
</template>

<script lang="ts" setup>
import type { DropdownOption } from "naive-ui";
import { NButton, NDataTable, NDropdown, NTag, useDialog } from "naive-ui";
import { h, reactive, Ref, ref } from "vue";
import {
  DelDiffTableHeader,
  FindDiffTableHeaderTree,
  ReloadDiffTableHeader,
  SupplyMissingBMSFromTable,
  UpdateHeaderOrder
} from "@wailsjs/go/main/App";
import { dto } from "@wailsjs/go/models";
import { useI18n } from "vue-i18n";
import DifficultTableAddForm from "./DifficultTableAddForm.vue";
import DifficultTableEditForm from "./DifficultTableEditForm.vue";
import { TagColor } from "naive-ui/es/tag/src/common-props";
import SortTableModal from "@/components/SortTableModal.vue";

const i18n = useI18n();
const { t } = i18n;
const showSortModal = ref(false);
const showAddModal = ref(false);
const editFormRef = ref<InstanceType<typeof DifficultTableEditForm>>(null);

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
    }).catch(err => window.$notifyError(err))
      .finally(() => loading.value = false);
  },
  title: t('title.refactorDifficultTableLevelOrder'),
  labelField: "Name",
  keyField: "ID"
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
    dialog.warning({
      title: t('message.supplyMissingBMS'),
      positiveText: t('message.yes'),
      negativeText: t('message.no'),
      onPositiveClick: () => {
        loading.value = true;
        SupplyMissingBMSFromTable(row.ID)
          .then(result => {
            if (result.Code != 200) {
              return Promise.reject(result.Msg);
            }
          })
          .catch(err => window.$notifyError(err))
          .finally(() => loading.value = false);
      }
    })
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
    sortTableSettings.show = true;
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
