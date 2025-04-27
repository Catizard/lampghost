<template>
  <perfect-scrollbar>
    <n-flex justify="space-start">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">{{ t('title') }}</n-text>
      </n-h1>
    </n-flex>
    <n-flex justify="flex-start" style="margin-bottom: 15px;">
      <n-select :loading="tableLoading" v-model:value="currentRivalID" :options="rivalOptions" style="width: 200px;" />
      <n-button style="margin-left: auto;" type="primary" @click="showAddModal = true">{{ t('button.add') }}</n-button>
    </n-flex>
    <n-data-table remote :columns="columns" :data="data" :pagination="pagination" :loading="tableLoading"
      :row-key="(row: dto.RivalTagDto) => row.ID" />
  </perfect-scrollbar>
  <n-modal v-model:show="showAddModal" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('modal.labelTagName')" path="TagName">
        <n-input v-model:value="formData.TagName" :placeholder="t('modal.placeholderTagName')" />
      </n-form-item>
      <n-form-item :label="t('modal.labelRecordTime')" path="RecordTimestamp">
        <n-date-picker clearable v-model:value="formData.RecordTimestamp" type="datetime" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { FindRivalInfoList } from '@wailsjs/go/main/App';
import { AddRivalTag, DeleteRivalTagByID, QueryRivalTagPageList, RevertRivalTagEnabledState } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, FormInst, FormRules, NButton, SelectOption, useDialog, useNotification } from 'naive-ui';
import { h, reactive, Ref, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { useI18n } from 'vue-i18n';
import YesNotTag from '@/components/YesNotTag.vue';

const i18n = useI18n();
const { t } = i18n;
const dialog = useDialog();
const notification = useNotification();

const tableLoading = ref(false);
const currentRivalID: Ref<number | null> = ref(null);
const rivalOptions: Ref<Array<SelectOption>> = ref([]);
function loadRivalOptions() {
  tableLoading.value = true;
  FindRivalInfoList()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Rows.length == 0) {
        return Promise.reject(t('message.noRivalError'));
      }
      rivalOptions.value = result.Rows.map((row: dto.RivalInfoDto) => {
        return {
          label: row.Name,
          value: row.ID,
        } as SelectOption
      });
      currentRivalID.value = rivalOptions.value[0].value as number;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      })
    }).finally(() => {
      tableLoading.value = false;
    });
}
loadRivalOptions();

const columns: DataTableColumns<dto.RivalTagDto> = [
  { title: t('column.tagName'), key: "TagName", width: "200px", ellipsis: { tooltip: true } },
  {
    title: t('column.generated'), key: "Generated",
    render: (row: dto.RivalTagDto) => {
      return h(
        YesNotTag,
        { state: row.Generated }
      );
    }
  },
  {
    title: t('column.enabled'), key: "Enabled",
    render: (row: dto.RivalTagDto) => {
      return h(
        YesNotTag,
        { state: row.Enabled }
      );
    }
  },
  {
    title: t('column.tagTime'),
    key: "RecordTime",
    render: (row: dto.RivalTagDto) => dayjs(row.RecordTime).format('YYYY-MM-DD HH:mm:ss')
  },
  {
    title: t('column.actions'), key: "Actions",
    render: (row: dto.RivalTagDto) => {
      const deleteTagButton = row.Generated == false ? h(
        NButton,
        {
          strong: true,
          tertiary: true,
          size: 'small',
          type: "error",
          onClick: () => deleteTag(row),
        },
        { default: () => t('button.delete') }
      ) : null;
      const revertEnabledStateButton = h(
        NButton,
        {
          strong: true,
          tertiary: true,
          size: "small",
          onClick: () => revertTagEnabledState(row),
        },
        { default: () => row.Enabled ? t('button.disable') : t('button.enable') }
      );
      if (deleteTagButton != null) {
        return [deleteTagButton, revertEnabledStateButton];
      }
      return revertEnabledStateButton;
    }
  }
];
const data: Ref<Array<dto.RivalTagDto>> = ref([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
    loadData();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    loadData();
  },
  reset: () => {
    pagination.page = 1;
    loadData();
  }
});
function loadData() {
  tableLoading.value = true;
  QueryRivalTagPageList({
    RivalId: currentRivalID.value,
    Pagination: pagination,
  } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows];
      pagination.pageCount = result.Pagination.pageCount;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      });
    }).finally(() => {
      tableLoading.value = false;
    });
}

