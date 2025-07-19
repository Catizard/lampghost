<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.addCustomCourse')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelName')" path="Name">
        <n-input v-model:value="formData.Name" :placeholder="t('form.placeholderName')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { AddCustomCourse } from '@wailsjs/go/main/App';
import { FormInst } from 'naive-ui';
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const show = defineModel<boolean>("show");
const props = defineProps<{
  customTableId?: number,
}>();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();

const loading = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = reactive({
  Name: null,
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
      loading.value = true;
      const result = await AddCustomCourse({
        ...formData,
        CustomTableID: props.customTableId,
      } as any);
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      resetFormData();
      show.value = false;
      emit('refresh');
    })
    .catch((err) => { })
    .finally(() => loading.value = false);
  return false;
}

function handleNegativeClick() {
  resetFormData();
}

function resetFormData() {
  formData.Name = null;
}
</script>
