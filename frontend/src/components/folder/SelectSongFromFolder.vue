<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.selectSongFromFolder')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <FolderTable :customTableId="customTableId" type="table" selectSong="single"
      v-model:checkedRowKeys="checkedRowKeys" />
  </n-modal>
</template>

<script setup lang="ts">
import { Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import FolderTable from '@/views/folder/FolderTable.vue';
import { BindFolderContentToCustomCourse } from '@wailsjs/go/main/App';

const props = defineProps<{
  customTableId: number
}>();

const show = defineModel<boolean>("show");

const { t } = useI18n();
const loading = ref(false);
const checkedRowKeys: Ref<number[]> = ref([]);

function handlePositiveClick() {
  if (checkedRowKeys.value.length != 1) {
    window.$notifyError(t('message.noSelectedSong'));
    return
  }
  loading.value = true;
  BindFolderContentToCustomCourse(checkedRowKeys[0], props.customTableId)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      show.value = false;
    }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function handleNegativeClick() {
  show.value = false;
}
</script>
