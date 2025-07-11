<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.designCustomTable') }}</n-text>
    </n-h1>
  </n-flex>
  <n-spin :show="loading">
    <n-flex justify="space-between">
      <SelectCustomTable v-model:value="currentCustomTableID" style="width: 200px;" ignoreDefaultTable />
      <n-flex justify="end">
        <n-button :disabled="currentCustomTableID == null" type="primary" @click="showAddModal = true">
          {{ t('button.addDifficultFolder') }}
        </n-button>
        <n-button :disabled="currentCustomTableID == null" type="info" @click="showSortTable = true">
          {{ t('button.sortLevels') }}
        </n-button>
      </n-flex>
    </n-flex>
    <FolderTable type="table" ref="folderTableRef" :customTableId="currentCustomTableID" />
    <FolderAddForm type="table" :customTableId="currentCustomTableID" v-model:show="showAddModal" @refresh="reload" />
    <SortTableModal v-model:show="showSortTable" :query-func="queryFunc" @select="handleUpdateSort"
      :title="t('title.refactorCustomTableLevelOrder')" labelField="FolderName" keyField="ID" />
  </n-spin>
</template>

<script lang="ts" setup>
import { ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import FolderTable from '../folder/FolderTable.vue';
import FolderAddForm from '../folder/FolderAddForm.vue';
import SelectCustomTable from '@/components/custom_table/SelectCustomTable.vue';
import { FindFolderList, UpdateFolderOrder } from '@wailsjs/go/main/App';
import { dto, result } from '@wailsjs/go/models';
import SortTableModal from '@/components/SortTableModal.vue';

const { t } = useI18n();

const loading = ref(false);
const showAddModal = ref(false);
const showSortTable = ref(false);
const folderTableRef = ref<InstanceType<typeof FolderTable>>(null);

const queryFunc: (() => Promise<result.RtnDataList>) = () => FindFolderList({
  CustomTableID: currentCustomTableID.value
} as any);

const currentCustomTableID: Ref<number | null> = ref(null);
function reload() {
  folderTableRef.value.loadData();
}

function handleUpdateSort(ids: number[]) {
  loading.value = true;
  UpdateFolderOrder(ids)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      reload();
    }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}
</script>
