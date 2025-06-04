<template>
  <n-spin :show="loading">
    <div class="container">
      <n-flex justify="space-between">
        <n-h1 prefix="bar" style="text-align: left;">
          <n-text type="primary">
            {{ t('title') }}
          </n-text>
        </n-h1>
        <n-button type="primary" @click="handleSaveSettings">{{ t('button.save') }}</n-button>
      </n-flex>

      <n-form ref="formRef" :model="model">
        <n-h2>
          {{ t('generalSettings.title') }}
          <n-p>
            <n-button type="info" @click="checkVersion">{{ t('button.checkVersion') }}</n-button>
          </n-p>
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
          <n-text>
            {{ t('folderSettings.title') }}
          </n-text>
          <!-- <n-p class="alert-p">
          {{ t('folderSettings.labelInternalServerPort') }}
        </n-p> -->
          <n-form-item :label="t('folderSettings.labelInternalServerPort')" path="internalServerPort">
            <n-input-number :show-button="false" v-model:value="model.InternalServerPort"
              :placeholder="t('folderSettings.placeholderInternalServerPort')" :maxlength="5" style="width: 150px;" />
          </n-form-item>
          <n-form-item path="folderSymbol">
            <template #label>
              <n-tooltip trigger="hover" style="max-width: 400px; text-align: left; white-space: pre-line;">
                <template #trigger>
                  <n-icon :component="HintIcon" />
                </template>
                {{ t('folderSettings.tipSymbol') }}
              </n-tooltip>
              {{ t('folderSettings.labelSymbol') }}
            </template>
            <n-input v-model:value="model.FolderSymbol" :placeholder="t('folderSettings.placeholderSymbol')"
              :maxlength="5" style="width: 200px;">
            </n-input>
          </n-form-item>
        </n-h2>
        <n-h2>
          <n-text>
            {{ t('downloadSettings.title') }}
          </n-text>
          <n-form-item :label="t('downloadSettings.labelDownloadSite')" path="downloadSite">
            <n-select v-model:value="model.DownloadSite" :options="downloadSiteOptions" style="width: 150px;" />
          </n-form-item>
          <n-form-item :label="t('downloadSettings.labelSeparateDownloadMD5')" path="separateDownloadMD5">
            <n-input v-model:value="model.SeparateDownloadMD5"
              :placeholder="t('downloadSettings.placeholderSeparateDownloadMD5')" style="width: 500px;" />
          </n-form-item>
          <n-form-item :label="t('downloadSettings.labelDownloadDirectory')" path="downloadDirectory">
            <n-button type="info" @click="chooseDownloadDirectory()">
              {{ t('button.chooseDownloadDirectory') }}
            </n-button>
            <n-divider vertical />
            <n-input v-model:value="model.DownloadDirectory"
              :placeholder="t('downloadSettings.placeholderDownloadDirectory')" style="width: 500px;" />
          </n-form-item>
          <n-form-item :label="t('downloadSettings.labelMaximumDownloadCount')" path="maximumDownloadCount">
            <n-input-number show-button v-model:value="model.MaximumDownloadCount" :min="1"
              style="width: 150px;" />
          </n-form-item>
        </n-h2>
      </n-form>
    </div>
  </n-spin>
</template>

<script setup lang="ts">
import { FormInst, SelectOption } from 'naive-ui';
import { ref } from 'vue';
import {
  ChatboxEllipsesOutline as HintIcon,
} from '@vicons/ionicons5';
import { QueryLatestVersion, ReadConfig, WriteConfig, OpenDirectoryDialog } from '@wailsjs/go/main/App';
import { config } from '../../../wailsjs/go/models';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const { t } = i18n;
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
  FolderSymbol: null,
  IgnoreVariantCourse: null,
  Locale: null,
  DownloadSite: null,
  SeparateDownloadMD5: null,
  DownloadDirectory: null,
  MaximumDownloadCount: null,
  EnableDownloadFeature: null // Unused
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

const downloadSiteOptions: Array<SelectOption> = [
  {
    label: "wriggle",
    value: "wriggle"
  },
];

function handleSaveSettings() {
  loading.value = true;
  WriteConfig(model.value)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      window.$notifySuccess(t('message.saveSuccess'));
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function chooseDownloadDirectory() {
  OpenDirectoryDialog("Choose Download Directory")
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Data != null && result.Data != undefined && result.Data != '') {
        model.value.DownloadDirectory = result.Data;
      }
    }).catch(err => window.$notifyError(err))
}

function loadSettings() {
  loading.value = true;
  ReadConfig()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      model.value = result.Data;
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function checkVersion() {
  QueryLatestVersion()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      window.$notifyInfo(result.Msg);
    }).catch(err => window.$notifyError(err));
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

:deep(.tooltip) {
  max-width: 400px;
  text-align: left;
  white-space: pre-line;
}
</style>

<i18n lang="json">{
  "en": {
    "title": "Settings",
    "button": {
      "save": "Save",
      "checkVersion": "Check Version",
      "chooseDownloadDirectory": "Choose Directory"
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
    "downloadSettings": {
      "title": "Download Settings",
      "labelDownloadSite": "Download Source Site",
      "labelSeparateDownloadMD5": "Download URL Pattern",
      "labelDownloadDirectory": "Download Directory",
      "labelMaximumDownloadCount": "Maximum Concurrent Download Count",
      "placeholderSeparateDownloadMD5": "Please input download url",
      "placeholderDownloadDirectory": "Please input download directory"
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
      "save": "保存",
      "checkVersion": "检查版本",
      "chooseDownloadDirectory": "选择文件夹"
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
    "downloadSettings": {
      "title": "下载设置",
      "labelDownloadSite": "下载源",
      "labelSeparateDownloadMD5": "下载地址",
      "labelDownloadDirectory": "下载路径",
      "labelMaximumDownloadCount": "最大同时下载数量",
      "placeholderSeparateDownloadMD5": "请输入下载地址",
      "placeholderDownloadDirectory": "请输入下载路径"
    },
    "options": {
      "yes": "是",
      "no": "否"
    },
    "message": {
      "saveSuccess": "保存成功"
    }
  }
}</i18n>
