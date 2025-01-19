<template>
  <perfect-scrollbar>
    <n-space justify="space-between">
      <n-h1 prefix="bar" style="text-align: start">
        <n-text type="primary"> 是的，这是难度表信息!! </n-text>
      </n-h1>
      <n-button @click="showAddModal = true"> 新增一个难度表 </n-button>
    </n-space>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" />
  </perfect-scrollbar>

  <n-modal v-model:show="showAddModal" preset="dialog" title="新增难度表信息" positive-text="新增" negative-text="取消"
    @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
    <n-form ref="formRef" :model="formData" :rules="rules">
      <n-form-item label="地址" path="url">
        <n-input v-model:value="formData.url" placeholder="输入地址" />
      </n-form-item>
    </n-form>
  </n-modal>

  <n-modal v-model:show="showDetailModal">
    <n-card style="width: 80%" title="详情" :bordered="false" size="huge" role="dialog" aria-modal="true" closable
      @close="() => { showDetailModal = false }">
      <difficult-table-detail :header-id="currentShowHeaderId" :level="currentShowLevel" />
    </n-card>
  </n-modal>
</template>

<script lang="ts" setup>
import type { DataTableColumns, FormInst } from "naive-ui";
import { NButton, NDataTable, useDialog, useMessage } from "naive-ui";
import { useNotification } from "naive-ui";
import { computed, defineComponent, h, Ref, ref } from "vue";
import {
  AddDiffTableHeader,
  DelDiffTableHeader,
  FindDiffTableHeaderTree
} from "@wailsjs/go/controller/DiffTableController";
import { dto, entity } from "@wailsjs/go/models";
import DifficultTableDetail from "./DifficultTableDetail.vue";

const showAddModal = ref(false);
const showDetailModal = ref(false);

const currentShowHeaderId: Ref<number> = ref(null);
const currentShowLevel: Ref<string> = ref(null);

const formRef = ref<FormInst | null>(null);
const formData = ref({
  url: "",
});
const rules = {
  url: {
    required: true,
    message: "请输入地址",
    trigger: ["input", "blur"],
  },
};

const notification = useNotification();
const dialog = useDialog();
loadDiffTableData();

function createColumns({
  deleteHeader,
}: {
  deleteHeader: (row: dto.DiffTableHeaderDto) => void;
}): DataTableColumns<dto.DiffTableHeaderDto> {
  return [
    {
      type: "expand",
      renderExpand: (rowData) => {
        return h(
          NDataTable,
          {
            columns: songDataColumns,
            data: rowData.Children,
            pagination: pagination,
            bordered: false
          }
        )
      }
    },
    { title: "Name", key: "Name", },
    { title: "url", key: "HeaderUrl", },
    {
      title: "Action",
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
          { default: () => "删除" },
        );
      },
    },
  ];
}

function createSongDataColumns(): DataTableColumns<dto.DiffTableDataDto> {
  return [
    { title: "Level Name", key: "Name", maxWidth: "200px", resizable: true },
    {
      title: "Action",
      key: "actions",
      render(row) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: "small",
            onClick: () => showLevelContent(row),
          },
          { default: () => "详情" }
        )
      }
    }
  ];
}
const songDataColumns = createSongDataColumns();

let data: Ref<Array<any>> = ref([]);
const pagination = false as const;
const columns = createColumns({
  deleteHeader(row: any) {
    dialog.warning({
      title: "确定要删除么?",
      positiveText: "确定",
      negativeText: "取消",
      onPositiveClick: () => {
        delDiffTableHeader(row.ID)
      }
    })
  }
});

function addDiffTableHeader(url: string) {
  AddDiffTableHeader(url)
    .then((result) => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      notifySuccess("新增难度表似乎成功了，后台返回结果: " + result.Msg);
      loadDiffTableData();
    })
    .catch((err) => {
      notifyError("新增难度表似乎失败了，后台返回错误: " + err)
      loadDiffTableData();
    });
}

function delDiffTableHeader(id: number) {
  DelDiffTableHeader(id)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      notifySuccess("删除成功")
      loadDiffTableData();
    }).catch(err => {
      notifyError(err)
      loadDiffTableData();
    })
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
  FindDiffTableHeaderTree()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      data.value = [...result.Rows].map(row => {
        return {
          key: row.ID,
          ...row
        }
      });
      console.log(data.value)
    })
    .catch((err) => {
      notification.error({
        content: "读取难度表信息出错:" + err,
        duration: 5000,
        keepAliveOnHover: true
      })
    });
}

function showLevelContent({ ID, Level }) {
  currentShowHeaderId.value = ID;
  currentShowLevel.value = Level;
  console.log(ID)
  console.log(Level)
  showDetailModal.value = true;
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
