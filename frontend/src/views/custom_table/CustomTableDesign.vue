<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title') }}</n-text>
    </n-h1>
  </n-flex>
  <n-spin :show="loading">
    <n-flex justify="space-between">
      <n-select v-model:value="currentCustomTableID" :options="customTableOptions" style="width: 200px;" />
      <n-flex justify="end">
        <n-button v-if="customTableOptions.length > 0" type="primary" @click="showAddModal = true">
          {{ t('button.add') }}
        </n-button>
      </n-flex>
    </n-flex>
    <template v-if="customTableOptions.length > 0">
      <FolderTable type="table" ref="folderTableRef" :customTableId="currentCustomTableID" />
    </template>
    <FolderAddForm type="table" :customTableId="currentCustomTableID" v-model:show="showAddModal" @refresh="reload" />
  </n-spin>
</template>

<script lang="ts" setup>
import { FindCustomDiffTableList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { SelectOption } from 'naive-ui';
import { ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import FolderTable from '../folder/FolderTable.vue';
import FolderAddForm from '../folder/FolderAddForm.vue';

const { t } = useI18n();

const loading = ref(false);
const showAddModal = ref(false);
const folderTableRef = ref<InstanceType<typeof FolderTable>>(null);

const currentCustomTableID: Ref<number | null> = ref(null);
const customTableOptions: Ref<SelectOption[]> = ref([]);
function loadCustomTableOptions() {
  loading.value = true;
  FindCustomDiffTableList({
    IgnoreDefaultTable: true
  } as any).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    if (result.Rows.length == 0) {
      return Promise.reject(t('message.noTableError'))
    }
    customTableOptions.value = result.Rows.map((row: dto.CustomDiffTableDto): SelectOption => {
      return {
        label: row.Name,
        value: row.ID
      }
    });
    currentCustomTableID.value = customTableOptions.value[0].value as number;
  })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

loadCustomTableOptions();

function reload() {
  folderTableRef.value.loadData();
}
</script>

<i18n lang="json">{
  "en": {
    "title": "Custom Table Design",
    "button": {
      "add": "Add Difficult"
    },
    "message": {
      "noTableError": "No custom difficult table found"
    }
  },
  "zh-CN": {
    "title": "自定义难度表设计",
    "button": {
      "add": "新增难度"
    },
    "message": {
      "noTableError": "当前无任何自定义难度表"
    }
  }
}</i18n>
