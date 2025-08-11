<template>
  <n-spin :show="loading">
    <div class="container">
      <n-flex justify="space-between">
        <n-h1 prefix="bar" style="text-align: left;">
          <n-text type="primary">
            {{ t('title.settings') }}
          </n-text>
        </n-h1>
        <n-float-button shape="square" style="font-size: 15px; height: 39px; width: 60px;" :right="50" :top="20"
          type="primary" @click="handleSaveSettings">{{ t('button.save') }}</n-float-button>
        <!-- <n-button type="primary" @click="handleSaveSettings">{{ t('button.save') }}</n-button> -->
      </n-flex>

      <n-form ref="formRef" :model="model">
        <n-h2>
          {{ t('title.generalSetting') }}
          <n-p>
            <n-button type="info" @click="checkVersion">{{ t('button.checkVersion') }}</n-button>
          </n-p>
          <n-p>
            <n-form-item :label="t('form.labelLanguageReboot')">
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
                  {{ t('hint.ignoreVariantCourse') }}
                </n-tooltip>
                <span style="white-space: pre-line;">{{ t('form.labelIgnoreVariantCourse') }}</span>
              </template>
              <n-select v-model:value="model.IgnoreVariantCourse" :options="yesnoOptions" style="width: 150px;" />
            </n-form-item>
          </n-p>
          <n-p>
            <n-form-item :label="t('form.labelEnableAutoReload')">
              <n-select v-model:value="model.EnableAutoReload" :options="yesnoOptions" style="width: 150px;" />
            </n-form-item>
          </n-p>
          <n-p>
            <n-form-item :label="t('form.labelPreviewSite')" path="previewSite">
              <n-select v-model:value="model.PreviewSite" :options="previewSiteOptions" style="width: 150px;" />
            </n-form-item>
          </n-p>
          <n-p>
            <n-form-item :label="t('form.labelUseScoredatalog')">
              <n-select v-model:value="model.UseScoredatalog" :options="yesnoOptions" style="width: 150px;" />
            </n-form-item>
          </n-p>
        </n-h2>
        <n-h2>
          <n-text>
            {{ t('title.folderSetting') }}
          </n-text>
          <n-form-item :label="t('form.labelInternalServerPort')" path="internalServerPort">
            <n-input-number :show-button="false" v-model:value="model.InternalServerPort"
              :placeholder="t('form.placeholderInternalServerPort')" :maxlength="5" style="width: 150px;" />
          </n-form-item>
        </n-h2>
        <n-h2>
          <n-text>
            {{ t('title.downloadSetting') }}
          </n-text>
          <n-form-item :label="t('form.labelDownloadSite')" path="downloadSite">
            <n-select v-model:value="model.DownloadSite" :options="downloadSiteOptions" style="width: 150px;" />
          </n-form-item>
          <n-form-item :label="t('form.labelDownloadDirectory')" path="downloadDirectory">
            <n-button type="info" @click="chooseDownloadDirectory()">
              {{ t('button.chooseDirectory') }}
            </n-button>
            <n-divider vertical />
            <n-input v-model:value="model.DownloadDirectory" :placeholder="t('form.placeholderDownloadDirectory')"
              style="width: 500px;" />
          </n-form-item>
          <n-form-item :label="t('form.labelMaximumDownloadCount')" path="maximumDownloadCount">
            <n-input-number show-button v-model:value="model.MaximumDownloadCount" :min="1" style="width: 150px;" />
          </n-form-item>
        </n-h2>
        <n-h2>
          <n-text>{{ t('title.contactUs') }}</n-text>
          <n-p>
            <n-flex>
              <n-tooltip placement="top" trigger="hover">
                <template #trigger>
                  <n-button type="info" circle @click="gotoGithubRepo">
                    <template #icon>
                      <GithubIcon />
                    </template>
                  </n-button>
                </template>
                <span>{{ t('button.gotoGithub') }}</span>
              </n-tooltip>
              <n-tooltip placement="top" trigger="hover">
                <template #trigger>
                  <n-button type="info" circle @click="gotoDocument">
                    <template #icon>
                      <DocumentIcon />
                    </template>
                  </n-button>
                </template>
                <span>{{ t('button.gotoDocument') }}</span>
              </n-tooltip>
              <n-tooltip placement="top" trigger="hover">
                <template #trigger>
                  <n-button type="info" circle @click="gotoJoinQQGroup()">
                    <template #icon>
                      <QQIcon />
                    </template>
                  </n-button>
                </template>
                <span>{{ t('button.joinQQGroup') }}</span>
              </n-tooltip>
            </n-flex>
          </n-p>
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
  LogoGithub as GithubIcon,
  DocumentTextOutline as DocumentIcon,
} from '@vicons/ionicons5';
import { Qq as QQIcon } from "@vicons/fa";
import { QueryLatestVersion, ReadConfig, WriteConfig, OpenDirectoryDialog } from '@wailsjs/go/main/App';
import { config } from '../../../wailsjs/go/models';
import { useI18n } from 'vue-i18n';
import { BrowserOpenURL } from '@wailsjs/runtime/runtime';

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
  IgnoreVariantCourse: null,
  Locale: null,
  DownloadSite: null,
  DownloadDirectory: null,
  MaximumDownloadCount: null,
  EnableAutoReload: null,
  EnableDownloadFeature: null, // Unused
  PreviewSite: null,
  UseScoredatalog: null,
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
  {
    label: "konmai",
    value: "konmai"
  }
];

const previewSiteOptions: Array<SelectOption> = [
  {
    label: "sayaka",
    value: "sayaka",
  },
  {
    label: "konmai",
    value: "konmai"
  }
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

function gotoGithubRepo() {
  BrowserOpenURL("https://github.com/Catizard/lampghost");
}

function gotoDocument() {
  BrowserOpenURL("https://docs.qq.com/doc/DSEROYVZydXdDUUd6?rtkey=dbaaa020bbeabb89938632c4xV4Yu1")
}

function gotoJoinQQGroup() {
  BrowserOpenURL("https://qm.qq.com/cgi-bin/qm/qr?k=KS1jP1ii8AbSAtM89L8YRScn907QDSdP&jump_from=webapi&authKey=LIIF9+agKNK1cnoLdBu4lP8UqbOS+dDG/gfLFO2iwQIaa4WA0PNBUX6ERHm2GKNy");
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
