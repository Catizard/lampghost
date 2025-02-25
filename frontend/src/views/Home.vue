<template>
  <n-h1 prefix="bar" style="text-align: start;">
    <n-text type="primary">
      {{ t('infoTitle') }}
    </n-text>
  </n-h1>
  <n-grid :cols="12">
    <n-gi :span="4">
      <n-flex>
        {{ t('playerInfo.name') }}: {{ playerData.playerName }}
        <n-divider />
        {{ t('playerInfo.count') }}: {{ playerData.playCount }}
        <n-divider />
        {{ t('playerInfo.lastSyncTime') }}: {{ playerData.lastUpdate }}
        <n-divider />
        <n-button @click="handleSyncClick" :loading="syncLoading">
          {{ t('button.sync') }}
        </n-button>
      </n-flex>
    </n-gi>
    <n-gi :span="8">
      <n-grid :cols="12">
        <n-gi :span="4">
          <n-dropdown trigger="hover" :options="playCountChartYearOptions" @select="handleSelect">
            <n-button>{{ t('button.chooseYear') }}</n-button>
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
      {{ t('lampStatusTitle') }}
    </n-text>
  </n-h1>
  <n-grid :cols="12" min-height="600px">
    <n-gi :span="4">
      <n-dropdown trigger="hover" :options="difftableStatusOptions" @select="handleChangeDifftableStatusOption">
        <n-button :loading="difftableChangingLoading">{{ currentChoosingDifftableName }}</n-button>
      </n-dropdown>
    </n-gi>
  </n-grid>
  <vue-apex-charts height="450px" type="bar" :options="lampCountChartOptions" :series="lampCountChartOptions.series" />
</template>

<script setup>
import VueApexCharts from 'vue3-apexcharts';
import { computed, reactive, ref } from 'vue';
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
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const { t } = i18n;

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
      name: t('playerInfo.count'),
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
  playerName: "U",
  playCount: 0,
  lastUpdate: "",
});

const difftableHeaderList = ref([]);
const difftableStatusOptions = ref([]);
const currentChoosingDifftableName = ref("");

function initUser() {
  QueryMainUser().then(result => {
    if (result.Code != 200) {
      notification.error({
        content: t('message.noMainUserError'),
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
      difftableStatusOptions.value = difftableHeaderList.value.map(header => {
        return {
          label: header.Name,
          key: header.ID,
        }
      })
      if (Rows.length > 0) {
        return Promise.resolve(Rows[0]);
      } else {
        return Promise.resolve(null);
      }
    }).then(result => {
      if (result == null) {
        // TODO: 正确地显示无数据的情况
        return Promise.reject(t('message.noTableError'))
      }
      QueryUserInfoWithLevelLayeredDiffTableLampStatus(mainUser.ID, result.ID).then(result => {
        if (result.Code != 200) {
          return Promise.reject(result.Msg)
        }
        const { Data } = result;
        applyPlayerData(Data);
        const { DiffTableHeader } = Data;
        applyDifftableStatus(DiffTableHeader);
        return QueryUserPlayCountInYear(1, 2024)
      }).then(result => {
        if (result.Code != 200) {
          return Promise.reject(result.Msg)
        }
        const { Rows } = result;
        playCountChartOptions.value.series[0].data = [...Rows];
      }).catch(err => {
        notification.error({
          content: t('message.loadUserDataErrorPrefix') + err,
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

function applyPlayerData(data) {
  playerData.playerName = data.Name;
  playerData.playerName = data.Name
  playerData.playCount = data.PlayCount
  playerData.lastUpdate = dayjs(data.UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
}

function applyDifftableStatus(header) {
  // (1) Set choosed option
  currentChoosingDifftableName.value = header.Name;
  // (2) Lamp status chart
  lampCountChartOptions.series = [];
  STR_LAMPS.forEach(lampName => {
    lampCountChartOptions.series.push({
      name: lampName,
      data: []
    })
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
}

const difftableChangingLoading = ref(false);
function handleChangeDifftableStatusOption(key) {
  difftableChangingLoading.value = true;
  QueryUserInfoWithLevelLayeredDiffTableLampStatus(mainUser.value.ID, key).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    console.log(result);
    let { DiffTableHeader } = result.Data;
    console.log(DiffTableHeader);
    applyDifftableStatus(DiffTableHeader);
  }).catch(err => {
    notification.error({
      content: err,
      duration: 3000,
      keepAliveOnHover: true
    })
  }).finally(() => {
    difftableChangingLoading.value = false;
  })
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
      content: t('message.reloadSuccess'),
      duration: 3000,
    });
  }).catch(err => {
    notification.error({
      content: t('message.reloadFailedPrefix') + err,
      duration: 3000,
      keepAliveOnHover: true,
    });
    syncLoading.value = false;
  })
}

function handleSelect() {
  notification.error({
    content: t('message.unfinishedFeature'),
    duration: 3000,
  });
}

initUser();
</script>

<i18n lang="json">{
  "en": {
    "infoTitle": "Player Info",
    "playerInfo": {
      "name": "Player Name",
      "count": "Player Count",
      "lastSyncTime": "Last Sync Time",
    },
    "button": {
      "sync": "Reload Save File",
      "chooseYear": "Choose Year",
    },
    "lampStatusTitle": "Lamp Status",
    "message": {
      "noMainUserError": "Found no main user, please first load your save file in",
      "noTableError": "Cannot handle no difficult table data currenlty, please add at least one table first",
      "loadUserDataErrorPrefix": "Cannot load user data: ",
      "reloadSuccess": "Successfully reloaded",
      "reloadFailedPrefix": "Failed to load save file, error message: ",
      "unfinishedFeature": "Sorry, haven't implemented yet"
    }
  },
  "zh-CN": {
    "infoTitle": "玩家信息",
    "playerInfo": {
      "name": "玩家名称",
      "count": "游玩次数",
      "lastSyncTime": "最后同步时间",
    },
    "button": {
      "sync": "同步最新存档",
      "chooseYear": "选择年份",
    },
    "lampStatusTitle": "点灯情况",
    "message": {
      "noMainUserError": "找不到主用户信息，请先导入你自己的存档",
      "noTableError": "目前无法处理一个难度表都没有的情况，请至少先添加一个难度表",
      "loadUserDataErrorPrefix": "获取用户信息失败: ",
      "reloadSuccess": "同步成功",
      "reloadFailedPrefix": "同步失败，返回结果: ",
      "unfinishedFeature": "不好意思，这个功能还没实现"
    }
  },
}</i18n>

<style scoped>
.n-button {
  width: 80%;
  white-space: normal;
}
</style>
