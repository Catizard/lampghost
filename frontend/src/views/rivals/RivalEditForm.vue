<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('modal.labelRivalName')" path="Name">
        <n-input v-model:value="formData.Name" :placeholder="t('modal.placeholderRivalName')" />
      </n-form-item>
      <n-form-item :label="t('modal.labelScoreLogPath')" path="ScoreLogPath">
        <n-button :disabled="!formData.MainUser" type="info" @click="chooseFile('Choose scorelog.db', 'scorelogPath')">
          {{ t('button.chooseFile') }}
        </n-button>
        <n-divider vertical />
        <n-input v-model:value="formData.ScoreLogPath" :placeholder="t('modal.placeholderScoreLogPath')" />
      </n-form-item>
      <n-form-item :label="t('modal.labelSongDataPath')" path="SongDataPath">
        <n-button type="info" @click="chooseFile('Choose songdata.db', 'songdataPath')">
          {{ t('button.chooseFile') }}
        </n-button>
        <n-divider vertical />
        <n-input :disabled="!formData.MainUser" v-model:value="formData.SongDataPath"
          :placeholder="t('modal.placeholderSongDataPath')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import { QueryUserInfoByID, UpdateRivalInfo } from '@wailsjs/go/main/App';
import { OpenFileDialog } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { FormInst, useNotification } from 'naive-ui';
import { reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const notification = useNotification();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose({ open });

const show = ref(false);

function open(rivalID: number) {
  if (rivalID == null || rivalID == 0) {
    notification.error({
      content: t('message.noChosenRivalError'),
      duration: 3000,
      keepAliveOnHover: true
    });
    show.value = false;
    return
  }

  formData.value.ID = rivalID;
  show.value = true;
  loading.value = true;

  QueryUserInfoByID(rivalID)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      const data: dto.RivalInfoDto = result.Data;
      formData.value.Name = data.Name;
      formData.value.ScoreLogPath = data.ScoreLogPath;
      formData.value.SongDataPath = data.SongDataPath;
      formData.value.MainUser = data.MainUser;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      });
      show.value = false;
    }).finally(() => {
      loading.value = false;
    });
}

const loading = ref(false);

const formRef = ref<FormInst | null>(null);
const formData = ref({
  ID: 0,
  Name: null,
  ScoreLogPath: null,
  SongDataPath: null,
  MainUser: false,
});
const rules = reactive({
  Name: {
    required: true,
    message: t('rules.missingRivalName'),
    trigger: ["input", "blur"],
  },
  ScoreLogPath: {
    required: true,
    message: t('rules.missingScoreLogPath'),
    trigger: ["input", "blur"],
  },
  SongDataPath: {
    required: false,
    message: t('rules.missingSongDataPath'),
    trigger: ["input", "blur"],
  }
});

function handlePositiveClick(): boolean {
  loading.value = true;
  formRef.value
    ?.validate()
    .then(() => {
      return UpdateRivalInfo(formData.value as any)
        .then(result => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          show.value = false;
          emit('refresh')
        });
    })
    .catch((err) => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true
      });
    }).finally(() => {
      loading.value = false;
    });
  return false;
}

function handleNegativeClick() {
  formData.value.ID = 0;
  formData.value.Name = null;
  formData.value.ScoreLogPath = null;
  formData.value.SongDataPath = null;
  formData.value.MainUser = false;
}

// target == "scorelogPath" | "songdataPath"
function chooseFile(title, target) {
  OpenFileDialog(title)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Data != null && result.Data != undefined && result.Data != "") {
        if (target == "scorelogPath") {
          formData.value.ScoreLogPath = result.Data;
        } else if (target == "songdataPath") {
          formData.value.SongDataPath = result.Data;
        }
      }
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true
      })
    });
}

watch(() => formData.value.MainUser, newValue => {
  rules.SongDataPath.required = newValue;
})
</script>

<i18n lang="json">{
  "en": {
    "modal": {
      "title": "Edit Rival Info",
      "positiveText": "Submit",
      "negativeText": "Cancel",
      "labelRivalName": "Name",
      "labelScoreLogPath": "scorelog.db file path",
      "labelSongDataPath": "songdata.db file path",
      "placeholderRivalName": "Please input rival's name",
      "placeholderScoreLogPath": "Please input scorelog.db file path",
      "placeholderSongDataPath": "Please input songdata.db file path"
    },
    "button": {
      "chooseFile": "Choose File"
    },
    "rules": {
      "missingRivalName": "Rival's name cannot be empty",
      "missingScoreLogPath": "scorelog.db file path cannot be empty",
      "missingSongDataPath": "songdata.db file path cannot be empty"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "修改用户信息",
      "positiveText": "提交",
      "negativeText": "取消",
      "labelRivalName": "名称",
      "labelScoreLogPath": "scorelog.db文件路径",
      "labelSongDataPath": "songdata.db文件路径",
      "placeholderRivalName": "请输入好友名称",
      "placeholderScoreLogPath": "请输入scorelog.db文件路径",
      "placeholderSongDataPath": "请输入songdata.db文件路径"
    },
    "button": {
      "chooseFile": "选择文件"
    },
    "rules": {
      "missingRivalName": "好友名称不可为空",
      "missingScoreLogPath": "scorelog.db文件路径不可为空",
      "missingSongDataPath": "songdata.db文件路径不可为空"
    }
  }
}</i18n>