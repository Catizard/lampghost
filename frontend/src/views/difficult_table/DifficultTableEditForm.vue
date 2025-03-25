<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData">
      <n-form-item :label="t('modal.labelEnableFallbackSort')" path="url">
        <n-select v-model:value="formData.EnableFallbackSort" :options="yesnoOptions" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { QueryDiffTableInfoById, UpdateDiffTableHeader } from '@wailsjs/go/controller/DiffTableController';
import { dto, vo } from '@wailsjs/go/models';
import { FormInst, SelectOption, useNotification } from 'naive-ui';
import { onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const notification = useNotification();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose({ open });

const show = ref(false);
const loading = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = ref({
  ID: 0,
  EnableFallbackSort: 0,
});

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
      return UpdateDiffTableHeader(formData.value as any)
        .then(result => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          formData.value.EnableFallbackSort = 0;
          show.value = false;
          emit('refresh');
        });
    })
    .catch((err) => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true
      })
    }).finally(() => loading.value = false);
  return false;
}

function handleNegativeClick() {
  formData.value.EnableFallbackSort = 0;
}

function open(headerId: number) {
  if (headerId == null || headerId == 0) {
    notification.error({
      content: t('message.noChosenHeaderError'),
      duration: 3000,
      keepAliveOnHover: true
    });
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
      formData.value.EnableFallbackSort = data.EnableFallbackSort;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      });
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
      "labelEnableFallbackSort": "Enable fallback sort strategy"
    },
    "options": {
      "yes": "Yes",
      "no": "No"
    },
    "message": {
      "noChosenHeaderError": "No header was chosed currently"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "修改难度表",
      "positiveText": "新增",
      "negativeText": "取消",
      "labelEnableFallbackSort": "是否启用保底排序机制"
    },
    "options": {
      "yes": "是",
      "no": "否"
    },
    "message": {
      "noChosenHeaderError": "当前没有选择任何难度表"
    }
  }
}</i18n>