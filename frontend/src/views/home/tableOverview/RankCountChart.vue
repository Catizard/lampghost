<template>
  <!-- This page is a copy of LampCountChart. Might be merged into one component in the future -->
  <n-data-table :columns="overviewColumns" :data="overviewData" />
  <vue-apex-charts :height="chartHeight()" type="bar" :options="rankCountChartOptions" :series="rankCountChartSeries" />
</template>

<script setup lang="ts">
import { ScoreRank } from "@/constants/scoreRank";
import { dto } from "@wailsjs/go/models";
import { DataTableColumns } from "naive-ui";
import { reactive, ref, Ref, watch } from "vue";
import VueApexCharts from "vue3-apexcharts";

const props = defineProps<{
  data: dto.DiffTableHeaderDto | null
  rivalId?: number;
  tableId?: number;
  dataOrder: "reverse" | "ordered"
}>();

// TODO: I need a better solution, eh
const RANKS = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
const REVERSE_RANKS = [9, 8, 7, 6, 5, 4, 3, 2, 1, 0];

const STR_RANKS = [
  "NO_PLAY",
  "F",
  "E",
  "D",
  "C",
  "B",
  "A",
  "AA",
  "AAA",
  "MAX",
];
const STR_REVERSE_RANKS = [
  "NO_PLAY",
  "MAX",
  "AAA",
  "AA",
  "A",
  "B",
  "C",
  "D",
  "E",
  "F",
];

const REVERSE_RANK_COLOR = [
  "#FFFFFF",
  "#CC5C76",
  "#CC5C76",
  "#49E670",
  "#4FBCF7",
  "#FF6B74",
  "#FFAD70",
  "#FFD251",
  "#FFFF21",
  "#FFE270",
];
const RANK_COLOR = [
  "#FFFFFF",
  "#FFD251",
  "#FFD251",
  "#FFD251",
  "#FFAD70",
  "#FF6B74",
  "#4FBCF7",
  "#49E670",
  "#CC5C76",
  "#CC5C76",
];

// NOTE: ApexCharts Vue doesn't give any responses if we changed rankCountChartOptions.colors
// The only way to correctly reassign the colors is by rebuilding whole rankCountChartOptions
// I believe this is a ApexCharts Vue bug or design failure
let rankCountChartOptions = reactive(buildRankCountChartOptions("reverse"));
function buildRankCountChartOptions(order: "reverse" | "ordered") {
  return {
    chart: {
      type: "bar",
      stacked: true,
      stackType: "100%",
    },
    colors: order == "reverse" ? REVERSE_RANK_COLOR : RANK_COLOR,
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
  };
}
let rankCountChartSeries = [];

function buildSeries(data: dto.DiffTableHeaderDto, order: "reverse" | "ordered") {
  const rankNames = order == "reverse" ? STR_REVERSE_RANKS : STR_RANKS;
  const rankValues = order == "reverse" ? REVERSE_RANKS : RANKS;
  rankCountChartOptions = buildRankCountChartOptions(order);
  rankCountChartSeries = rankNames.map((rankName) => {
    return {
      name: rankName,
      data: [],
    };
  });
  rankCountChartOptions.xaxis.categories = [];
  const header = data;
  for (let i = 0; i < header.SortedLevels.length; ++i) {
    const level = header.SortedLevels[i];
    const levelName = header.Symbol + level;
    rankCountChartOptions.xaxis.categories.push(levelName);
    const dataList = header.LevelLayeredContents[level];
    for (let j = 0; j < rankValues.length; ++j) {
      const rankValue = rankValues[j];
      const count = dataList.filter((data) => data.ScoreRank == rankValue).length;
      let v = {
        x: levelName,
        y: count,
        fillColor: null,
        strokeColor: null,
      };
      if (rankValue == 0) {
        v.fillColor = "#FFFFFF";
        v.strokeColor = "#FFFFFF";
      }
      rankCountChartSeries[j].data.push(v);
    }
  }
}

const baseHeight = 400;
function chartHeight(): string {
  const addition = rankCountChartOptions.xaxis.categories.length * 5;
  return `${baseHeight + addition}px`;
}

const overviewColumns: DataTableColumns<dto.DiffTableHeaderDto> = [
  {
    title: "No Play",
    key: ScoreRank.NO_PLAY,
    width: "60px",
  },
  {
    title: "F",
    key: ScoreRank.F,
    width: "60px",
  },
  {
    title: "E",
    key: ScoreRank.E,
    width: "60px",
  },
  {
    title: "D",
    key: ScoreRank.D,
    width: "60px",
  },
  {
    title: "C",
    key: ScoreRank.C,
    width: "60px",
  },
  {
    title: "B",
    key: ScoreRank.B,
    width: "60px"
  },
  {
    title: "A",
    key: ScoreRank.A,
    width: "60px",
  },
  {
    title: "AA",
    key: ScoreRank.AA,
    width: "60px",
  },
  {
    title: "AAA",
    key: ScoreRank.AAA,
    width: "60px",
  },
  {
    title: "MAX",
    key: ScoreRank.MAX,
    width: "60px",
  }
];
// overviewData is a two-dimensional array
// overviewData[0]: the percentage for each clear type
// overviewData[1]: a string for each clear type, which structure is "count/total"
const overviewData: Ref<Array<any>> = ref([]);

function buildOverviewTable(data: dto.DiffTableHeaderDto, order: "reverse" | "ordered") {
  const rankCount = data.RankCount;
  const songCount = data.SongCount ?? 0;
  let sum = 0;
  overviewData.value = [{}, {}];
  const reversedScoreRankKeys = Object.keys(ScoreRank);
  for (let i = reversedScoreRankKeys.length - 1; i >= 0; --i) {
    // This is not cool typescript :(
    const v: any = ScoreRank[reversedScoreRankKeys[i]];
    const count = rankCount[v] ?? 0;
    sum += count;
    // We should continue to calculate sum here
    if (v == ScoreRank.NO_PLAY) {
      continue;
    }
    overviewData.value[0][v] = `${((100 * count) / songCount).toFixed(2)}%`;
    overviewData.value[1][v] = `${count}/${songCount}`;
  }
  // We cannot have NO_PLAY count directly
  overviewData.value[0][ScoreRank.NO_PLAY] = `${((100 * (songCount - sum)) / songCount).toFixed(2)}%`;
  overviewData.value[1][ScoreRank.NO_PLAY] = `${songCount - sum}/${songCount}`;
}

function rebuild(data: dto.DiffTableHeaderDto, order: "reverse" | "ordered") {
  if (props.data == null) {
    return;
  }
  try {
    buildSeries(data, order);
    buildOverviewTable(data, order);
  } catch (err) {
    window.$notifyError(err);
  }
}

watch([() => props.data, () => props.dataOrder], ([newData, newOrder]) => {
  rebuild(newData, newOrder);
}, { immediate: true });
</script>
