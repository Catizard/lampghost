<template>
  <n-space justify="space-between" vertical style="width: 65%">
    <n-form ref="formRef" :rules="rules" :model="formData" :loading="loading">
      <n-form-item :label="t('form.labelLanguage')">
        <n-select v-model:value="locale" :options="localeOptions" style="width: 150px;" />
      </n-form-item>
      <n-form-item :label="t('form.labelUserName')" path="name">
        <n-input v-model:value="formData.name" :placeholder="t('form.placeholderUserName')" />
      </n-form-item>
      <n-radio-group v-model:value="importStrategy" name="radiogroup">
        <n-radio-button key="directory" value="directory" :label="t('form.labelBeatorajaDirectory')" />
        <n-radio-button key="separate" value="separate" :label="t('form.labelSeparateFiles')" />
        <n-radio-button key="LR2" value="LR2" :label="t('form.labelLR2')" />
      </n-radio-group>
      <template v-if="importStrategy == 'directory'">
        <div style="margin-top: 10px;">
          <n-button type="info" @click="chooseBeatorajaDirectory()">
            {{ t('button.chooseBeatorajaDirectory') }}
          </n-button>
          <n-select :options="playerDirectories" v-model:value="formData.playerDirectory"
            style="margin-top: 10px;width: 250px;">
          </n-select>
        </div>
      </template>
      <template v-if="importStrategy == 'separate'">
        <n-form-item :label="t('form.labelSongdataPath')" path="songdataPath" style="margin-top: 10px;">
          <n-button type="info" @click="chooseFile('Choose songdata.db', 'songdataPath')">
            {{ t('button.chooseFile') }}
          </n-button>
          <n-divider vertical />
          <n-input v-model:value="formData.songdataPath" :placeholder="t('form.placeholderSongdataPath')" />
        </n-form-item>
        <n-form-item :label="t('form.labelScorelogPath')" path="scorelogPath">
          <n-button type="info" @click="chooseFile('Choose scorelog.db', 'scorelogPath')">
            {{ t('button.chooseFile') }}
          </n-button>
          <n-divider vertical />
          <n-input v-model:value="formData.scorelogPath" :placeholder="t('form.placeholderScorelogPath')" />
        </n-form-item>
        <n-form-item :label="t('form.labelScoredatalogPath')" path="scoredatalogPath">
          <n-button type="info" @click="chooseFile('Choose scoredatalog.db', 'scoredatalogPath')">
            {{ t('button.chooseFile') }}
          </n-button>
          <n-divider vertical />
          <n-input v-model:value="formData.scoredatalogPath" :placeholder="t('form.placeholderScoredatalogPath')" />
        </n-form-item>
        <n-form-item :label="t('form.labelScoredataPath')" path="scoredataPath">
          <n-button type="info" @click="chooseFile('Choose score.db', 'scoredataPath')">
            {{ t('button.chooseFile') }}
          </n-button>
          <n-divider vertical />
          <n-input v-model:value="formData.scoredataPath" :placeholder="t('form.placeholderScoredataPath')" />
        </n-form-item>
      </template>
      <template v-if="importStrategy == 'LR2'">
        <n-flex justify="space-between">
          <n-form-item :label="t('form.labelScanBMSFiles')" style="margin-top: 10px;">
            <DirectoryTable v-model:directories="formData.BMSDirectories" />
          </n-form-item>
          <n-button type="primary" @click="chooseBMSDirectory">
            {{ t('button.chooseDirectory') }}
          </n-button>
        </n-flex>
        <n-form-item :label="t('form.labelUserDBPath')" path="scorelogPath">
          <n-button type="info" @click="chooseFile('Choose user database file', 'scorelogPath')">
            {{ t('button.chooseFile') }}
          </n-button>
          <n-divider vertical />
          <n-input v-model:value="formData.scorelogPath" :placeholder="t('form.placeholderUserDBPath')" />
        </n-form-item>
      </template>
      <n-form-item>
        <n-button :loading="loading" attr-type="button" @click="handleSubmit" type="primary">
          {{ t('button.submit') }}
        </n-button>
      </n-form-item>
    </n-form>
  </n-space>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { ref, watch, reactive, Ref } from "vue";
import { ChooseBeatorajaDirectory, InitializeMainUser, OpenDirectoryDialog, OpenFileDialog, QueryMainUser } from "@wailsjs/go/main/App";
import { SelectOption } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import DirectoryTable from './directoryTable.vue';
import { useUserStore } from '@/stores/user';

const { t } = useI18n();
const globalI18n = useI18n({ useScope: 'global' });

const props = defineProps<{
  moveOn: () => void
}>();

const importStrategy: Ref<"directory" | "separate" | "LR2"> = ref("directory");

