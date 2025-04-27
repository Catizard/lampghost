<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('modal.labelAddress')" path="url">
        <n-input v-model:value="formData.url" :placeholder="t('modal.placeholderAddress')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { AddDiffTableHeader } from '@wailsjs/go/main/App';
import { FormInst, useNotification } from 'naive-ui';
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
  url: "",
});
const rules = {
  url: {
    required: true,
    message: t('rules.missingAddress'),
    trigger: ["input", "blur"],
  },
};

function handlePositiveClick(): boolean {
  formRef.value
    ?.validate()
    .then(() => {
      addDiffTableHeader(formData.value.url)
    })
    .catch((err) => { });
  return false;
}

function handleNegativeClick() {
  formData.value.url = "";
}

function addDiffTableHeader(url: string) {
  loading.value = true;
  AddDiffTableHeader(url)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      formData.value.url = "";
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
}
</script>

<i18n lang="json">{
  "en": {
    "modal": {
      "title": "Add a new table",
      "positiveText": "add",
      "negativeText": "cancel",
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
      "positiveText": "新增",
      "negativeText": "取消",
      "labelAddress": "地址",
      "placeholderAddress": "请输入地址"
    },
    "rules": {
      "missingAddress": "请输入地址"
    }
  }
}</i18n>