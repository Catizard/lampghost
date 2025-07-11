<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="title" :positive-text="t('button.submit')"
    :negative-text="t('button.cancel')" @positive-click="handlePositiveClick" @negative-click="handleNegativeClick"
    :mask-closable="false" style="width: 85vw;">
    <n-data-table ref="tableRef" :columns="columns" :data="data" :bordered="false" :row-key="(row: T) => row[keyField]"
      max-height="500px" />
  </n-modal>
</template>

<script setup lang="ts" generic="T">
import { result } from '@wailsjs/go/models';
import { DataTableColumns, NDataTable } from 'naive-ui';
import { nextTick, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import Sortable from "sortablejs";

const props = defineProps<{
  queryFunc: () => Promise<result.RtnDataList>,
  title: string,
  labelField: string,
  keyField: string,
}>();
const emit = defineEmits<{
  (e: 'select', value: any[]): void,
}>();
defineExpose({ open });
const { t } = useI18n();

const loading = ref(false);
const show = defineModel<boolean>("show");

const tableRef = ref<InstanceType<typeof NDataTable>>(null);
let data = reactive<T[]>([]);
const columns: DataTableColumns<T> = [
  { title: t('column.name'), key: props.labelField },
  {
    title: t('column.rowIndex'), key: "RowIndex",
    render(_, rowIndex: number) {
      return rowIndex + 1;
    }
  }
];

function hookSortable() {
  const el: HTMLDivElement = tableRef.value?.$el;
  const tbody: HTMLElement = el.querySelector('.n-data-table-tbody');
  // Naive-UI's data table could not completely render in one tick
  // Therefore we need to use this hack to hook sortable correctly
  if (tbody == null) {
    nextTick(() => {
      hookSortable();
    });
    return;
  }
  Sortable.create(tbody, {
    onEnd(evt) {
      const oldElem = data[evt.oldIndex!];
      data.splice(evt.oldIndex!, 1);
      data.splice(evt.newIndex!, 0, oldElem);
    }
  });
}

function handlePositiveClick() {
  emit('select', data.map(v => v[props.keyField]));
  show.value = false;
}

function handleNegativeClick() {
  show.value = false;
}

watch(show, (newValue) => {
  if (newValue == false) {
    return;
  }
  loading.value = true;
  props.queryFunc()
    .then((result: result.RtnDataList) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data = [...result.Rows];
      // when data.length == 0, hookSortable would make whole application stall
      // because the body of the table can never be rendered successfully
      if (data.length > 0) {
        // Here goes the magic
        hookSortable();
      }
    }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
})
</script>
