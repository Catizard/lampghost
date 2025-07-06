<template>
  <n-modal v-model:show="show" preset="dialog" :title="t('title.selectSongFromFolder')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false" style="width: 95vh; height: 95vh;">
    <FolderTable :customTableId="customTableId" type="table" selectSong="single" v-model:checkedRowKeys="checkedRowKeys"
      noActions :max-height="'75vh'" />
  </n-modal>
</template>

<script setup lang="ts">
import { Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import FolderTable from '@/views/folder/FolderTable.vue';

const props = defineProps<{
  customTableId?: number
}>();
const emit = defineEmits(['select']);

const show = defineModel<boolean>("show");

const { t } = useI18n();
const checkedRowKeys: Ref<number[]> = ref([]);

function handlePositiveClick() {
  if (checkedRowKeys.value.length == 0) {
    window.$notifyError(t('message.noSelectedSong'));
    return
  }
  emit('select', checkedRowKeys.value);
}

function handleNegativeClick() {
  show.value = false;
}
</script>
