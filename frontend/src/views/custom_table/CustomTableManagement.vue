<template>
	<n-flex justify="space-between">
		<n-h1 prefix="bar" style="text-align: start;">
			<n-text type="primary">{{ t('title') }}</n-text>
		</n-h1>
		<n-flex justify="end">
			<n-button :loading="loading" type="primary" @click="showAddModal = true">
				{{ t('button.add') }}
			</n-button>
		</n-flex>
	</n-flex>

	<n-data-table remote :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
		:row-key="(row: dto.CustomDiffTableDto) => row.ID" />

	<CustomTableAddForm v-model:show="showAddModal" @refresh="loadData()" />
	<CustomTableEditForm ref="editFormRef" @refresh="loadData()" />
</template>

<script lang="ts" setup>
import { DataTableColumns, NButton, useDialog } from 'naive-ui';
import { h, reactive, Ref, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { dto } from '@wailsjs/go/models';
import { DeleteCustomDiffTable, QueryCustomDiffTablePageList } from '@wailsjs/go/main/App';
import CustomTableEditForm from './CustomTableEditForm.vue';
import CustomTableAddForm from './CustomTableAddForm.vue';

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
			return [
				h(NButton, { type: "primary", size: "small", onClick: () => editFormRef.value.open(row.ID) }, { default: () => t('button.update') }),
				h(NButton, {
					style: {
						"margin-left": "5px",
					},
					type: "error", size: "small", onClick: () => {
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
						})
					}
				}, { default: () => t('button.delete') })
			];
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

<i18n lang="json">{
	"en": {
		"title": "Custom Table Management",
		"column": {
			"name": "Name",
			"actions": "Actions"
		},
		"button": {
			"add": "Add Custom Table",
			"update": "Edit",
			"delete": "Delete"
		},
		"deleteDialog": {
			"title": "Delete Custom Table",
			"positiveText": "Yes",
			"negativeText": "Cancel"
		}
	},
	"zh-CN": {
		"title": "自定义难度表管理",
		"column": {
			"name": "名称",
			"actions": "操作"
		},
		"button": {
			"add": "新增自定义难度表",
			"update": "修改",
			"delete": "删除"
		},
		"deleteDialog": {
			"title": "删除自定义难度表",
			"positiveText": "确认",
			"negativeText": "取消"
		}
	}
}</i18n>
