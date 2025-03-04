<template>
  <perfect-scrollbar>
    <n-space justify="space-between">
      <n-h1 prefix="bar" style="text-align: start">
        <n-text type="primary">{{ t('title') }}</n-text>
      </n-h1>
      <n-button :loading="loading" type="primary" @click="showAddModal = true">{{ t('button.add') }}</n-button>
    </n-space>
    <n-data-table :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
      :row-key="(row: dto.DiffTableHeaderDto) => row.ID" />
  </perfect-scrollbar>

  <n-modal v-model:show="showAddModal" preset="dialog" :title="t('modal.title')"
    :positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item :label="t('modal.labelAddress')" path="url">
        <n-input v-model:value="formData.url" :placeholder="t('modal.placeholderAddress')" />
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<script lang="ts" setup>
import type { DataTableColumns, FormInst } from "naive-ui";
import { NButton, NDataTable, useDialog } from "naive-ui";
import { useNotification } from "naive-ui";
import { h, Ref, ref } from "vue";
import {
  AddDiffTableHeader,
  DelDiffTableHeader,
  FindDiffTableHeaderTree
} from "@wailsjs/go/controller/DiffTableController";
import { dto } from "@wailsjs/go/models";
import { useI18n } from "vue-i18n";

const i18n = useI18n();
const { t } = i18n;
const showAddModal = ref(false);

const formRef = ref<FormInst | null>(null);
const formData = ref({
  url: "",
});
const rules = {
  url: {
    required: true,
    message: t('rules.missingAddress'),
    trigger: ["input", "blur"],
  },
};

const notification = useNotification();
const dialog = useDialog();
const loading = ref(false);
loadDiffTableData();

function createColumns({
  deleteHeader,
}: {
  deleteHeader: (row: dto.DiffTableHeaderDto) => void;
}): DataTableColumns<dto.DiffTableHeaderDto> {
  return [
    { title: t('column.name'), key: "Name", },
    { title: t('column.url'), key: "HeaderUrl", },
    {
      title: t('column.actions'),
      key: "actions",
      render(row) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: "small",
            onClick: () => deleteHeader(row),
          },
          { default: () => t('button.delete') },
        );
      },
    },
  ];
}

let data: Ref<Array<any>> = ref([]);
const pagination = false as const;
const columns = createColumns({
  deleteHeader(row: any) {
    dialog.warning({
      title: t('deleteDialog.title'),
      positiveText: t('deleteDialog.postiveText'),
      negativeText: t('deleteDialog.negativeText'),
      onPositiveClick: () => {
        delDiffTableHeader(row.ID)
      }
    })
  }
});

function addDiffTableHeader(url: string) {
  loading.value = true;
  AddDiffTableHeader(url)
    .then((result) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      formData.value.url = "";
      loadDiffTableData();
    })
    .catch((err) => {
      notifyError(t('message.addTableFailedPrefix') + err)
      loadDiffTableData();
    }).finally(() => loading.value = false);
}

function delDiffTableHeader(id: number) {
  loading.value = true;
  DelDiffTableHeader(id)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      notifySuccess(t('message.deleteSuccess'))
      loadDiffTableData();
    }).catch(err => {
      notifyError(err)
      loadDiffTableData();
    }).finally(() => loading.value = false);
}

function handlePositiveClick(): boolean {
  formRef.value
    ?.validate()
    .then(() => {
      addDiffTableHeader(formData.value.url);
      showAddModal.value = false;
    })
    .catch((err) => { });
  return false;
}

function handleNegativeClick() {
  formData.value.url = "";
}

function loadDiffTableData() {
  loading.value = true;
  FindDiffTableHeaderTree(null)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows]
    }).catch((err) => {
      notification.error({
        content: t('message.loadTableDataFailedPrefix') + err,
        duration: 5000,
        keepAliveOnHover: true
      })
    }).finally(() => loading.value = false);
}

function notifySuccess(msg: string) {
  notification.success({
    content: msg,
    duration: 5000,
    keepAliveOnHover: true
  })
}

function notifyError(msg: string) {
  notification.error({
    content: msg,
    duration: 5000,
    keepAliveOnHover: true
  })
}
</script>

<i18n lang="json">{
  "en": {
    "title": "Table Management",
    "button": {
      "add": "Add Table",
      "delete": "Delete"
    },
    "modal": {
      "title": "Add a new table",
      "positiveText": "add",
      "negativeText": "cancel",
      "labelAddress": "Address",
      "placeholderAddress": "Input address"
    },
    "rules": {
      "missingAddress": "Please input address"
    },
    "column": {
      "name": "Name",
      "url": "URL",
      "actions": "Actions"
    },
    "deleteDialog": {
      "title": "Confirm to delete?",
      "positiveText": "Yes",
      "negativeText": "No"
    },
    "message": {
      "addTableFailedPrefix": "Failed to add table, error message: ",
      "deleteSuccess": "Deleted successfully",
      "loadTableDataFailedPrefix": "Failed to load table, error message: "
    }
  },
  "zh-CN": {
    "title": "难度表管理",
    "button": {
      "add": "新增",
      "delete": "删除"
    },
    "modal": {
      "title": "新增难度表",
      "positiveText": "新增",
      "negativeText": "取消",
      "labelAddress": "地址",
      "placeholderAddress": "请输入地址"
    },
    "rules": {
      "missingAddress": "请输入地址"
    },
    "column": {
      "name": "名称",
      "url": "地址",
      "actions": "操作"
    },
    "deleteDialog": {
      "title": "确定要删除吗？",
      "positiveText": "是",
      "negativeText": "否"
    },
    "message": {
      "addTableFailedPrefix": "新增难度表失败，错误信息: ",
      "deleteSuccess": "删除成功",
      "loadTableDataFailedPrefix": "读取难度表信息失败, 错误信息: "
    }
  }
}</i18n>