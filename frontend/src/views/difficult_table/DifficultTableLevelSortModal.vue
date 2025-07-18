<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.refactorDifficultTableLevelOrder')" style="width: 75%;"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-data-table ref="tableRef" :loading="loading" :columns="columns" :data="data" :bordered="false"
      :row-key="(row: dto.DiffTableHeaderDto) => row.ID" max-height="500px" />
  </n-modal>
</template>

<script lang="ts" setup>
import { nextTick, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { dto } from '@wailsjs/go/models';
import { FindDiffTableHeaderTree, UpdateHeaderLevelOrders } from '@wailsjs/go/main/App';
import { NDataTable } from 'naive-ui';
import Sortable from "sortablejs";

const show = defineModel<boolean>("show");
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose({ open });

const { t } = useI18n();
const loading = ref(false);

const columns = [
  { title: t('column.name'), key: "Name" },
  {
    title: t('column.rowIndex'), key: "RowIndex",
    render: (_, rowIndex: number) => {
      return rowIndex + 1
    }
  }
];

const tableRef = ref<InstanceType<typeof NDataTable>>();
let data = reactive([]);
const currentHeaderId = ref(null);
function open(headerId: number) {
  if (headerId == null || headerId == 0) {
    window.$notifyError(t('message.noChosenHeaderError'));
    show.value = false;
    return;
  }

  show.value = true;
  loading.value = true;
  currentHeaderId.value = headerId;
  FindDiffTableHeaderTree({ ID: headerId } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      const { Children } = result.Rows[0] as dto.DiffTableHeaderDto;
      data = [...Children];
      // when data.length == 0, hookSortable would make whole application stall
      // because the body of the table can never be rendered successfully
      if (data.length > 0) {
        // Here goes the magic
        hookSortable();
      }
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

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

function handlePositiveClick(): boolean {
  loading.value = true;
  if (currentHeaderId.value == null) {
    window.$notifyError(t('message.noChosenHeaderError'));
    return;
  }
  const updateParam: dto.DiffTableHeaderDto = {
    ID: currentHeaderId.value,
    LevelOrders: data.map((u: dto.DiffTableHeaderDto) => u.Level).join(",")
  } as any;
  UpdateHeaderLevelOrders(updateParam as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      show.value = false;
      emit('refresh');
    }).catch(err => window.$notifyError(err)).finally(() => loading.value = false);
  return false
}

function handleNegativeClick() {
  show.value = false;
}
</script>
