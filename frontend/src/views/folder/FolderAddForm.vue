<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelName')" path="name">
        <n-input v-model:value="formData.name" :placeholder="t('form.placeholderName')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import { AddFolder } from '@wailsjs/go/main/App';
import { FormInst } from 'naive-ui';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const show = defineModel<boolean>("show");
const emit = defineEmits<{
  (e: 'refresh'): void
}>();

const { t } = useI18n();

const loading = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = ref({
  name: "",
});
const rules = {
  name: {
    required: true,
    message: t('message.missingName'),
    trigger: ["input", "blur"],
  },
};

function handlePositiveClick(): boolean {
  loading.value = true;
  formRef.value
    ?.validate()
    .then(() => {
      // TODO: Remove the magic 1 here, we don't allow add folder into
      // custom difficult table now
      return AddFolder({
        FolderName: formData.value.name,
        CustomTableID: 1
      } as any)
        .then((result) => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          console.log('reset')
          show.value = false;
          formData.value.name = null;
          emit('refresh');
        })
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
  return false;
}

function handleNegativeClick() {
  formData.value.name = "";
}

watch(show, newValue => {
  console.log("show =>", newValue);
})
</script>

<i18n lang="json">{
  "en": {
    "modal": {
      "title": "New folder",
      "positiveText": "Submit",
      "negativeText": "Cancel"
    },
    "form": {
      "labelName": "Name",
      "placeholderName": "Input name"
    },
    "message": {
      "missingName": "Please input name"
    }
  },
  "zh-CN": {
    "modal": {
      "title": "新增收藏夹",
      "positiveText": "提交",
      "negativeText": "取消"
    },
    "form": {
      "labelName": "名称",
      "placeholderName": "请输入名称"
    },
    "message": {
      "missingName": "请输入名称"
    }
  }
}</i18n>