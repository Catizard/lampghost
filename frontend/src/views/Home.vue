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
          你是?: {{ playerData.playerName }}
          <n-divider />
          玩了多少次: {{ playerData.playCount }}
          <n-divider />
          最后同步时间 : {{ playerData.lastUpdate }}
          <n-divider />
          <n-button @click="handleSyncClick" :loading="syncLoading">
            同步最新存档数据,大概
          </n-button>
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
import { reactive, ref } from 'vue';
import { QueryUserInfoByID, SyncRivalScoreLog, QueryUserPlayCountInYear } from '../../wailsjs/go/controller/RivalInfoController'
import { useNotification } from 'naive-ui';
import dayjs from 'dayjs';

const notification = useNotification();

const playCountChartOptions = ref({
  chart: {
    id: "chart-play-count",
    type: 'line'
  },
  xaxis: {
    categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
  },
  series: [
    {
      name: "游玩次数",
      data: [30, 40, 35, 50, 49, 60, 70, 91, 125, 0, 0, 0]
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
  lastUpdate: "",
})

function initUser() {
  // TODO: Remove another magic 1
  QueryUserInfoByID(1).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg)
    }
    const { Data } = result;
    console.log(Data)
    playerData.playerName = Data.Name
    playerData.playCount = Data.PlayCount
    playerData.lastUpdate = dayjs(Data.UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
    console.log(playerData.lastUpdate)
    return QueryUserPlayCountInYear(1, 2024)
  }).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg)
    }
    const { Rows } = result;
    console.log(Rows);
    playCountChartOptions.value.series[0].data = [...Rows];
  }).catch(err => {
    notification.error({
      content: "获取用户数据失败: " + err,
      duration: 3000,
      keepAliveOnHover: true
    })
  })
}

const syncLoading = ref(false);
function handleSyncClick() {
  // TODO: Remove magic 1
  syncLoading.value = true;
  SyncRivalScoreLog(1).then(result => {
    if (result.Code != 200) {
      return Promise.reject();
    }
    syncLoading.value = false;
    notification.success({
      content: "同步成功",
      duration: 3000,
    });
  }).catch(err => {
    notification.error({
      content: "同步数据失败, 错误结果: " + err,
      duration: 3000,
      keepAliveOnHover: true,
    });
    syncLoading.value = false;
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
  width: 80%;
  white-space: normal;
}

.ps {
  height: 100%;
}
</style>
