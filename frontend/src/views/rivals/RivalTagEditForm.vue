<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.editRivalTag')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelEnable')">
        <n-switch v-model:value="formData.Enabled" name="enable">{{ t('form.enable') }}</n-switch>
      </n-form-item>
      <n-form-item :label="t('form.labelSymbol')" path="Name">
        <n-input v-model:value="formData.Symbol" clearable />
      </n-form-item>
      <n-form-item :label="t('form.labelTagRecordTime')" path="RecordTimestamp">
        <!-- NOTE: Cannot edit generated tag's record time -->
        <n-date-picker :disabled="formData.Generated" clearable v-model:value="formData.RecordTimestamp"
          type="datetime" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { FindRivalTagByID, UpdateCustomCourse, UpdateRivalTag } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { FormInst } from 'naive-ui';
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose(({ open }));

const loading = ref(false);
const show = ref(false);

const formRef = ref<FormInst | null>(null);
const formData = reactive({
  ID: 0,
  Symbol: "",
  RecordTimestamp: null,
  Generated: false,
  Enabled: false
});
const rules = {
  Symbol: {
    message: t('rule.missingSymbol'),
    trigger: ["input", "blur"]
  }
};

function handlePositiveClick() {
  formRef.value
    ?.validate()
    .then(async () => {
      const result = await UpdateRivalTag(formData as any);
      if (result.Code != 200) {
        throw result.Msg;
      }
      show.value = false;
      emit('refresh');
    }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function handleNegativeClick() {
  show.value = false;
  formData.ID = 0;
  formData.Generated = false;
  formData.RecordTimestamp = null;
  formData.Symbol = "";
  formData.Enabled = false;
}

function open(rivalTagId: number) {
  if (rivalTagId == null || rivalTagId == 0) {
    window.$notifyError(t('message.noChosenRivalTagError'));
    show.value = false;
    return;
  }

  formData.ID = rivalTagId;
  show.value = true;
  loading.value = true;
  FindRivalTagByID(rivalTagId).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    const data: dto.RivalTagDto = result.Data;
    formData.Symbol = data.Symbol;
    formData.RecordTimestamp = data.RecordTimestamp;
    formData.Generated = data.Generated;
    formData.Enabled = data.Enabled;
  })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}
</script>
