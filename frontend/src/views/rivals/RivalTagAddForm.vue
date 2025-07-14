<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.addPlayerTag')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelTagName')" path="TagName">
        <n-input v-model:value="formData.TagName" :placeholder="t('form.placeholderTagName')" />
      </n-form-item>
      <n-form-item :label="t('form.labelTagRecordTime')" path="RecordTimestamp">
        <n-date-picker clearable v-model:value="formData.RecordTimestamp" type="datetime" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { AddRivalTag } from '@wailsjs/go/main/App';
import { FormInst, FormRules } from 'naive-ui';
import { ref, defineEmits } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { rivalId } = defineProps<{
  rivalId?: number
}>();
const emit = defineEmits<{
  (e: 'refresh'): void
}>()
const loading = ref(false);
const show = defineModel<boolean>("show");

const formRef = ref<FormInst | null>(null);
const formData = ref({
  RivalID: null,
  TagName: null,
  RecordTimestamp: null
});
const rules: FormRules = {
  RecordTimestamp: {
    type: 'number',
    required: true,
    message: t('rule.missingTagRecordTime'),
    trigger: ["input", "blur"]
  }
}
function handlePositiveClick(): boolean {
  loading.value = true;
  formRef.value
    ?.validate()
    .then(async () => {
      const result = await AddRivalTag({
        RivalId: rivalId,
        ...formData.value
      } as any);
      if (result.Code != 200) {
        throw result.Msg;
      }
      show.value = false;
      emit('refresh');
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
  return false;
}
function handleNegativeClick() {
  formData.value.TagName = null;
  formData.value.RecordTimestamp = null;
}
</script>
