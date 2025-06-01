<template>
	<n-flex justify="space-between">
		<n-h1 prefix="bar" style="text-align: start">
			<n-text type="primary">{{ t('title') }}</n-text>
		</n-h1>
		<n-flex justify="flex-end">
			<n-button :loading="loading" type="info" @click="showSortModal = true">
				{{ t('button.sort') }}
			</n-button>
			<n-button :loading="loading" type="primary" @click="showAddModal = true">
				{{ t('button.add') }}
			</n-button>
		</n-flex>
	</n-flex>
	<n-data-table :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
		:row-key="(row: dto.DiffTableHeaderDto) => row.ID" />

	<DifficultTableSortModal v-model:show="showSortModal" @refresh="loadDiffTableData()" />
	<DifficultTableAddForm v-model:show="showAddModal" @refresh="loadDiffTableData()" />
	<DifficultTableEditForm ref="editFormRef" @refresh="loadDiffTableData()" />
	<DifficultTableLevelSortModal ref="levelSortModalRef" />
</template>

<script lang="ts" setup>
import type { DropdownOption, TagProps } from "naive-ui";
import { NButton, NDataTable, NDropdown, NTag, useDialog } from "naive-ui";
import { useNotification } from "naive-ui";
import { h, Ref, ref } from "vue";
import {
	DelDiffTableHeader,
	FindDiffTableHeaderTree,
	ReloadDiffTableHeader,
	SupplyMissingBMSFromTable
} from "@wailsjs/go/main/App";
import { dto } from "@wailsjs/go/models";
import { useI18n } from "vue-i18n";
import DifficultTableAddForm from "./DifficultTableAddForm.vue";
import DifficultTableEditForm from "./DifficultTableEditForm.vue";
import DifficultTableSortModal from "./DifficultTableSortModal.vue";
import DifficultTableLevelSortModal from "./DifficultTableLevelSortModal.vue";
import { TagColor } from "naive-ui/es/tag/src/common-props";

const i18n = useI18n();
const { t } = i18n;
const showSortModal = ref(false);
const showAddModal = ref(false);
const editFormRef = ref<InstanceType<typeof DifficultTableEditForm>>(null);
const levelSortModalRef = ref<InstanceType<typeof DifficultTableLevelSortModal>>(null);

const notification = useNotification();
const dialog = useDialog();
const loading = ref(false);
loadDiffTableData();

const columns = [
	{ title: t('column.name'), key: "Name", },
	{
		title: t('column.tag'), key: "Tag",
		render(row: dto.DiffTableHeaderDto) {
			let tagColorProp: TagColor = {};
			if (row.TagColor != '') {
				tagColorProp.color = row.TagColor;
			}
			if (row.TagTextColor != '') {
				tagColorProp.textColor = row.TagTextColor;
			}
			return h(
				NTag,
				{ color: tagColorProp },
				{ default: () => row.Symbol == '' ? '/' : row.Symbol },
			)
		}
	},
	{ title: t('column.url'), key: "HeaderUrl", },
	{
		title: t('column.actions'),
		key: "actions",
		render(row) {
			return h(
				NDropdown,
				{
					trigger: "hover",
					options: otherActionOptions,
					onSelect: (key: string) => handleSelectOtherAction(row, key)
				},
				{
					default: () => h(
						NButton,
						null,
						{ default: () => '...' }
					)
				},
			);
		},
	},
];

let data: Ref<Array<any>> = ref([]);
const pagination = false as const;

const otherActionOptions: Array<DropdownOption> = [
	{ label: t('button.reload'), key: "Reload" },
	{ label: t('button.edit'), key: "Edit", },
	{ label: t('button.supply'), key: "Supply", disabled: true },
	{ label: t('button.sortLevels'), key: "SortLevels", },
	{
		label: t('button.delete'),
		key: "Delete",
		props: {
			style: "color: red"
		}
	}
];

function handleSelectOtherAction(row: dto.DiffTableHeaderDto, key: string) {
	if ("Reload" === key) {
		reloadTableHeader(row.ID);
	}
	if ("Supply" === key) {
		dialog.warning({
			title: t('supplyDialog.title'),
			positiveText: t('supplyDialog.positiveText'),
			negativeText: t('supplyDialog.negativeText'),
			onPositiveClick: () => {
				loading.value = true;
				SupplyMissingBMSFromTable(row.ID)
					.then(result => {
						if (result.Code != 200) {
							return Promise.reject(result.Msg);
						}
					})
					.catch(err => window.$notifyError(err))
					.finally(() => loading.value = false);
			}
		})
	}
	if ("Delete" === key) {
		dialog.warning({
			title: t('deleteDialog.title'),
			positiveText: t('deleteDialog.positiveText'),
			negativeText: t('deleteDialog.negativeText'),
			onPositiveClick: () => {
				delDiffTableHeader(row.ID)
			}
		})
	}
	if ("Edit" === key) {
		editFormRef.value.open(row.ID);
	}
	if ("SortLevels" === key) {
		levelSortModalRef.value.open(row.ID);
	}
}

function reloadTableHeader(id: number) {
	loading.value = true;
	ReloadDiffTableHeader(id)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			loadDiffTableData();
		}).catch(err => {
			notifyError(err);
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
			"delete": "Delete",
			"edit": "Edit",
			"sort": "Sort",
			"sortLevels": "Sort Levels",
			"reload": "Reload",
			"supply": "Supply all missing bms"
		},
		"column": {
			"name": "Name",
			"url": "URL",
			"actions": "Actions",
			"tag": "Tag"
		},
		"deleteDialog": {
			"title": "Confirm to delete?",
			"positiveText": "Yes",
			"negativeText": "No"
		},
		"supplyDialog": {
			"title": "Do you really want to supply all missing bms?",
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
			"delete": "删除",
			"edit": "修改",
			"sort": "排序",
			"sortLevels": "设定难度排序",
			"reload": "重新导入",
			"supply": "补充所有缺少的BMS"
		},
		"column": {
			"name": "名称",
			"url": "地址",
			"actions": "操作",
			"tag": "难度表标签"
		},
		"deleteDialog": {
			"title": "确定要删除吗？",
			"positiveText": "是",
			"negativeText": "否"
		},
		"supplyDialog": {
			"title": "确定要添加所有缺少的BMS吗？",
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