const showAddModal = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = ref({
  TagName: null,
  RecordTimestamp: null
});
const rules: FormRules = {
  RecordTimestamp: {
    type: 'number',
    required: true,
    message: t('rules.missingRecordTime'),
    trigger: ["input", "blur"]
  }
}
function handlePositiveClick(): boolean {
  formRef.value
    ?.validate()
    .then(() => {
      AddRivalTag(formData.value as any)
        .then(result => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          showAddModal.value = false;
        }).catch(err => {
          notification.error({
            content: err,
            duration: 3000,
            keepAliveOnHover: true
          });
        });
    })
    .catch((err) => { });
  return false;
}
function handleNegativeClick() {
  formData.value.TagName = null;
  formData.value.RecordTimestamp = null;
}

function deleteTag(tag: dto.RivalTagDto) {
  dialog.warning({
    title: t('deleteDialog.title'),
    positiveText: t('deleteDialog.positiveText'),
    negativeText: t('deleteDialog.negativeText'),
    onPositiveClick: () => {
      DeleteRivalTagByID(tag.ID)
        .then(result => {
          if (result.Code != 200) {
            return Promise.reject(result.Msg);
          }
          loadData();
        }).catch(err => {
          notification.error({
            content: err,
            duration: 3000,
            keepAliveOnHover: true,
          });
        });
    }
  })
}

function revertTagEnabledState(tag: dto.RivalTagDto) {
  RevertRivalTagEnabledState(tag.ID)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      loadData();
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true
      });
    });
}

// Watch: Whenever changing current rival, reset current page to the first one
// NOTE: There is no need to call loadData(), reset method would call it
watch(currentRivalID, () => {
  pagination.reset();
})
</script>

<i18n lang="json">{
  "en": {
    "title": "Rival Tags",
    "column": {
      "tagName": "Tag Name",
      "enabled": "Enabled",
      "generated": "Generated",
      "tagTime": "Tag Time",
      "actions": "Actions"
    },
    "message": {
      "noRivalError": "FATAL ERROR: no rival data found"
    },
    "modal": {
      "title": "Add Custom Tag",
      "positiveText": "Submit",
      "negativeText": "Cancel",
      "labelTagName": "Tag Name",
      "labelRecordTime": "Tag Time",
      "placeholderTagName": "Please input tag name",
      "placeholderRecordTime": "Please input tag time"
    },
    "button": {
      "add": "Add Custom Tag",
      "delete": "Delete",
      "enable": "Enable",
      "disable": "Disable"
    },
    "rules": {
      "missingRecordTime": "Tag time cannot be empty"
    },
    "deleteDialog": {
      "title": "Confirm to delete?",
      "positiveText": "Yes",
      "negativeText": "No"
    }
  },
  "zh-CN": {
    "title": "玩家标签",
    "column": {
      "tagName": "标签名称",
      "enabled": "是否启用",
      "generated": "自动生成",
      "tagTime": "标签时间",
      "actions": "操作"
    },
    "message": {
      "noRivalError": "未知错误: 找不到任何玩家信息?"
    },
    "modal": {
      "title": "添加自定义标签",
      "positiveText": "提交",
      "negativeText": "取消",
      "labelTagName": "名称",
      "labelRecordTime": "标签时间",
      "placeholderTagName": "请输入标签名称",
      "placeholderRecordTime": "请输入标签时间"
    },
    "button": {
      "add": "添加自定义标签",
      "enable": "启用",
      "disable": "禁用"
    },
    "rules": {
      "missingRecordTime": "标签时间不可为空"
    },
    "deleteDialog": {
      "title": "确定要删除吗?",
      "positiveText": "是",
      "negativeText": "否"
    }
  }
}</i18n>