// target == "scorelogPath" | "songdataPath" | "scoredatalogPath"
function chooseFile(title: string, target: "scorelogPath" | "songdataPath" | "scoredatalogPath" | "scoredataPath") {
  OpenFileDialog(title)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Data != null && result.Data != undefined && result.Data != "") {
        if (target == "scorelogPath") {
          formData.scorelogPath = result.Data;
        } else if (target == "songdataPath") {
          formData.songdataPath = result.Data;
        } else if (target == "scoredatalogPath") {
          formData.scoredatalogPath = result.Data;
        } else if (target == "scoredataPath") {
          formData.scoredataPath = result.Data;
        }
      }
    }).catch(err => window.$notifyError(err));
}

const playerDirectories: Ref<SelectOption[]> = ref([]);
function chooseBeatorajaDirectory() {
  OpenDirectoryDialog(t('title.chooseBeatorajaDirectory')).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    return Promise.resolve(result.Data);
  }).then(path => {
    ChooseBeatorajaDirectory(path)
      .then(result => {
        if (result.Code != 200) {
          return Promise.reject(result.Msg)
        }
        const data: dto.BeatorajaDirectoryMeta = result.Data;
        playerDirectories.value = data.PlayerDirectories.map(row => {
          return {
            label: row,
            value: row
          } as SelectOption
        });
        formData.beatorajaDirectoryPath = data.BeatorajaDirectoryPath;
        formData.playerDirectory = playerDirectories.value[0].value as string;
      })
  }).catch(err => window.$notifyError(err));
}

function chooseBMSDirectory() {
  OpenDirectoryDialog(t('title.chooseBMSDirectory')).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    const path = result.Data;
    if (path != "") {
      formData.BMSDirectories.push(path);
    }
  }).catch(err => window.$notifyError(err))
}

const loading = ref(false);
const userStore = useUserStore();
const formRef = ref(null);
const formData = reactive({
  locale: null,
  songdataPath: "",
  scorelogPath: "",
  scoredatalogPath: "",
  scoredataPath: "",
  name: "",
  beatorajaDirectoryPath: "",
  playerDirectory: "",
  BMSDirectories: [],
});

const localeOptions = [
  {
    label: "English",
    value: "en",
  },
  {
    label: "中文",
    value: "zh-CN"
  }
];
// I have no idea why need this hack
const locale = globalI18n.locale;
watch(locale, (newLocale) => {
  (globalI18n as any).locale = ref(newLocale);
})

const rules = {
  name: {
    required: true,
    message: t('rule.missingUserName'),
    trigger: ["input", "blur"],
  },
  songdataPath: {
    required: true,
    message: t('rule.missingSongdataPath'),
    trigger: ["input", "blur"],
  },
  scorelogPath: {
    required: true,
    message: t('rule.missingScorelogPath'),
    trigger: ["input", "blur"],
  },
  scoredatalogPath: {
    required: true,
    message: t('rule.missingScoredatalogPath'),
    trigger: ["input", "blur"],
  },
  scoredataPath: {
    required: true,
    message: t('rule.missingScoredataPath'),
    trigger: ["input", "blur"],
  }
};

function handleSubmit(e) {
  e.preventDefault();
  formRef.value?.validate(async (errors) => {
    if (errors) {
      console.error(errors);
    } else {
      loading.value = true;
      // TODO: This was a mistake...
      const rivalInfo = {
        Name: formData.name,
        ScoreLogPath: formData.scorelogPath,
        SongDataPath: formData.songdataPath,
        ScoreDataLogPath: formData.scoredatalogPath,
        ScoreDataPath: formData.scoredataPath,
        BeatorajaDirectoryPath: formData.beatorajaDirectoryPath,
        PlayerDirectory: formData.playerDirectory,
        ImportStrategy: importStrategy.value,
        BMSDirectories: [],
      };
      if (importStrategy.value == "LR2") {
        rivalInfo.ScoreDataLogPath = rivalInfo.ScoreLogPath;
        rivalInfo.ScoreDataPath = rivalInfo.ScoreLogPath;
        formData.BMSDirectories.forEach(p => rivalInfo.BMSDirectories.push(p));
      }
      console.log('param: ', rivalInfo);
      try {
        const result = await InitializeMainUser(rivalInfo as any);
        if (result.Code != 200) {
          throw result.Msg;
        }
        const mainUserResult = await QueryMainUser();
        if (mainUserResult.Code != 200) {
          throw result.Msg;
        }
        userStore.setter(mainUserResult.Data);
        props.moveOn();
      } catch (err) {
        window.$notifyError(err);
      } finally {
        loading.value = false;
      }
    }
  });
}
</script>
