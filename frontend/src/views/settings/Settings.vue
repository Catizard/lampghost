<template>
  <div class="container">
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: left;">
        <n-text type="primary">
          {{ t('title') }} 
        </n-text>
      </n-h1>
      <n-button type="primary" @click="handleSaveSettings" :loading="loading">{{ t('button.save') }}</n-button>
    </n-flex>

    <n-form ref="formRef" :model="model">
      <n-h2>
        {{ t('generalSettings.title') }}
        <n-p>
          <n-form-item :label="t('generalSettings.labelLanguage')">
            <n-select v-model:value="model.Locale" :options="localeOptions" style="width: 150px;" />
          </n-form-item>
        </n-p>
        <n-p>
          <n-form-item>
            <template #label>
              <n-tooltip trigger="hover" style="max-width: 400px; text-align: left; white-space: pre-line;">
                <template #trigger>
                  <n-icon :component="HintIcon" />
                </template>
                {{ t('generalSettings.tipIgnoreVariantCourse') }}
              </n-tooltip>
              <span style="white-space: pre-line;">{{ t('generalSettings.labelIgnoreVariantCourse') }}</span>
            </template>
            <n-select v-model:value="model.IgnoreVariantCourse" :options="yesnoOptions" style="width: 150px;" />
          </n-form-item>
        </n-p>
      </n-h2>
      <n-h2>
        {{ t('saveSettings.title') }}
        <n-p>
          <n-form-item :label="t('saveSettings.labelUserName')" path="userName">
            <n-input show-count v-model:value="model.UserName" :placeholder="t('saveSettings.placeholderUserName')" style="width: 50%;"
              :loading="loading" />
          </n-form-item>
        </n-p>
        <n-p class="alert-p">
          {{ t('saveSettings.alert') }}
        </n-p>
        <n-p>
          <n-form-item :label="t('saveSettings.labelScorelogPath')" path="scorelogFilePath">
            <n-input clearable v-model:value="model.ScorelogFilePath" :placeholder="t('saveSettings.placeholderScorelogPath')"
              style="width: 50%;" :loading="loading" />
          </n-form-item>
        </n-p>
        <n-p>
          <n-form-item :label="t('saveSettings.labelSongdataPath')" path="songdataFilePath">
            <n-input clearable v-model:value="model.SongdataFilePath" :placeholder="t('saveSettings.placeholderSongdataPath')"
              style="width: 50%;" :loading="loading" />
          </n-form-item>
        </n-p>
        <!-- <n-p>
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
            <n-input disabled clearable v-model:value="model.ScoreFilePath" placeholder="请输入score.db文件路径"
              style="width: 50%;" :loading="loading" />
          </n-form-item>
        </n-p> -->
      </n-h2>
      <n-h2>
        <n-text>
          {{ t('folderSettings.title') }} 
        </n-text>
        <!-- <n-p class="alert-p">
          {{ t('folderSettings.labelInternalServerPort') }} 
        </n-p> -->
        <n-form-item :label="t('folderSettings.labelInternalServerPort')" path="internalServerPort">
          <n-input-number :show-button="false" v-model:value="model.InternalServerPort" :placeholder="t('folderSettings.placeholderInternalServerPort')"
            :maxlength="5" style="width: 150px;" :loading="loading" />
        </n-form-item>
        <n-form-item path="folderSymbol">
          <template #label>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-icon :component="HintIcon" />
              </template>
              {{ t('folderSettings.tipSymbol') }}
            </n-tooltip>
            {{ t('folderSettings.labelSymbol') }}
          </template>
          <n-input v-model:value="model.FolderSymbol" :placeholder="t('folderSettings.placeholderSymbol')" :maxlength="5" style="width: 200px;"
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
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const { t } = i18n;
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
    label: t('options.yes'),
    value: 1,
  },
  {
    label: t('options.no'),
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
        content: t('message.saveSuccess'),
        duration: 3000,
        keepAliveOnHover: false
      });
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

<i18n lang="json">
{
  "en": {
    "title": "Settings",
    "button": {
      "save": "Save"
    },
    "generalSettings": {
      "title": "General Settings",
      "labelLanguage": "Language Options(Reboot needed)",
      "labelIgnoreVariantCourse": "Controls whether we should ignore variant courses or not",
      "tipIgnoreVariantCourse": "When opend, ignore all the courses which constraint contains 'no_good', 'no_speed', 'no_great', only reserve the normal course\nWarning: rival's tag won't be re-generated again by changing this option, you still need to reload manually"
    },
    "saveSettings": {
      "title": "Save File Settings",
      "alert": "If you changed the file path, by pressing the save button your save file would be reloaded automatically, there's no need to reload again manually",
      "labelUserName": "Your user name",
			"labelSongdataPath": "songdata.db file path",
			"labelScorelogPath": "scorelog.db file path",
			"placeholderUserName": "Please input your name",
			"placeholderSongdataPath": "Please input songdata.db file path",
			"placeholderScorelogPath": "Please input scorelog.db file path"
    },
    "folderSettings": {
      "title": "Custom Folder Settings",
      "alert": "If you have imported some difficult tables before changing below options, you need to reload all tables to clean up the old folder content",
      "labelInternalServerPort": "Internal Server Port",
      "placeholderInternalServerPort": "Please input internal server port",
      "labelSymbol": "Table symbol",
      "tipSymbol": "Equals to the sl/st from Satellite, recommend to leave empty",
      "placeholderSymbol": "Defaults to empty"
    },
    "options": {
      "yes": "Yes",
      "no": "No"
    },
    "message": {
      "saveSuccess": "Saved successfully"
    }
  },
  "zh-CN": {
    "title": "设置",
    "button": {
      "save": "保存"
    },
    "generalSettings": {
      "title": "通用设置",
      "labelLanguage": "语言设置(需要重启)",
      "labelIgnoreVariantCourse": "是否忽略带有特殊变化的段位",
      "tipIgnoreVariantCourse": "开启时移除所有带有no_good/no_great/no_speed的段位\n注意: 修改该配置不会自动更新用户的标签,你必须手动刷新用户的标签信息"
    },
    "saveSettings": {
      "title": "存档设置",
      "alert": "如果你修改了这里的设置,保存时会自动重新加载存档信息, 不需要手动加载数据",
      "labelUserName": "用户名称",
			"labelSongdataPath": "songdata.db文件路径",
			"labelScorelogPath": "scorelog.db文件路径",
			"placeholderUserName": "请输入用户名",
			"placeholderSongdataPath": "请输入songdata.db文件路径",
			"placeholderScorelogPath": "请输入scorelog.db文件路径"
    },
    "folderSettings": {
      "title": "收藏夹设置",
      "labelInternalServerPort": "内部服务器端口号",
      "placeholderInternalServerPort": "请输入内部服务器端口号",
      "labelSymbol": "难度表标志",
      "tipSymbol": "标志即satellite表的sl/st等, 如果你不知道这是什么, 建议留空",
      "placeholderSymbol": "默认为空"
    },
    "options": {
      "yes": "是",
      "no": "否"
    },
    "message": {
      "saveSuccess": "保存成功"
    }
  }
}
</i18n>