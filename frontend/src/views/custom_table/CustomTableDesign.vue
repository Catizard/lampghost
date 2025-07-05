<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.designCustomTable') }}</n-text>
    </n-h1>
  </n-flex>
  <n-spin :show="loading">
    <n-flex justify="space-between">
      <SelectCustomTable v-model:value="currentCustomTableID" style="width: 200px;" />
      <n-flex justify="end">
        <n-button :disabled="currentCustomTableID == null" type="primary" @click="showAddModal = true">
          {{ t('button.addDifficultFolder') }}
        </n-button>
        <n-button :disabled="currentCustomTableID == null" type="info"
          @click="levelSortModalRef.open(currentCustomTableID)">
          {{ t('button.sortLevels') }}
        </n-button>
      </n-flex>
    </n-flex>
    <FolderTable type="table" ref="folderTableRef" :customTableId="currentCustomTableID" />
    <FolderAddForm type="table" :customTableId="currentCustomTableID" v-model:show="showAddModal" @refresh="reload" />
    <FolderSortModal ref="levelSortModalRef" @refresh="reload" />
  </n-spin>
</template>

<script lang="ts" setup>
import { ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import FolderTable from '../folder/FolderTable.vue';
import FolderAddForm from '../folder/FolderAddForm.vue';
import FolderSortModal from '../folder/FolderSortModal.vue';
import SelectCustomTable from '@/components/custom_table/SelectCustomTable.vue';

const { t } = useI18n();

const loading = ref(false);
const showAddModal = ref(false);
const folderTableRef = ref<InstanceType<typeof FolderTable>>(null);
const levelSortModalRef = ref<InstanceType<typeof FolderSortModal>>(null);

const currentCustomTableID: Ref<number | null> = ref(null);
function reload() {
  folderTableRef.value.loadData();
}
</script>
