<template>
  <n-grid :cols="12">
    <n-gi :span="8">
      <div>
    <apexchart
      height="108px"
      type="bar"
      :options="chartOptions"
      :series="series"
    ></apexchart>
  </div>
    </n-gi>
    <n-gi :span="4">
      <div class="green" />
    </n-gi>
  </n-grid>
</template>

<script lang="ts">
import ApexCharts from 'apexcharts';
import { reactive, ref, toRefs } from 'vue';

export default {
  components: {
    ApexCharts,
  },
  setup(props, ctx) {
    const datas = reactive({
      chartOptions: {
        chart: {
          id: "vuechart-example",
          type: 'line'
        },
        xaxis: {
          categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep'],
        },
      },
      series: [
        {
          name: "series 1",
          data: [30, 40, 35, 50, 49, 60, 70, 91, 125]
        }
      ],
    });

    return {
      ...toRefs(datas)
    }
  },
}
</script>

<style scoped>
.light-green {
  height: 108px;
  background-color: rgba(0, 128, 0, 0.12);
}

.green {
  height: 108px;
  background-color: rgba(0, 128, 0, 0.24);
}
</style>