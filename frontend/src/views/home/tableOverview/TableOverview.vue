<template>
  <n-h1 prefix="bar" style="text-align: start">
    <n-text type="primary">
      {{ t("title.overview") }}
    </n-text>
  </n-h1>
  <n-flex justify="start">
    <SelectDifficultTable v-model:value="currentTableId" style="width: 150px;" />
    <n-radio-group v-model:value="overviewType">
      <n-radio-button :key="'lamp'" :value="'lamp'" :label="t('button.lamp')" />
      <n-radio-button :key="'score'" :value="'score'" :label="t('button.score')" />
    </n-radio-group>
    <n-radio-group v-model:value="dataOrder">
      <n-radio-button :key="'reverse'" :value="'reverse'" :label="t('button.reverse')" />
      <n-radio-button :key="'ordered'" :value="'ordered'" :label="t('button.ordered')" />
    </n-radio-group>
  </n-flex>
  <LampCountChart v-if="overviewType == 'lamp'" :data="header" :rivalId="rivalId" :tableId="currentTableId"
    :dataOrder="dataOrder" />
  <RankCountChart v-if="overviewType == 'score'" :data="header" :rivalId="rivalId" :tableId="currentTableId"
    :dataOrder="dataOrder" />
</template>

<script setup lang="ts">
import SelectDifficultTable from "@/components/difficult_table/SelectDifficultTable.vue";
import { ref, Ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import LampCountChart from "./LampCountChart.vue";
import { QueryUserInfoWithLevelLayeredDiffTable } from "@wailsjs/go/main/App";
import { dto } from "@wailsjs/go/models";
import RankCountChart from "./RankCountChart.vue";

const { t } = useI18n();
const currentTableId: Ref<number | null> = ref(null);
const overviewType: Ref<"lamp" | "score"> = ref("lamp");
const dataOrder: Ref<"reverse" | "ordered"> = ref("reverse");
const header: Ref<dto.DiffTableHeaderDto> = ref(null);
const { rivalId } = defineProps<{
  rivalId?: number
}>();

function loadData() {
  QueryUserInfoWithLevelLayeredDiffTable(
    rivalId,
    currentTableId.value,
  )
    .then((result) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      header.value = result.Data.DiffTableHeader;
    })
    .catch(err => window.$notifyError(err));
}

watch([currentTableId, () => rivalId], ([newTableId, newRivalId]) => {
  if (newTableId == null || newTableId == undefined) {
    return;
  }
  if (newRivalId == null || newRivalId == undefined) {
    return;
  }
  console.log(newTableId, newRivalId);
  loadData();
});
</script>
