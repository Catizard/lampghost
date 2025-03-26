<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('modal.labelRivalName')" path="Name">
        <n-input v-model:value="formData.Name" :placeholder="t('modal.placeholderRivalName')" />
      </n-form-item>
      <n-form-item :label="t('modal.labelScoreLogPath')" path="ScoreLogPath">
        <n-input v-model:value="formData.ScoreLogPath" :placeholder="t('modal.placeholderScoreLogPath')" />
      </n-form-item>
      <n-form-item :label="t('modal.labelSongDataPath')" path="SongDataPath">
        <n-input disabled v-model:value="formData.SongDataPath" :placeholder="t('modal.placeholderSongDataPath')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import { AddRivalInfo } from '@wailsjs/go/controller/RivalInfoController';
import { entity } from '@wailsjs/go/models';
import { FormInst, useNotification } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const show = defineModel<boolean>("show");
const emit = defineEmits<{
  (e: 'refresh'): void
}>();

const { t } = useI18n();
const notification = useNotification();
const loading = ref(false);

const formRef = ref<FormInst | null>(null);
const formData = ref({
  Name: null,
  ScoreLogPath: null,
  SongDataPath: null,
});
const rules = {
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
  // SongDataPath: {
  // 	required: true,
  // 	message: t('rules.missingSongDataPath'),
  // 	trigger: ["input", "blur"],
  // }
};

function handlePositiveClick(): boolean {
  loading.value = true;
  formRef.value
    ?.validate()
    .then(() => {
      return AddRivalInfo(formData.value as any as entity.RivalInfo)
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
    }).finally(() => loading.value = false);
  return false;
}

function handleNegativeClick() {
  formData.value.Name = null;
  formData.value.ScoreLogPath = null;
  formData.value.SongDataPath = null;
}
</script>

<i18n lang="json">{
  "en": {
    "modal": {
      "title": "Add Rival",
      "positiveText": "Submit",
      "negativeText": "Cancel",
      "labelRivalName": "Name",
      "labelScoreLogPath": "scorelog.db file path",
      "labelSongDataPath": "songdata.db file path",
      "placeholderRivalName": "Please input rival's name",
      "placeholderScoreLogPath": "Please input scorelog.db file path",
      "placeholderSongDataPath": "Please input songdata.db file path"
    },
    "rules": {
      "missingRivalName": "Rival's name cannot be empty",
      "missingScoreLogPath": "scorelog.db file path cannot be empty",
      "missingSongDataPath": "songdata.db file path cannot be empty"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "新增好友",
      "positiveText": "提交",
      "negativeText": "取消",
      "labelRivalName": "名称",
      "labelScoreLogPath": "scorelog.db文件路径",
      "labelSongDataPath": "songdata.db文件路径",
      "placeholderRivalName": "请输入好友名称",
      "placeholderScoreLogPath": "请输入scorelog.db文件路径",
      "placeholderSongDataPath": "请输入songdata.db文件路径"
    },
    "rules": {
      "missingRivalName": "好友名称不可为空",
      "missingScoreLogPath": "scorelog.db文件路径不可为空",
      "missingSongDataPath": "songdata.db文件路径不可为空"
    }
  }
}</i18n>