<template>
  <n-flex justify="flex-start">
    <n-radio-group v-model:value="currentDisplayType" name="radiogroup">
      <n-radio-button v-for="dt in displayTypes" :key="dt.type" :value="dt.type" :label="dt.label" />
    </n-radio-group>
    <n-select v-model:value="currentTableId" :options="tableOptions" style="width: 150px;" />
  </n-flex>
  <vue-apex-charts height="450px" type="bar" :options="lampCountChartOptions" :series="lampCountChartOptions.series" />
</template>

<script setup lang="ts">
import { FindDiffTableHeaderList } from "@wailsjs/go/controller/DiffTableController";
import { QueryUserInfoWithLevelLayeredDiffTableLampStatus } from "@wailsjs/go/controller/RivalInfoController";
import { dto } from "@wailsjs/go/models";
import { SelectOption, useNotification } from "naive-ui";
import { reactive, ref, Ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import VueApexCharts from "vue3-apexcharts";

const props = defineProps<{
  rivalId?: number
}>();

const { t } = useI18n();
const notification = useNotification();

const LAMPS = [1, 4, 5, 11, 0];
const STR_LAMPS = ["FAILED", "Easy", "Normal", "Hard+", "NO_PLAY"];

const currentTableId: Ref<number | null> = ref(null);
const currentHeader: Ref<dto.DiffTableHeaderDto | null> = ref(null);
const tableOptions: Ref<Array<SelectOption>> = ref([]);
function loadTableData() {
  FindDiffTableHeaderList()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Rows.length == 0) {
        return Promise.reject(t('message.noTableError'));
      }
      tableOptions.value = result.Rows.map((header: dto.DiffTableHeaderDto) => {
        return {
          label: header.Name,
          value: header.ID,
        } as SelectOption
      });
      currentTableId.value = tableOptions.value[0].value as number;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      })
    });
}
loadTableData();

interface DisplayType {
  type: number,
  label: string
}
const displayTypes: Array<DisplayType> = [
  {
    type: 0,
    label: "Count",
  },
  {
    type: 1,
    label: "Rate"
  }
];
const currentDisplayType: Ref<number | null> = ref(0);

const lampCountChartOptions = reactive({
  chart: {
    type: "bar",
    stacked: true,
    stackType: "100%",
  },
  colors: ["#FF0000", "#00FF00", "#0000FF", "#00FFFF", "#FFFFFF"],
  series: [],
  plotOptions: {
    bar: {
      horizontal: true,
      columnHeight: "100%",
    },
  },
  xaxis: {
    categories: [],
  },
  fill: {
    opacity: 1,
  },
  legend: {
    position: "top",
    horizontalAlign: "left",
    offsetX: 40,
  },
  dataLabels: {
    enabled: false,
  },
});
function loadLampData() {
  QueryUserInfoWithLevelLayeredDiffTableLampStatus(props.rivalId, currentTableId.value)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      currentHeader.value = result.Data.DiffTableHeader;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      })
    });
}

function buildSeries() {
  lampCountChartOptions.series = STR_LAMPS.map(lampName => {
    return {
      name: lampName,
      data: [],
    }
  });
  lampCountChartOptions.xaxis.categories = [];
  const header = currentHeader.value;
  for (let i = 0; i < header.SortedLevels.length; ++i) {
    const level = header.SortedLevels[i];
    const levelName = header.Symbol + level;
    lampCountChartOptions.xaxis.categories.push(levelName);
    const dataList = header.LevelLayeredContents[level];
    for (let j = 0; j < LAMPS.length; ++j) {
      const lampValue = LAMPS[j];
      // TODO: wtf
      const count = dataList.filter(
        (data) =>
          data.Lamp == lampValue ||
          (data.Lamp >= 6 && lampValue == 11) ||
          (lampValue == 1 && data.Lamp < 4 && data.Lamp > 0),
      ).length;
      let v = {
        x: levelName,
        y: count,
        fillColor: null,
        strokeColor: null,
      };
      if (lampValue == 0) {
        v.fillColor = "#FFFFFF";
        v.strokeColor = "#FFFFFF";
      }
      lampCountChartOptions.series[j].data.push(v);
    }
  }
}

watch([currentTableId, () => props.rivalId], ([newTableId, newRivalId]) => {
  if (newTableId == null || newTableId == undefined || newRivalId == null || newRivalId == undefined) {
    return;
  }
  loadLampData();
});

watch(currentHeader, _ => {
  buildSeries();
})
</script>

<i18n lang="json">{
  "en": {
    "message": {
      "noTableError": "Cannot handle no difficult table data currenlty, please add at least one table first"
    }
  },
  "zh-CN": {
    "message": {
      "noTableError": "目前无法处理一个难度表都没有的情况，请至少先添加一个难度表"
    }
  }
}</i18n>