<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.editPlayerReverseImport')"
    :positive-text="t('button.submit')" :negative-text="t('button.submit')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData">
      <n-form-item :label="t('form.labelReverseImport')" path="ReverseImport">
        <n-select style="width: 300px;" v-model:value="formData.ReverseImport" :options="reverseImportOptions" />
      </n-form-item>
      <n-form-item :label="t('form.labelLockTag')" path="LockTagID">
        <SelectRivalTag v-model:value="formData.LockTagID" :rivalId="formData.ID" width="300px" clearable />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import SelectRivalTag from '@/components/rivals/SelectRivalTag.vue';
import { QueryUserInfoByID, UpdateRivalReverseImportInfo } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { FormInst, SelectOption } from 'naive-ui';
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose({ open });

const show = ref(false);
const loading = ref(false);
const reverseImportOptions: SelectOption[] = [
  {
    label: t('options.noReverseImport'),
    value: 0,
  },
  {
    label: t('options.reverseImport'),
    value: 1,
  },
];

function open(rivalID: number) {
  if (rivalID == null || rivalID == 0) {
    window.$notifyError(t('message.noChosenRivalError'));
    show.value = false;
    return
  }

  formData.ID = rivalID;
  show.value = true;
  loading.value = true;

  QueryUserInfoByID(rivalID)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      const data: dto.RivalInfoDto = result.Data;
      formData.LockTagID = data.LockTagID == 0 ? null : data.LockTagID;
      formData.ReverseImport = data.ReverseImport;
    }).catch(err => {
      window.$notifyError(err);
      show.value = false;
    }).finally(() => {
      loading.value = false;
    });
}

const formRef = ref<FormInst | null>(null);
const formData = reactive({
  ID: 0,
  LockTagID: null,
  ReverseImport: 0
});

function handlePositiveClick(): boolean {
  loading.value = true;
  formRef.value
    ?.validate()
    .then(async () => {
      const result = await UpdateRivalReverseImportInfo(formData as any);
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
  formData.ID = 0;
  formData.ReverseImport = 0;
  formData.LockTagID = null;
}
</script>
