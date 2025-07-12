<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.customTableManagement') }}</n-text>
    </n-h1>
    <n-flex justify="end">
      <n-button :loading="loading" type="primary" @click="showAddModal = true">
        {{ t('button.addCustomTable') }}
      </n-button>
    </n-flex>
  </n-flex>

  <n-data-table remote :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
    :row-key="(row: dto.CustomDiffTableDto) => row.ID" />

  <CustomTableAddForm v-model:show="showAddModal" @refresh="loadData()" />
  <CustomTableEditForm ref="editFormRef" @refresh="loadData()" />
</template>

<script lang="ts" setup>
import { DataTableColumns, NButton, NDropdown, useDialog } from 'naive-ui';
import { h, reactive, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { dto } from '@wailsjs/go/models';
import { DeleteCustomDiffTable, QueryCustomDiffTablePageList } from '@wailsjs/go/main/App';
import CustomTableEditForm from './CustomTableEditForm.vue';
import CustomTableAddForm from './CustomTableAddForm.vue';
import { ClipboardSetText } from '@wailsjs/runtime/runtime';

const { t } = useI18n();
const dialog = useDialog();
const editFormRef = ref<InstanceType<typeof CustomTableEditForm>>(null);

const loading = ref(false);
const showAddModal = ref(false);

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
  }
});
const columns: DataTableColumns<dto.CustomDiffTableDto> = [
  { title: t('column.name'), key: "Name" },
  {
    title: t('column.actions'), key: "actions", width: "200px",
    render(row: dto.CustomDiffTableDto) {
      return h(
        NDropdown,
        {
          trigger: "hover",
          options: [
            { label: t("button.edit"), key: "Edit" },
            { label: t("button.delete"), key: "Delete" },
            { label: t("button.link"), key: "Link" }
          ],
          onSelect: (key: string) => {
            switch (key) {
              case "Edit": editFormRef.value.open(row.ID); break;
              case "Delete":
                dialog.warning({
                  title: t('deleteDialog.title'),
                  positiveText: t('deleteDialog.positiveText'),
                  negativeText: t('deleteDialog.negativeText'),
                  onPositiveClick: () => {
                    loading.value = true;
                    DeleteCustomDiffTable(row.ID)
                      .then(result => {
                        if (result.Code != 200) {
                          return Promise.reject(result.Msg);
                        }
                        loadData();
                      })
                      .catch(err => window.$notifyError(err))
                      .finally(() => loading.value = true);
                  }
                });
                break;
              case "Link":
                try {
                  const name = row.Name;
                  ClipboardSetText(`http://localhost:7391/table/${encodeURIComponent(name)}.json`)
                  window.$notifySuccess(t('message.setClipboardSuccess'));
                } catch (e) {
                  window.$notifyError(t('message.setClipboardError', { msg: String(e) }));
                }
                break;
            }
          }
        },
        { default: () => h(NButton, null, { default: () => '...' }) }
      );
    }
  }
];
const data: Ref<dto.CustomDiffTableDto[]> = ref([]);

function loadData() {
  loading.value = true;
  QueryCustomDiffTablePageList({
    Pagination: pagination,
    IgnoreDefaultTable: true
  } as any).then(result => {
    if (result.Code != 200) {
      return Promise.reject(result.Msg);
    }
    data.value = [...result.Rows];
    pagination.pageCount = result.Pagination.pageCount;
  })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false)
}

loadData();
</script>
