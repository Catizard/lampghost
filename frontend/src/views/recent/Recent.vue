<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.recentPlay') }}</n-text>
    </n-h1>
  </n-flex>
  <n-flex justify="start">
    <n-input :placeholder="t('form.placeholderSearchSongOrSabunName')" v-model:value="searchNameLike"
      @keyup.enter="tableRef.loadData()" style="width: 350px;" />
    <!-- <n-button>{{ t('button.chooseClearType') }}</n-button>
      <n-button>{{ t('button.minimumClearType') }}</n-button> -->
  </n-flex>
  <RecentTable ref="tableRef" :songNameLike="searchNameLike" />
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import RecentTable from './RecentTable.vue';

const i18n = useI18n();
const { t } = i18n;

const tableRef = ref<InstanceType<typeof RecentTable>>(null);
const searchNameLike = ref("");
</script>
