<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.addPlayer')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelRivalName')" path="Name">
        <n-input v-model:value="formData.Name" :placeholder="t('form.placeholderRivalName')" />
      </n-form-item>
      <n-form-item :label="t('form.labelScorelogPath')" path="ScoreLogPath">
        <n-button type="info" @click="chooseFile('Choose scorelog.db', 'scorelogPath')">
          {{ t('button.chooseFile') }}
        </n-button>
        <n-divider vertical />
        <n-input v-model:value="formData.ScoreLogPath" :placeholder="t('form.placeholderScorelogPath')" />
      </n-form-item>
      <n-form-item :label="t('form.labelScoredatalogPath')" path="ScoreDataLogPath">
        <n-input disabled v-model:value="formData.ScoreDataLogPath"
          :placeholder="t('form.placeholderScoredatalogPath')" />
      </n-form-item>
      <n-form-item :label="t('form.labelSongdataPath')" path="SongDataPath">
        <n-input disabled v-model:value="formData.SongDataPath" :placeholder="t('form.placeholderSongdataPath')" />
      </n-form-item>
      <n-form-item :label="t('form.labelScoredataPath')" path="ScoreDataPath">
        <n-input disabled v-model:value="formData.ScoreDataPath" :placeholder="t('form.placeholderScoredataPath')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import { AddRivalInfo } from '@wailsjs/go/main/App';
import { OpenFileDialog } from '@wailsjs/go/main/App';
import { FormInst } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const show = defineModel<boolean>("show");
const emit = defineEmits<{
  (e: 'refresh'): void
}>();

const { t } = useI18n();
const loading = ref(false);

const formRef = ref<FormInst | null>(null);
const formData = ref({
  Name: null,
  ScoreLogPath: null,
  SongDataPath: null,
  ScoreDataPath: null,
  ScoreDataLogPath: null,
  Type: "beatoraja" // TODO: We don't support import LR2 user as rival currently
});
const rules = {
  Name: {
    required: true,
    message: t('rule.missingUserName'),
    trigger: ["input", "blur"],
  },
  ScoreLogPath: {
    required: true,
    message: t('rule.missingScorelogPath'),
    trigger: ["input", "blur"],
  },
  // SongDataPath: {
  // 	required: true,
  // 	message: t('rule.missingSongdataPath'),
  // 	trigger: ["input", "blur"],
  // },
  // ScoreDataLogPath: {
  //   required: false,
  // },
};

function handlePositiveClick(): boolean {
  loading.value = true;
  formRef.value
    ?.validate()
    .then(() => {
      return AddRivalInfo(formData.value as any)
        .then(result => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          show.value = false;
          emit('refresh')
          refreshFormData();
        });
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
  return false;
}

function handleNegativeClick() {
  refreshFormData();
}

// target == "scorelogPath" | "songdataPath"
function chooseFile(title, target: "scorelogPath" | "songdataPath" | "scoredatalogPath" | "scoredataPath") {
  OpenFileDialog(title)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Data != null && result.Data != undefined && result.Data != "") {
        switch (target) {
          case "scorelogPath": formData.value.ScoreLogPath = result.Data; break;
          case "songdataPath": formData.value.SongDataPath = result.Data; break;
          case "scoredatalogPath": formData.value.ScoreDataLogPath = result.Data; break;
          case "scoredataPath": formData.value.ScoreDataPath = result.Data; break;
        }
      }
    }).catch(err => window.$notifyError(err));
}

function refreshFormData() {
  formData.value.Name = null;
  formData.value.ScoreLogPath = null;
  formData.value.SongDataPath = null;
  formData.value.ScoreDataLogPath = null;
  formData.value.ScoreDataPath = null;
}
</script>
