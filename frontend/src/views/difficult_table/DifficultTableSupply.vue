<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.supplyMissingBMSFromTable')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-data-table :columns="columns" :data="data" :bordered="false" :row-key="(row: Level) => row.name"
      v-model:checked-row-keys="selectedLevels" max-height="75vh" />
  </n-modal>
</template>

<script setup lang="ts">
import { FindDiffTableLevelList, SupplyMissingBMSFromTable } from '@wailsjs/go/main/App';
import { DataTableColumns } from 'naive-ui';
import { Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';

type Level = {
  name: string
};

const { t } = useI18n();
const show = ref(false);
const loading = ref(false);
const tableId = ref(null);
const selectedLevels: Ref<string[]> = ref([]);
defineExpose({ open });

function open(difficultTableId: number) {
  show.value = true;
  loading.value = true;
  tableId.value = difficultTableId;
  FindDiffTableLevelList(difficultTableId).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    data.value = [...result.Rows.map(level => {
      return {
        name: level
      }
    })];
    selectedLevels.value = [...result.Rows];
  }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false)
}

async function handlePositiveClick() {
  if (selectedLevels.value.length == 0) {
    window.$notifyError(t('message.noSelectedLevel'));
    return;
  }
  loading.value = true;
  try {
    const result = await SupplyMissingBMSFromTable(tableId.value, selectedLevels.value);
    if (result.Code != 200) {
      throw result.Msg;
    }
  } catch (err) {
    window.$notifyError(err);
  } finally {
    loading.value = false;
  }
}

function handleNegativeClick() {
  show.value = false;
}

let data: Ref<Level[]> = ref([]);
const columns: DataTableColumns<Level> = [
  { type: "selection" },
  { title: t('column.name'), key: "name" },
  {
    title: t('column.rowIndex'), key: "RowIndex",
    render(_, rowIndex: number) {
      return rowIndex + 1;
    }
  }
]
</script>
