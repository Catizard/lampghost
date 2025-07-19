<!-- Similar to SelectFolder, but candidates are difficult folders from custom tables -->
<template>
  <n-modal :loading="loading" v-model:show="show" :title="t('title.bindToCustomTable')" preset="dialog"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" closable @close="() => { show = false }">
    <n-flex justify="space-between">
      <SelectCustomTable v-model:value="currentCustomTableID" style="width: 200px;" ignoreDefaultTable />
      <!-- <n-button type="primary" @click="handleClickAddFolder">{{ t('button.addDifficultFolder') }}</n-button> -->
    </n-flex>
    <SelectUnboundFolder ref="selectUnboundFolderRef" type="folder" v-model:checkedFolderIds="checkedFolderIds"
      :sha256="sha256" :customTableId="currentCustomTableID" />
  </n-modal>

  <FolderAddForm type="table" v-model:show="showAddModal" @refresh="reload" />
</template>

<script setup lang="ts">
import { Ref, ref, watch } from 'vue';
import { dto } from '@wailsjs/go/models';
import { useI18n } from 'vue-i18n';
import FolderAddForm from '../folder/FolderAddForm.vue';
import SelectUnboundFolder from '../folder/SelectUnboundFolder.vue';
import SelectCustomTable from '@/components/custom_table/SelectCustomTable.vue';

const { t } = useI18n();
const show = defineModel<boolean>("show");
const props = defineProps<{
  sha256?: string
}>();
const emit = defineEmits<{
  (e: 'submit', selected: Array<any>): void
}>();

const loading = ref(false);
const showAddModal = ref(false);
const data = ref<dto.FolderDto[]>([]);
let checkedFolderIds: Ref<number[]> = ref([]);
const selectUnboundFolderRef: Ref<InstanceType<typeof SelectUnboundFolder>> = ref(null);

const currentCustomTableID: Ref<number | null> = ref(null);

function reload() {
  selectUnboundFolderRef.value.reload();
}

watch(show, (newValue, oldValue) => {
  if (newValue == true) {
    watch(selectUnboundFolderRef, (r) => {
      r.reload();
    }, { once: true });
  }
}, { deep: true });

function handlePositiveClick() {
  emit('submit', checkedFolderIds.value);
}

function handleNegativeClick() {
  data.value = [];
  show.value = false;
}
</script>
