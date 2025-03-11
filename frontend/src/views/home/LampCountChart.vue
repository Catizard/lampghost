<template>
  <n-select v-model:value="currentTableId" :options="tableOptions" width="150px" />
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
      const header: dto.DiffTableHeaderDto = result.Data.DiffTableHeader;
      lampCountChartOptions.series = [];
      STR_LAMPS.forEach((lampName) => {
        lampCountChartOptions.series.push({
          name: lampName,
          data: [],
        });
      });
      lampCountChartOptions.xaxis.categories = [];
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
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      })
    });
}

watch([currentTableId, () => props.rivalId], ([newTableId, newRivalId]) => {
  if (newTableId == null || newTableId == undefined || newRivalId == null || newRivalId == undefined) {
    return;
  }
  loadLampData();
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
}

</i18n>