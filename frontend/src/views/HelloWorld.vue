<template>
  <perfect-scrollbar>
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">
        也许是玩家信息
      </n-text>
    </n-h1>
    <n-grid :cols="12">
      <n-gi :span="4">
        <n-flex>
          Player Name: {{ playerData.playerName }}
          <n-divider />
          Total Play: {{ playerData.playCount }}
          <n-divider />
          Play Time : {{ playerData.playTime }}
          <n-divider />
          <n-button>同步最新存档数据,大概</n-button>
        </n-flex>
      </n-gi>
      <n-gi :span="8">
        <vue-apex-charts height="100%" type="line" :options="playCountChartOptions"
          :series="playCountChartOptions.series" />
      </n-gi>
    </n-grid>
    <n-divider />
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">
        报告,我是点灯信息
      </n-text>
    </n-h1>
    <n-grid :cols="12" min-height="200px">
      <n-gi :span="4">
        这里应该有个难度表选择列表？
        <perfect-scrollbar style="height: 350px">
          <n-flex vertical>
            <n-button>発狂BMS難易度表</n-button>
            <n-button>Satellite</n-button>
            <n-button>NEW GENERATION 発狂難易度表</n-button>
            <n-button>発狂BMS難易度表</n-button>
            <n-button>Satellite</n-button>
            <n-button>NEW GENERATION 発狂難易度表</n-button>
            <n-button>発狂BMS難易度表</n-button>
            <n-button>Satellite</n-button>
            <n-button>NEW GENERATION 発狂難易度表</n-button>
          </n-flex>
        </perfect-scrollbar>
      </n-gi>
      <n-gi :span="8">
        <vue-apex-charts height="350px" type="bar" :options="lampCountChartOptions"
          :series="lampCountChartOptions.series" />
      </n-gi>
    </n-grid>
  </perfect-scrollbar>
</template>

<script setup>
import VueApexCharts from 'vue3-apexcharts';
import { reactive, ref, toRefs } from 'vue';
import { QueryUserInfo } from '../../wailsjs/go/controller/RivalInfoController'

const playCountChartOptions = ref({
  chart: {
    id: "chart-play-count",
    type: 'line'
  },
  xaxis: {
    categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep'],
  },
  series: [
    {
      name: "series 1",
      data: [30, 40, 35, 50, 49, 60, 70, 91, 125]
    }
  ]
});

const lampCountChartOptions = ref({
  chart: {
    type: 'bar',
    stacked: true,
  },
  series: [{
    name: 'Marine Sprite',
    data: [44, 55, 41, 37, 22, 43, 21]
  }, {
    name: 'Striking Calf',
    data: [53, 32, 33, 52, 13, 43, 32]
  }, {
    name: 'Tank Picture',
    data: [12, 17, 11, 9, 15, 11, 20]
  }, {
    name: 'Bucket Slope',
    data: [9, 7, 5, 8, 6, 9, 4]
  }, {
    name: 'Reborn Kid',
    data: [25, 12, 19, 32, 25, 24, 10]
  }],
  plotOptions: {
    bar: {
      horizontal: true,
      dataLabels: {
        total: {
          enabled: true,
          offsetX: 0,
          style: {
            fontSize: '13px',
            fontWeight: 900
          }
        }
      }
    },
  },
  stroke: {
    width: 1,
    colors: ['#fff']
  },
  title: {
    text: '我觉得我是当前难度表的进度'
  },
  xaxis: {
    categories: [2008, 2009, 2010, 2011, 2012, 2013, 2014],
    labels: {
      formatter: function (val) {
        return val + "K"
      }
    }
  },
  yaxis: {
    title: {
      text: undefined
    },
  },
  tooltip: {
    y: {
      formatter: function (val) {
        return val + "K"
      }
    }
  },
  fill: {
    opacity: 1
  },
  legend: {
    position: 'top',
    horizontalAlign: 'left',
    offsetX: 40
  }
});

const playerData = reactive({
  playerName: "Catizard",
  playCount: 114514,
  playTime: "24:17:30"
})

function initUser() {
  QueryUserInfo().then(result => {
    playerData.playerName = result.Name
  })
}

initUser();
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

.n-button {
  width: 80%
}

.ps {
  height: 100%;
}
</style>