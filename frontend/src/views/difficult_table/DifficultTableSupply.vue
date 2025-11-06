<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="title" :positive-text="t('button.submit')"
    :negative-text="t('button.cancel')" @positive-click="handlePositiveClick" @negative-click="handleNegativeClick"
    :mask-closable="false">
    <n-data-table :columns="columns" :data="data" :bordered="false"
      :row-key="(row: dto.DiffTableHeaderDto) => row.Level" v-model:checked-row-keys="selectedLevels"
      max-height="75vh" />
  </n-modal>
</template>

<script setup lang="ts">
import { FindDownloadableLevelList, SupplyMissingBMSFromTable } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns } from 'naive-ui';
import { computed, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const show = ref(false);
const loading = ref(false);
const tableId = ref(null);
const tableSymbol = ref("");
const tableName = ref("");
const selectedLevels: Ref<string[]> = ref([]);
defineExpose({ open });

// I'm lazy to query symbol in this function
function open(difficultTableId: number, symbol: string, name: string) {
  show.value = true;
  loading.value = true;
  tableId.value = difficultTableId;
  tableSymbol.value = symbol;
  tableName.value = name;
  FindDownloadableLevelList(difficultTableId).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    data.value = [...result.Rows];
    selectedLevels.value = [...result.Rows.map(header => {
      return header.Level;
    })];
  }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false)
}

const title = computed(() => {
  return `${t('title.supplyMissingBMSFromTable')} - ${tableName.value}`;
});

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

let data: Ref<dto.DiffTableHeaderDto[]> = ref([]);
const columns: DataTableColumns<dto.DiffTableHeaderDto> = [
  { type: "selection" },
  {
    title: t('column.level'), key: "Level", render(row: dto.DiffTableHeaderDto) {
      return `${tableSymbol.value}${row.Level}`;
    }
  },
  { title: t('column.lostCount'), key: "LostCount" },
  { title: t('column.songCount'), key: "SongCount" },
]
</script>
