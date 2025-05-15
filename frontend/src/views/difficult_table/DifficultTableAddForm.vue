<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('modal.labelAddress')" path="HeaderUrl">
        <n-input v-model:value="formData.HeaderUrl" :placeholder="t('modal.placeholderAddress')" />
      </n-form-item>
      <n-form-item :label="t('modal.labelNoTagBuild')" path="NoTagBuild">
        <n-select v-model:value="formData.NoTagBuild" :options="yesnoOptions"></n-select>
      </n-form-item>
      <n-form-item :label="t('modal.labelTagColor')" path="TagColor">
        <n-color-picker v-model:value="formData.TagColor" :show-alpha="false" :modes="['hex', 'rgb']" />
      </n-form-item>
      <n-form-item :label="t('modal.labelTagTextColor')" path="TagTextColor">
        <n-color-picker v-model:value="formData.TagTextColor" :show-alpha="false" :modes="['hex', 'rgb']" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { AddDiffTableHeader } from '@wailsjs/go/main/App';
import { FormInst, SelectOption, useNotification } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const notification = useNotification();
const show = defineModel<boolean>("show");
const emit = defineEmits<{
  (e: 'refresh'): void
}>();

const loading = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = ref({
  HeaderUrl: "",
  TagColor: "",
  TagTextColor: "",
  NoTagBuild: 0,
});
const rules = {
  HeaderUrl: {
    required: true,
    message: t('rules.missingAddress'),
    trigger: ["input", "blur"],
  },
};

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

function handlePositiveClick(): boolean {
  formRef.value
    ?.validate()
    .then(() => {
      loading.value = true;
      AddDiffTableHeader(formData.value as any)
        .then(result => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          resetFormData();
          show.value = false;
          emit('refresh');
        })
        .catch((err) => {
          notification.error({
            content: err,
            duration: 3000,
            keepAliveOnHover: true
          })
        }).finally(() => loading.value = false);
    })
    .catch((err) => { });
  return false;
}

function handleNegativeClick() {
  resetFormData();
}

function resetFormData() {
  formData.value.HeaderUrl = "";
  formData.value.TagColor = "";
  formData.value.TagTextColor = "";
  formData.value.NoTagBuild = 0;
}
</script>

<i18n lang="json">{
  "en": {
    "modal": {
      "title": "Add a new table",
      "positiveText": "Submit",
      "negativeText": "Cancel",
      "labelAddress": "Address",
      "placeholderAddress": "Input address"
    },
    "rules": {
      "missingAddress": "Please input address"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "新增难度表",
      "positiveText": "提交",
      "negativeText": "取消",
      "labelAddress": "地址",
      "placeholderAddress": "请输入地址"
    },
    "rules": {
      "missingAddress": "请输入地址"
    }
  }
}</i18n>