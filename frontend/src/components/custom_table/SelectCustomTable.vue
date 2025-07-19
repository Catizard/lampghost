<template>
  <n-select :loading="loading" v-model:value="customTableId" :options="customTableOptions" />
</template>

<script setup lang="ts">
import { useSelectMemo } from '@/stores/selectMemo';
import { FindCustomDiffTableList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { SelectOption } from 'naive-ui';
import { onMounted, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const customTableId = defineModel<number | null>("value");
const customTableOptions: Ref<SelectOption[]> = ref([]);
const loading = ref(false);
const selectMemoStore = useSelectMemo();

interface Props {
  ignoreDefaultTable?: boolean
};

const { ignoreDefaultTable = false } = defineProps<Props>();

function loadCustomTableOptions() {
  loading.value = true;
  FindCustomDiffTableList({
    IgnoreDefaultTable: ignoreDefaultTable
  } as any).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    if (result.Rows.length == 0) {
      return Promise.reject(t('message.noCustomTableError'))
    }
    customTableOptions.value = result.Rows.map((row: dto.CustomDiffTableDto): SelectOption => {
      return {
        label: row.Name,
        value: row.ID
      }

    });
    const memoId = selectMemoStore.$state.customTableId;
    if (memoId != null && memoId != 0 && result.Rows.find((customTable: dto.CustomDiffTableDto) => customTable.ID == memoId)) {
      customTableId.value = memoId;
    } else {
      customTableId.value = customTableOptions.value[0].value as number;
    }
  })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

onMounted(() => {
  loadCustomTableOptions();
});

watch(customTableId, newId => {
  selectMemoStore.setCustomTableId(newId);
})
</script>
