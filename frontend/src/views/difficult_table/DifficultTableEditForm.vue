<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.editTable')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData">
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
import { QueryDiffTableInfoById, UpdateDiffTableHeader } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { FormInst, SelectOption } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose({ open });

const show = ref(false);
const loading = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = ref({
  ID: 0,
  TagColor: "",
  TagTextColor: "",
  NoTagBuild: 0,
});

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
    .then(async () => {
      console.log('form:', formData.value);
      const result = await UpdateDiffTableHeader(formData.value as any);
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

function open(headerId: number) {
  if (headerId == null || headerId == 0) {
    window.$notifyError(t('message.noChosenHeaderError'));
    show.value = false;
    return;
  }

  formData.value.ID = headerId;
  show.value = true;
  loading.value = true;
  QueryDiffTableInfoById(headerId)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      const data: dto.DiffTableHeaderDto = result.Data;
      formData.value.TagColor = data.TagColor;
      formData.value.TagTextColor = data.TagTextColor;
      formData.value.NoTagBuild = data.NoTagBuild;
    }).catch(err => {
      window.$notifyError(err);
      show.value = false;
    }).finally(() => {
      loading.value = false;
    });
}
</script>