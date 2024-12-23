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
        <n-grid :cols="12">
          <n-gi :span="4">
            <n-dropdown trigger="hover" :options="playCountChartYearOptions" @select="handleSelect">
              <n-button>选择年份</n-button>
            </n-dropdown>
          </n-gi>
        </n-grid>
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
    <n-grid :cols="12" min-height="600px">
      <n-gi :span="4">
        这里应该有个难度表选择列表？
        <perfect-scrollbar style="height: 350px">
          <n-flex vertical>
            <n-button v-for="{ Name } in difftableHeaderList">
              {{ Name }}
            </n-button>
          </n-flex>
        </perfect-scrollbar>
      </n-gi>
      <n-gi :span="8">
        <vue-apex-charts height="450px" type="bar" :options="lampCountChartOptions"
          :series="lampCountChartOptions.series" />
      </n-gi>
    </n-grid>
  </perfect-scrollbar>
</template>

<script setup>
import VueApexCharts from 'vue3-apexcharts';
import { reactive, ref } from 'vue';
import {
  QueryUserInfoWithLevelLayeredDiffTableLampStatus,
  SyncRivalScoreLog,
  QueryUserPlayCountInYear,
  QueryMainUser
} from '../../wailsjs/go/controller/RivalInfoController'
import { FindDiffTableHeaderList } from '../../wailsjs/go/controller/DiffTableController';
import { useNotification } from 'naive-ui';
import dayjs from 'dayjs';
import router from '../router';

const LAMPS = [1, 4, 5, 11, 0];
const STR_LAMPS = ["FAILED", "Easy", "Normal", "Hard+", "NO_PLAY"];
const mainUser = ref(null);

const notification = useNotification();

const playCountChartOptions = ref({
  chart: {
    id: "chart-play-count",
    type: 'line',
    zoom: {
      enabled: false,
      allowMouseWheelZoom: false,
    }
  },
  xaxis: {
    categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
  },
  series: [
    {
      name: "游玩次数",
      data: [],
    }
  ],
});

const playCountChartYearOptions = ref([
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
])

const lampCountChartOptions = reactive({
  chart: {
    type: 'bar',
    stacked: true,
    stackType: "100%",
  },
  colors: [
    "#FF0000",
    "#00FF00",
    "#0000FF",
    "#00FFFF",
    "#FFFFFF"
  ],
  series: [],
  plotOptions: {
    bar: {
      horizontal: true,
      columnHeight: '100%',
    },
  },
  title: {
    text: '我觉得我是当前难度表的进度'
  },
  xaxis: {
    categories: [],
  },
  fill: {
    opacity: 1
  },
  legend: {
    position: 'top',
    horizontalAlign: 'left',
    offsetX: 40
  },
  dataLabels: {
    enabled: false,
  }
});

const playerData = reactive({
  playerName: "Catizard",
  playCount: 114514,
  lastUpdate: "",
});

const difftableHeaderList = ref([]);

function initUser() {
  QueryMainUser().then(result => {
    if (result.Code != 200) {
      // TODO: 看起来现在还没有初始化主用户？
      notification.error({
        content: "没有主用户，请先录入你自己的存档信息",
        duration: 3000,
      });
      router.push("/initialize")
      return Promise.reject()
    }
    mainUser.value = result.Data
    return Promise.resolve(result.Data);
  }).then(mainUser => {
    FindDiffTableHeaderList().then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      const { Rows } = result;
      difftableHeaderList.value = [...Rows];
      if (Rows.length > 0) {
        return Promise.resolve(Rows[0]);
      } else {
        return Promise.resolve(null);
      }
    }).then(result => {
      console.log("query difftable result:", result)
      if (result == null) {
        // TODO: 正确地显示无数据的情况
        return Promise.reject("目前无法处理一个难度表都没有的情况，请至少先加入一个数据")
      }
      QueryUserInfoWithLevelLayeredDiffTableLampStatus(mainUser.ID, result.ID).then(result => {
        if (result.Code != 200) {
          return Promise.reject(result.Msg)
        }
        const { Data } = result;
        console.log(Data)
        playerData.playerName = Data.Name
        playerData.playCount = Data.PlayCount
        playerData.lastUpdate = dayjs(Data.UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
        // Apply lamp status
        lampCountChartOptions.series = [];
        STR_LAMPS.forEach(lampName => {
          lampCountChartOptions.series.push({
            name: lampName,
            data: []
          })
        });
        lampCountChartOptions.xaxis.categories = [];
        const { DiffTableHeader } = Data;
        for (let i = 0; i < DiffTableHeader.SortedLevels.length; ++i) {
          const level = DiffTableHeader.SortedLevels[i];
          const levelName = DiffTableHeader.Symbol + level;
          lampCountChartOptions.xaxis.categories.push(levelName);
          const dataList = DiffTableHeader.LevelLayeredContents[level];
          for (let j = 0; j < LAMPS.length; ++j) {
            const lampValue = LAMPS[j];
            // TODO: wtf
            const count = dataList.filter(data => data.Lamp == lampValue || (data.Lamp >= 6 && lampValue == 11) || (lampValue == 1 && data.Lamp < 4 && data.Lamp > 0)).length;
            let v = {
              x: levelName,
              y: count
            }
            if (lampValue == 0) {
              v.fillColor = '#FFFFFF'
              v.strokeColor = '#FFFFFF'
            }
            lampCountChartOptions.series[j].data.push(v);
          }
        }

        console.log(lampCountChartOptions);
        return QueryUserPlayCountInYear(1, 2024)
      }).then(result => {
        if (result.Code != 200) {
          return Promise.reject(result.Msg)
        }
        const { Rows } = result;
        console.log(Rows);
        playCountChartOptions.value.series[0].data = [...Rows];
      }).catch(err => {
        console.log('error: ', err)
        notification.error({
          content: "获取用户数据失败: " + err,
          duration: 3000,
          keepAliveOnHover: true
        })
      })
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
      })
    })
  })
}

function initDiffTable() {

}

const syncLoading = ref(false);
function handleSyncClick() {
  syncLoading.value = true;
  SyncRivalScoreLog(mainUser.value.ID).then(result => {
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

function handleSelect() {
  notification.error({
    content: "噢，不好意思，这个功能还没实现",
    duration: 3000,
  });
}

initUser();
initDiffTable();
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
