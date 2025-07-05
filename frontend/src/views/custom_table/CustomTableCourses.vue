<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.customTableCourses') }}</n-text>
    </n-h1>
  </n-flex>
  <n-spin :show="loading">
    <n-flex justify="space-between">
      <SelectCustomTable v-model:value="currentCustomTableID" style="width: 200px;" />
    </n-flex>
  </n-spin>
</template>

<script setup lang="ts">
import { ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const loading = ref(false);
const currentCustomTableID: Ref<number | null> = ref(null);
</script>
