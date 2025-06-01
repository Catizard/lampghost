<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData">
      <n-form-item :label="t('modal.labelSymbol')" path="Symbol">
        <n-input v-model:value="formData.Symbol" clearable />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { FindCustomDiffTableByID, UpdateDiffTableHeader } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { FormInst } from 'naive-ui';
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
  Symbol: "",
});

function handlePositiveClick(): boolean {
  formRef.value
    ?.validate()
    .then(async () => {
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

function open(customTableID: number) {
  if (customTableID == null || customTableID == 0) {
    window.$notifyError(t('message.noChosenTableError'))
    show.value = false;
    return;
  }

  formData.value.ID = customTableID;
  show.value = true;
  loading.value = true;
  FindCustomDiffTableByID(customTableID)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      const data: dto.DiffTableHeaderDto = result.Data;
      data.Symbol = data.Symbol;
    }).catch(err => {
      window.$notifyError(err)
      show.value = false;
    }).finally(() => loading.value = false);
}
</script>

<i18n lang="json">{
  "en": {
    "modal": {
      "title": "Edit Custom Table",
      "positiveText": "Submit",
      "negativeText": "Cancel",
      "labelSymbol": "Symbol"
    },
    "message": {
      "noChosenTableError": "No custom table was chosed currently"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "修改自定义难度表",
      "positiveText": "新增",
      "negativeText": "取消",
      "labelSymbol": "标志"
    },
    "message": {
      "noChosenTableError": "当前没有选中任何自定义难度表"
    }
  }
}</i18n>
