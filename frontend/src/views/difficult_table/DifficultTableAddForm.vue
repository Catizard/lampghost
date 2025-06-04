<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('modal.labelAddress')" path="HeaderUrl">
        <n-input v-model:value="formData.HeaderUrl" :placeholder="t('modal.placeholderAddress')" />
      </n-form-item>
      <n-form-item :label="t('modal.labelNoTagBuild')" path="NoTagBuild">
        <n-select v-model:value="formData.NoTagBuild" :options="tableTagOptions" />
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
import { FormInst, SelectOption } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
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

const tableTagOptions: Array<SelectOption> = [
  {
    label: t('options.noDisplay'),
    value: 1,
  },
  {
    label: t('options.display'),
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
      "placeholderAddress": "Input address",
      "labelNoTagBuild": "Display table tag or not",
      "labelTagColor": "Tag color",
      "labelTagTextColor": "Tag text color"
    },
    "rules": {
      "missingAddress": "Please input address"
    },
    "options": {
      "noDisplay": "Don't display",
			"display": "Display"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "新增难度表",
      "positiveText": "提交",
      "negativeText": "取消",
      "labelAddress": "地址",
      "placeholderAddress": "请输入地址",
      "labelNoTagBuild": "是否展示标签",
      "labelTagColor": "标签颜色",
      "labelTagTextColor": "标签嵌字颜色"
    },
    "rules": {
      "missingAddress": "请输入地址"
    },
    "options": {
			"noDisplay": "不展示",
			"display": "展示"
    }
  }
}</i18n>
