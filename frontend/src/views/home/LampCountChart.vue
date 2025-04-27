<template>
  <n-flex justify="flex-start">
    <n-select v-model:value="currentTableId" :options="tableOptions" style="width: 150px" />
    <n-radio-group v-model:value="currentLampOrder">
      <n-radio-button :key='"reverse"' :value='"reverse"' :label="t('radio.reverse')" />
      <n-radio-button :key='"ordered"' :value='"ordered"' :label="t('radio.ordered')" />
    </n-radio-group>
  </n-flex>
  <n-data-table :columns="overviewColumns" :data="overviewData" :row-class-name="'test'" />
  <vue-apex-charts :height="chartHeight()" type="bar" :options="lampCountChartOptions"
    :series="lampCountChartOptions.series" />
</template>

<script setup lang="ts">
import { ClearType } from "@/constants/cleartype";
import { FindDiffTableHeaderList } from "@wailsjs/go/main/App";
import { QueryUserInfoWithLevelLayeredDiffTableLampStatus } from "@wailsjs/go/main/App";
import { dto } from "@wailsjs/go/models";
import { DataTableColumns, SelectOption, useNotification } from "naive-ui";
import { computed, reactive, ref, Ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import VueApexCharts from "vue3-apexcharts";

const props = defineProps<{
  rivalId?: number;
}>();

const { t } = useI18n();
const notification = useNotification();

// TODO: I need a better solution, err
const LAMPS = [1, 4, 5, 6, 7, 8, 9, 10, 0];
const REVERSE_LAMPS = [10, 9, 8, 7, 6, 5, 4, 1, 0];
const STR_LAMPS = [
  "FAILED",
  "Easy",
  "Normal",
  "Hard",
  "EXH",
  "FC",
  "PERFECT",
  "MAX",
  "NO_PLAY",
];
const REVERSE_STR_LAMPS = [
  "MAX",
  "PERFECT",
  "FC",
  "EXH",
  "Hard",
  "Normal",
  "Easy",
  "FAILED",
  "NO_PLAY"
];
const LAMP_COLOR = [
  "#CC5C76",
  "#49E670",
  "#4FBCF7",
  "#FF6B74",
  "#FFAD70",
  "#FFD251",
  "#FFD251",
  "#FFD251",
  "#FFFFFF",
];
const REVERSE_LAMP_COLOR = [
  "#FFD251",
  "#FFD251",
  "#FFD251",
  "#FFAD70",
  "#FF6B74",
  "#4FBCF7",
  "#49E670",
  "#CC5C76",
  "#FFFFFF",
];

const currentLampOrder: Ref<"reverse" | "ordered"> = ref("reverse")
const currentTableId: Ref<number | null> = ref(null);
const currentHeader: Ref<dto.DiffTableHeaderDto | null> = ref(null);
const tableOptions: Ref<Array<SelectOption>> = ref([]);
function loadTableData() {
  FindDiffTableHeaderList()
    .then((result) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Rows.length == 0) {
        return Promise.reject(t("message.noTableError"));
      }
      tableOptions.value = result.Rows.map((header: dto.DiffTableHeaderDto) => {
        return {
          label: header.Name,
          value: header.ID,
        } as SelectOption;
      });
      currentTableId.value = tableOptions.value[0].value as number;
    })
    .catch((err) => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      });
    });
}
loadTableData();

const lampCountChartOptions = reactive({
  chart: {
    type: "bar",
    stacked: true,
    stackType: "100%",
  },
  colors: [

  ],
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
    enabled: true,
  },
});
function loadLampData() {
  QueryUserInfoWithLevelLayeredDiffTableLampStatus(
    props.rivalId,
    currentTableId.value,
  ).then((result) => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    currentHeader.value = result.Data.DiffTableHeader;
  }).catch((err) => {
    notification.error({
      content: err,
      duration: 3000,
      keepAliveOnHover: true,
    });
  });
}

