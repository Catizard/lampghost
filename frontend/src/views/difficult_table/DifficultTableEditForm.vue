<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData">
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

<i18n lang="json">{
  "en": {
    "modal": {
      "title": "Edit table",
      "positiveText": "submit",
      "negativeText": "cancel",
      "labelNoTagBuild": "Display table tag or not",
      "labelTagColor": "Tag color",
      "labelTagTextColor": "Tag text color"
    },
    "options": {
      "noDisplay": "Don't display",
      "display": "Display"
    },
    "message": {
      "noChosenHeaderError": "No header was chosed currently"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "修改难度表",
      "positiveText": "提交",
      "negativeText": "取消",
      "labelNoTagBuild": "是否展示标签",
      "labelTagColor": "标签颜色",
      "labelTagTextColor": "标签嵌字颜色"
    },
    "options": {
      "noDisplay": "不展示",
      "display": "展示"
    },
    "message": {
      "noChosenHeaderError": "当前没有选择任何难度表"
    }
  }
}</i18n>
