<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.addDifficultTable')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelAddress')" path="HeaderUrl">
        <n-input v-model:value="formData.HeaderUrl" :placeholder="t('form.placeholderAddress')" />
      </n-form-item>
      <n-form-item :label="t('form.labelNoTagBuild')" path="NoTagBuild">
        <n-select v-model:value="formData.NoTagBuild" :options="tableTagOptions" />
      </n-form-item>
      <n-form-item :label="t('form.labelTagColor')" path="TagColor">
        <n-color-picker v-model:value="formData.TagColor" :show-alpha="false" :modes="['hex', 'rgb']" />
      </n-form-item>
      <n-form-item :label="t('form.labelTagTextColor')" path="TagTextColor">
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
    message: t('rule.missingAddress'),
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