function buildSeries() {
  const lampNames = currentLampOrder.value == "reverse" ? REVERSE_STR_LAMPS : STR_LAMPS;
  const lampValues = currentLampOrder.value == "reverse" ? REVERSE_LAMPS : LAMPS;
  lampCountChartOptions.colors = (currentLampOrder.value == "reverse" ? REVERSE_LAMP_COLOR : LAMP_COLOR);
  lampCountChartOptions.series = lampNames.map((lampName) => {
    return {
      name: lampName,
      data: [],
    };
  });
  lampCountChartOptions.xaxis.categories = [];
  const header = currentHeader.value;
  for (let i = 0; i < header.SortedLevels.length; ++i) {
    const level = header.SortedLevels[i];
    const levelName = header.Symbol + level;
    lampCountChartOptions.xaxis.categories.push(levelName);
    const dataList = header.LevelLayeredContents[level];
    for (let j = 0; j < lampValues.length; ++j) {
      const lampValue = lampValues[j];
      // TODO: wtf
      const count = dataList.filter(
        (data) =>
          data.Lamp == lampValue ||
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

const baseHeight = 400;
function chartHeight(): string {
  const addition = lampCountChartOptions.xaxis.categories.length * 5;
  return `${baseHeight + addition}px`;
}

const overviewColumns: DataTableColumns<dto.DiffTableHeaderDto> = [
  {
    title: "No Play",
    key: ClearType.NO_PLAY,
    width: "60px",
  },
  {
    title: "Failed",
    key: ClearType.Failed,
    width: "60px",
  },
  {
    title: "Easy Clear+",
    key: ClearType.Easy,
    width: "60px",
  },
  {
    title: "Normal Clear+",
    key: ClearType.Normal,
    width: "60px",
  },
  {
    title: "Hard Clear+",
    key: ClearType.Hard,
    width: "60px",
  },
  {
    title: "EX Hard Clear+",
    key: ClearType.ExHard,
    width: "60px",
  },
  {
    title: "Full Combo+",
    key: ClearType.FullCombo,
    width: "60px",
  },
];
// overviewData is a two-dimensional array
// overviewData[0]: the percentage for each clear type
// overviewData[1]: a string for each clear type, which structure is "count/total"
const overviewData: Ref<Array<any>> = ref([]);

function buildOverviewTable() {
  const lampCount = currentHeader.value.LampCount;
  const songCount = currentHeader.value.SongCount ?? 0;
  let sum = 0;
  overviewData.value = [{}, {}];
  // TODO: currently, Assist/Light Assist would be calculated as fail
  // Some of the ClearType has special calculate way, I haven't found a good way to write this code
  const failCount = (lampCount[ClearType.Failed] ?? 0)
    + (lampCount[ClearType.AssistEasy] ?? 0)
    + (lampCount[ClearType.LightAssistEasy] ?? 0);
  overviewData.value[0][ClearType.Failed] = `${(100 * failCount / songCount).toFixed(2)}%`;
  overviewData.value[1][ClearType.Failed] = `${failCount}/${songCount}`
  const reversedClearTypeKeys = Object.keys(ClearType)
  for (let i = reversedClearTypeKeys.length - 1; i >= 0; --i) {
    // This is not cool typescript :(
    const v: any = ClearType[reversedClearTypeKeys[i]];
    const count = lampCount[v] ?? 0;
    sum += count;
    // We should continue to calculate sum here
    if (v == ClearType.LightAssistEasy || v == ClearType.AssistEasy || v == ClearType.Failed || v == ClearType.NO_PLAY) {
      continue;
    }
    overviewData.value[0][v] = `${(100 * sum / songCount).toFixed(2)}%`;
    overviewData.value[1][v] = `${sum}/${songCount}`
  }
  // We cannot have NO_PLAY count directly
  overviewData.value[0][ClearType.NO_PLAY] = `${(100 * (songCount - sum) / songCount).toFixed(2)}%`;
  overviewData.value[1][ClearType.NO_PLAY] = `${songCount - sum}/${songCount}`
}

watch([currentTableId, () => props.rivalId], ([newTableId, newRivalId]) => {
  if (newTableId == null || newTableId == undefined) {
    return;
  }
  if (newRivalId == null || newRivalId == undefined) {
    return;
  }
  loadLampData();
});

watch([currentHeader, currentLampOrder], (_) => {
  buildSeries();
  buildOverviewTable();
});
</script>

<i18n lang="json">{
  "en": {
    "message": {
      "noTableError": "Cannot handle no difficult table data currenlty, please add at least one table first"
    },
    "radio": {
      "reverse": "Reverse",
      "ordered": "Ordered"
    }
  },
  "zh-CN": {
    "message": {
      "noTableError": "目前无法处理一个难度表都没有的情况，请至少先添加一个难度表"
    },
    "radio": {
      "reverse": "倒序",
      "ordered": "正序"
    }
  }
}</i18n>
