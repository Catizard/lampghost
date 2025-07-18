<template>
  <n-select v-model:value="currentYear" :options="yearOptions" style="width: 100px;" />
  <vue-apex-charts height="100%" type="line" :options="chartOptions" :series="chartOptions.series" />
</template>

<script setup lang="ts">
import VueApexCharts from "vue3-apexcharts";
import { QueryRivalPlayedYears, QueryUserPlayCountInYear } from "@wailsjs/go/main/App";
import { SelectOption } from "naive-ui";
import { reactive, Ref, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useSpecifyYearStore } from "@/stores/home";

const props = defineProps<{
  rivalId?: number
}>();

const { t } = useI18n();

const currentYear: Ref<string | null> = ref(null);
const yearOptions: Ref<Array<SelectOption>> = ref([
  {
    label: "2024",
    key: "2024",
  },
  {
    label: "2023",
    key: "2023",
  },
  {
    label: "2022",
    keys: "2022",
  }
]);


const chartOptions = reactive({
  chart: {
    id: "chart-play-count",
    type: "line",
    zoom: {
      enabled: false,
      allowMouseWheelZoom: false,
    },
  },
  xaxis: {
    categories: [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ],
  },
  series: [
    {
      name: t("title.playCount"),
      data: [],
    },
  ],
});

function loadPlayedYears() {
  QueryRivalPlayedYears(props.rivalId)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      yearOptions.value = result.Rows.map((year: number) => {
        return {
          label: year.toString(),
          value: year.toString()
        } as SelectOption
      });
      if (yearOptions.value.length == 0) {
        return Promise.reject(t('message.noRecordError'));
      }
      currentYear.value = yearOptions.value[yearOptions.value.length - 1].value as string;
    }).catch(err => window.$notifyError(err))
}

function loadPlayCountData() {
  QueryUserPlayCountInYear(props.rivalId, currentYear.value)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      chartOptions.series[0].data = [...result.Rows];
    }).catch(err => window.$notifyError(err))
}

watch(() => props.rivalId, (newId) => {
  if (newId != null && newId != undefined) {
    loadPlayedYears();
  }
});

const specifyYearStore = useSpecifyYearStore();
watch(currentYear, (newValue) => {
  console.log('setting new currentYear:', newValue);
  loadPlayCountData();
  specifyYearStore.setter(newValue);
});

</script>
