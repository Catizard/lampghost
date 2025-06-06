<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('modal.labelName')" path="Name">
        <n-input v-model:value="formData.Name" :placeholder="t('modal.placeholderName')" />
      </n-form-item>
      <n-form-item :label="t('modal.labelSymbol')" path="Symbol">
        <n-input v-model:value="formData.Symbol" :placeholder="t('modal.placeholderSymbol')" />
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
    message: t('rules.missingName'),
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

<i18n lang="json">{
  "en": {
    "modal": {
      "title": "Add a new custom table",
      "positiveText": "Submit",
      "negativeText": "Cancel",
      "labelName": "Name",
      "placeholderName": "Please Input name",
      "labelSymbol": "Symbol",
      "placeHolderSymbol": "Customize Table Symbol"
    },
    "rules": {
      "missingName": "Name cannot be empty"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "新增自定义难度表",
      "positiveText": "提交",
      "negativeText": "取消",
      "labelName": "名称",
      "placeholderName": "请输入名称",
      "labelSymbol": "标志",
      "placeHolderSymbol": "自定义难度表标志"
    },
    "rules": {
      "missingName": "名称不能为空"
    }
  }
}</i18n>
