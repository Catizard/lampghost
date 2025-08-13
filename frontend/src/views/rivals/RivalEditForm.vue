<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.editPlayer')"
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
      <n-form-item :disabled="!formData.MainUser" :label="t('form.labelSongdataPath')" path="SongDataPath">
        <n-button :disabled="!formData.MainUser" type="info" @click="chooseFile('Choose songdata.db', 'songdataPath')">
          {{ t('button.chooseFile') }}
        </n-button>
        <n-divider vertical />
        <n-input :disabled="!formData.MainUser" v-model:value="formData.SongDataPath"
          :placeholder="t('form.placeholderSongdataPath')" />
      </n-form-item>
      <n-form-item :label="t('form.labelScoredatalogPath')" path="ScoreDataLogPath">
        <n-button :disabled="!formData.MainUser" type="info"
          @click="chooseFile('Choose scoredatalog.db', 'scoredatalogPath')">
          {{ t('button.chooseFile') }}
        </n-button>
        <n-divider vertical />
        <n-input :disabled="!formData.MainUser" v-model:value="formData.ScoreDataLogPath"
          :placeholder="t('form.placeholderScoredatalogPath')" />
      </n-form-item>
      <n-form-item :label="t('form.labelScoredataPath')" path="ScoreDataPath">
        <n-button :disabled="!formData.MainUser" type="info" @click="chooseFile('Choose score.db', 'scoredataPath')">
          {{ t('button.chooseFile') }}
        </n-button>
        <n-divider vertical />
        <n-input :disabled="!formData.MainUser" v-model:value="formData.ScoreDataPath"
          :placeholder="t('form.placeholderScoredataPath')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import { QueryUserInfoByID, UpdateRivalInfo } from '@wailsjs/go/main/App';
import { OpenFileDialog } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { FormInst } from 'naive-ui';
import { reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose({ open });

const show = ref(false);

function open(rivalID: number) {
  if (rivalID == null || rivalID == 0) {
    window.$notifyError(t('message.noChosenRivalError'));
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
      formData.value.ScoreDataLogPath = data.ScoreDataLogPath;
      formData.value.ScoreDataPath = data.ScoreDataPath;
      formData.value.MainUser = data.MainUser;
    }).catch(err => {
      window.$notifyError(err);
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
  ScoreDataLogPath: null,
  ScoreDataPath: null,
  MainUser: false,
});
const rules = reactive({
  Name: {
    ignored: false,
    required: true,
    message: t('rule.missingUserName'),
    trigger: ["input", "blur"],
  },
  ScoreLogPath: {
    ignored: false,
    required: true,
    message: t('rule.missingScorelogPath'),
    trigger: ["input", "blur"],
  },
  ScoreDataLogPath: {
    ignored: true,
    required: true,
    message: t('rule.missingScoredatalogPath'),
    trigger: ["input", "blur"],
  },
  SongDataPath: {
    ignored: true,
    required: true,
    message: t('rule.missingSongdataPath'),
    trigger: ["input", "blur"],
  },
  ScoreDataPath: {
    ignored: true,
    required: true,
    message: t('rule.missingScoredataPath'),
    trigger: ["input", "blur"],
  }
});

function handlePositiveClick(): boolean {
  loading.value = true;
  formRef.value
    ?.validate(() => { }, rule => {
      if (formData.value.MainUser) {
        return true; // Everything
      }
      return (rule as any).ignored == false;
    })
    .then(async () => {
      const result = await UpdateRivalInfo(formData.value as any);
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      show.value = false;
      emit('refresh');
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
  return false;
}

function handleNegativeClick() {
  formData.value.ID = 0;
  formData.value.Name = null;
  formData.value.ScoreLogPath = null;
  formData.value.SongDataPath = null;
  formData.value.ScoreDataLogPath = null;
  formData.value.ScoreDataPath = null;
  formData.value.MainUser = false;
}

function chooseFile(title: string, target: "scorelogPath" | "songdataPath" | "scoredatalogPath" | "scoredataPath") {
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
        } else if (target == "scoredatalogPath") {
          formData.value.ScoreDataLogPath = result.Data;
        } else if (target == "scoredataPath") {
          formData.value.ScoreDataPath = result.Data;
        }
      }
    }).catch(err => window.$notifyError(err));
}

watch(() => formData.value.MainUser, newValue => {
  rules.SongDataPath.required = newValue;
})
</script>
