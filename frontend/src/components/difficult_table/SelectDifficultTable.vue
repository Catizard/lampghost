<template>
  <n-select :loading="loading" v-model:value="tableId" :options="tableOptions" />
</template>

<script setup lang="ts">
import { FindDiffTableHeaderList } from '@wailsjs/go/main/App';
import { dto, result } from '@wailsjs/go/models';
import { SelectOption } from 'naive-ui';
import { onMounted, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const loading = ref(false);
const tableId = defineModel<number | null>("value");
const tableOptions: Ref<SelectOption[]> = ref([]);

interface Props {
  slientWhenNoTable?: boolean
};
const { slientWhenNoTable = false } = defineProps<Props>();

function loadData() {
  loading.value = true;
  FindDiffTableHeaderList()
    .then((result: result.RtnDataList) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Rows.length == 0) {
        if (!slientWhenNoTable) {
          return Promise.reject(t("message.noTableError"))
        }
        return Promise.resolve();
      }
      tableOptions.value = result.Rows.map((header: dto.DiffTableHeaderDto) => {
        return {
          label: header.Name,
          value: header.ID,
        } as SelectOption;
      });
      tableId.value = tableOptions.value[0].value as number;
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

onMounted(() => {
  loadData();
});
</script>
