<template>
  <div class="container">
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: left;">
        <n-text type="primary">
          Settings
        </n-text>
      </n-h1>
      <n-button type="primary" @click="handleSaveSettings" :loading="loading">保存</n-button>
    </n-flex>

    <n-form ref="formRef" :model="model">
      <n-h2>
        通用设置
        <n-p>
          <n-form-item label="语言设置">
            <n-select v-model:value="model.Locale" :options="localeOptions" style="width: 150px;" />
          </n-form-item>
        </n-p>
        <n-p>
          <n-form-item>
            <template #label>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-icon :component="HintIcon" />
                </template>
                当开启时, 移除所有no_great/no_speed/no_good等变化出来的段位,仅保留原版段位。
                <br>
                注意调整该选项并不会重新生成rival的tag标记,你需要手动重新生成。
              </n-tooltip>
              是否隐藏变化段位
            </template>
            <n-select v-model:value="model.IgnoreVariantCourse" :options="yesnoOptions" style="width: 150px;" />
          </n-form-item>
        </n-p>
      </n-h2>
      <n-h2>
        存档设置
        <n-p>
          <n-form-item label="用户名, 如果你想换个用户名的话" path="userName">
            <n-input show-count v-model:value="model.UserName" placeholder="请输入用户名" style="width: 50%;"
              :loading="loading" />
          </n-form-item>
        </n-p>
        <n-p class="alert-p">
          注意: 如果你修改了文件路径地址, 保存设置时会重新再次加载你的存档数据, 因此你不需要再手动同步一次存档
        </n-p>
        <n-p>
          <n-form-item label="scorelog.db文件路径" path="scorelogFilePath">
            <n-input clearable v-model:value="model.ScorelogFilePath" placeholder="请输入scorelog.db文件路径"
              style="width: 50%;" :loading="loading" />
          </n-form-item>
        </n-p>
        <n-p>
          <n-form-item label="songdata.db文件路径" path="songdataFilePath">
            <n-input clearable v-model:value="model.SongdataFilePath" placeholder="请输入songdata.db文件路径"
              style="width: 50%;" :loading="loading" />
          </n-form-item>
        </n-p>
        <n-p>
          <n-form-item path="scoreFilePath">
            <template #label>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-icon :component="HintIcon" />
                </template>
                添加该路径之后可以显示玩家当前的游玩时间等信息, 可加可不加
              </n-tooltip>
              score.db文件路径
            </template>
            <n-input clearable v-model:value="model.ScoreFilePath" placeholder="请输入score.db文件路径" style="width: 50%;"
              :loading="loading" />
          </n-form-item>
        </n-p>
      </n-h2>
      <n-h2>
        <n-text>
          自定义收藏夹
        </n-text>
        <n-p class="alert-p">
          注意: 如果你修改了下列设置之前已经导入过难度表, 需要修改表定义之后重新刷新所有难度表信息来移出此前导入的表信息！
        </n-p>
        <n-form-item label="内部服务器端口号" path="internalServerPort">
          <n-input-number :show-button="false" v-model:value="model.InternalServerPort" placeholder="请输入端口号"
            :maxlength="5" style="width: 150px;" :loading="loading" />
        </n-form-item>
        <n-form-item path="folderSymbol">
          <template #label>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-icon :component="HintIcon" />
              </template>
              难度表标志即发狂表的🌟或satellite的sl, 建议留空。
            </n-tooltip>
            难度表标志
          </template>
          <n-input v-model:value="model.FolderSymbol" placeholder="默认为空" :maxlength="5" style="width: 150px;"
            :loading="loading">
          </n-input>
        </n-form-item>
      </n-h2>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { FormInst, SelectOption, useNotification } from 'naive-ui';
import { ref } from 'vue';
import {
  ChatboxEllipsesOutline as HintIcon,
} from '@vicons/ionicons5';
import { ReadConfig, WriteConfig } from '@wailsjs/go/controller/ConfigController';
import { config } from '../../../wailsjs/go/models';

const notification = useNotification();
const localeOptions: Array<SelectOption> = [
  {
    label: "English",
    value: "en",
  },
  {
    label: "中文",
    value: "zh-CN"
  }
]
const formRef = ref<FormInst | null>(null);
const model = ref<config.ApplicationConfig>({
  InternalServerPort: null,
  UserName: null,
  ScorelogFilePath: null,
  SongdataFilePath: null,
  ScoreFilePath: null,
  FolderSymbol: null,
  IgnoreVariantCourse: null,
  Locale: null,
});
const loading = ref(false);
const yesnoOptions: Array<SelectOption> = [
  {
    label: "是",
    value: 1,
  },
  {
    label: "否",
    value: 0,
  }
];

function handleSaveSettings() {
  loading.value = true;
  WriteConfig(model.value)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      notification.success({
        content: "保存成功",
        duration: 3000,
        keepAliveOnHover: false
      })
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true
      })
    }).finally(() => {
      loading.value = false;
    })
}

function loadSettings() {
  loading.value = true;
  ReadConfig()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      model.value = result.Data;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true
      })
    }).finally(() => {
      loading.value = false;
    })
}

loadSettings();
</script>

<style scoped>
.container {
  text-align: left;
}

.alert-p {
  color: #ff1100;
}
</style>