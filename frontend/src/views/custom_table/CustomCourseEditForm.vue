<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.editCustomCourse')"
    :positive-text="t('button.submit')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('form.labelName')" path="Name">
        <n-input v-model:value="formData.Name" clearable />
      </n-form-item>
      <n-form-item :label="t('form.labelConstraints')" path="Constraint">
        <n-select v-model:value="formData.Constraint" :options="constraintOptions" multiple />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { FindCustomCourseByID, UpdateCustomCourse } from '@wailsjs/go/main/App';
import { entity } from '@wailsjs/go/models';
import { FormInst, SelectOption } from 'naive-ui';
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const emit = defineEmits<{
  (e: 'refresh'): void
}>();
defineExpose(({ open }));

const loading = ref(false);
const show = ref(false);

const constraintOptions: SelectOption[] = [
  { label: "grade", value: "grade" },
  { label: "grade_mirror", value: "grade_mirror" },
  { label: "grade_random", value: "grade_random" },
  { label: "no_speed", value: "no_speed" },
  { label: "no_good", value: "no_good" },
  { label: "no_great", value: "no_great" },
  { label: "gauge_lr2", value: "gauge_lr2" },
  { label: "gauge_5k", value: "gauge_5k" },
  { label: "gauge_7k", value: "gauge_7k" },
  { label: "gauge_9k", value: "gauge_9k" },
  { label: "gauge_24k", value: "gauge_24k" },
  { label: "ln", value: "ln" },
  { label: "cn", value: "cn" },
  { label: "hcn", value: "hcn" },
];

const formRef = ref<FormInst | null>(null);
const formData = reactive({
  ID: 0,
  Name: "",
  Constraint: []
});
const rules = {
  Name: {
    message: t('rule.missingName'),
    trigger: ["input", "blur"]
  }
};

function handlePositiveClick() {
  formRef.value
    ?.validate()
    .then(async () => {
      const result = await UpdateCustomCourse({
        ID: formData.ID,
        Name: formData.Name,
        Constraints: formData.Constraint.join(",")
      } as any);
      if (result.Code != 200) {
        throw result.Msg;
      }
      show.value = false;
      emit('refresh');
    }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function handleNegativeClick() {
  show.value = false;
}

function open(customCourseID: number) {
  if (customCourseID == null || customCourseID == 0) {
    window.$notifyError(t('message.noChosenCustomCourseError'));
    show.value = false;
    return;
  }

  formData.ID = customCourseID;
  show.value = true;
  loading.value = true;
  FindCustomCourseByID(customCourseID).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    const data: entity.CustomCourse = result.Data;
    formData.Name = data.Name;
    if (data.Constraints != "") {
      formData.Constraint = data.Constraints.split(",")
    }
  }).catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}
</script>
