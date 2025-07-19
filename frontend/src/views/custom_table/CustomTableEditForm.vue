<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.editCustomTable')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelName')" path="Name">
        <n-input v-model:value="formData.Name" :placeholder="t('form.placeholderName')" clearable />
      </n-form-item>
      <n-form-item :label="t('form.labelSymbol')" path="Symbol">
        <n-input v-model:value="formData.Symbol" :placeholder="t('form.placeholderCustomTableSymbol')" clearable />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { FindCustomDiffTableByID, UpdateCustomDiffTable } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { FormInst } from 'naive-ui';
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose({ open });

const show = ref(false);
const loading = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = reactive({
  ID: 0,
  Name: "",
  Symbol: "",
});

const rules = {
  Name: {
    required: true,
    message: t('rule.missingName'),
    trigger: ["input", "blur"],
  },
};

function handlePositiveClick(): boolean {
  formRef.value
    ?.validate()
    .then(async () => {
      const result = await UpdateCustomDiffTable(formData as any);
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
}

function open(customTableID: number) {
  if (customTableID == null || customTableID == 0) {
    window.$notifyError(t('message.noChosenCustomTableError'))
    show.value = false;
    return;
  }

  formData.ID = customTableID;
  show.value = true;
  loading.value = true;
  FindCustomDiffTableByID(customTableID)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      const data: dto.CustomDiffTableDto = result.Data;
      formData.Symbol = data.Symbol;
      formData.Name = data.Name;
    }).catch(err => {
      window.$notifyError(err)
      show.value = false;
    }).finally(() => loading.value = false);
}
</script>
