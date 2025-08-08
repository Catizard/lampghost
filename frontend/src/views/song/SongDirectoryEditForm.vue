<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.editSongDirectories')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-button type="primary" @click="chooseBMSDirectory">
      {{ t('button.chooseDirectory') }}
    </n-button>
    <DirectoryTable v-model:directories="directories" />
  </n-modal>
</template>

<script setup lang="ts">
import { FindSongDirectories, OpenDirectoryDialog, SaveSongDirectories } from '@wailsjs/go/main/App';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import DirectoryTable from '../initialize/directoryTable.vue';
import { entity } from '@wailsjs/go/models';

const { t } = useI18n();
const show = ref(false);
const loading = ref(false);
const directories = ref([]);
defineExpose({ open });

function open() {
  show.value = true;
  loadData();
}

function loadData() {
  loading.value = true;
  FindSongDirectories().then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    directories.value = [...result.Rows.map((directory: entity.SongDirectory) => directory.DirectoryPath)];
  }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function handlePositiveClick() {
  loading.value = true;
  SaveSongDirectories(directories.value).then(result => {
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

function chooseBMSDirectory() {
  OpenDirectoryDialog(t('title.chooseBMSDirectory')).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    const path = result.Data;
    if (path != "") {
      directories.value.push(path);
    }
  }).catch(err => window.$notifyError(err))
}
</script>
