<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.addCustomTable')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelName')" path="Name">
        <n-input v-model:value="formData.Name" :placeholder="t('form.placeholderName')" />
      </n-form-item>
      <n-form-item :label="t('form.labelSymbol')" path="Symbol">
        <n-input v-model:value="formData.Symbol" :placeholder="t('form.placeholderCustomTableSymbol')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { AddCustomDiffTable } from '@wailsjs/go/main/App';
import { FormInst } from 'naive-ui';
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const show = defineModel<boolean>("show");
const emit = defineEmits<{
  (e: 'refresh'): void
}>();

const loading = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = reactive({
  Name: null,
  Symbol: null
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
    .then(() => {
      loading.value = true;
      AddCustomDiffTable(formData as any)
        .then(result => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          resetFormData();
          show.value = false;
          emit('refresh');
        })
        .catch(err => window.$notifyError(err))
        .finally(() => loading.value = false);
    })
    .catch((err) => { });
  return false;
}

function handleNegativeClick() {
  resetFormData();
}

function resetFormData() {
  formData.Name = null;
  formData.Symbol = null;
}
</script>
