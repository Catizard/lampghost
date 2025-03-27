<template>
	<perfect-scrollbar>
		<n-flex justify="space-between">
			<n-h1 prefix="bar" style="text-align: left">
				<n-text type="primary">{{ t('title') }}</n-text>
			</n-h1>
			<n-flex justify="end">
				<n-button type="primary" @click="showAddModal = true">{{ t('button.addFolder') }}</n-button>
			</n-flex>
		</n-flex>
		<n-data-table :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
			:row-key="(row: dto.FolderDto) => row.ID" />
	</perfect-scrollbar>

	<FolderAddForm v-model:show="showAddModal" @refresh="loadFolderData" />
</template>

<script setup lang="ts">
import { h, Ref, ref, watch } from "vue";
import { dto, vo } from "@wailsjs/go/models";
import {
	DataTableColumns,
	NButton,
	NDataTable,
	useDialog,
	useNotification,
} from "naive-ui";
import {
	DelFolder,
	DelFolderContent,
	FindFolderTree,
} from "@wailsjs/go/controller/FolderController";
import { useI18n } from "vue-i18n";
import FolderAddForm from "./FolderAddForm.vue";
import ClearTag from "@/components/ClearTag.vue";

const i18n = useI18n();
const { t } = i18n;
const notification = useNotification();
const dialog = useDialog();

const showAddModal = ref(false);
const loading = ref(false);
const pagination = false as const;
let data: Ref<Array<any>> = ref([]);
const columns = createColumns({
	deleteFolder(row: any) {
		dialog.warning({
			title: t('message.confirmToDelete'),
			positiveText: t('dialog.positiveText'),
			negativeText: t('dialog.negativeText'),
			onPositiveClick: () => {
				deleteFolder(row.ID);
			},
		});
	},
});
const folderContentColumns = createContentColumns({
	deleteFolderContent(row: any) {
		dialog.warning({
			title: t('message.confirmToDelete'),
			positiveText: t('dialog.positiveText'),
			negativeText: t('dialog.negativeText'),
			onPositiveClick: () => {
				deleteFolderContent(row.ID);
			},
		});
	},
});

loadFolderData();

function createColumns({
	deleteFolder,
}: {
	deleteFolder: (row: dto.FolderDto) => void;
}): DataTableColumns<dto.FolderDto> {
	return [
		{
			type: "expand",
			renderExpand: (rowData) => {
				return h(NDataTable, {
					columns: folderContentColumns,
					data: rowData.Contents,
					pagination: pagination,
					bordered: false,
				});
			},
		},
		{ title: t('column.name'), key: "FolderName" },
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
						type: "error",
						onClick: () => deleteFolder(row),
					},
					{ default: () => t('button.delete') },
				);
			},
		},
	];
}

function createContentColumns({
	deleteFolderContent,
}: {
	deleteFolderContent: (row: dto.FolderContentDto) => void;
}): DataTableColumns<dto.FolderContentDto> {
	return [
		{ title: t('contentColumn.name'), key: "Title" },
		{
			title: t('contentColumn.clear'), key: "Clear",
			render: (row: dto.FolderContentDto) => {
				return h(ClearTag, { clear: row.Lamp },)
			}
		},
		{
			title: t('contentColumn.actions'),
			key: "actions",
			render(row) {
				return h(
					NButton,
					{
						strong: true,
						tertiary: true,
						size: "small",
						type: "error",
						onClick: () => deleteFolderContent(row),
					},
					{ default: () => t('button.delete') },
				);
			},
		},
	];
}

function deleteFolder(id: number) {
	DelFolder(id)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			notifySuccess(t('message.deleteSuccess'));
			loadFolderData();
		})
		.catch((err) => {
			notifyError(err);
			loadFolderData();
		});
}

function deleteFolderContent(id: number) {
	DelFolderContent(id)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			notifySuccess(t('message.deleteSuccess'));
			loadFolderData();
		})
		.catch((err) => {
			notifyError(err);
			loadFolderData();
		});
}

function loadFolderData() {
	// TODO: remove magical 1
	loading.value = true;
	FindFolderTree({ RivalID: 1 } as any)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			data.value = [...result.Rows];
		}).catch((err) => {
			notifyError(t('message.loadFolderDataFailedPrefix') + err);
		}).finally(() => loading.value = false);
}

function notifySuccess(msg: string) {
	notification.success({
		content: msg,
		duration: 5000,
		keepAliveOnHover: true,
	});
}

function notifyError(msg: string) {
	notification.error({
		content: msg,
		duration: 5000,
		keepAliveOnHover: true,
	});
}
</script>

<i18n lang="json">{
	"en": {
		"title": "Folder Management",
		"button": {
			"addFolder": "Add Folder",
			"delete": "Delete"
		},
		"column": {
			"name": "Folder Name",
			"actions": "Actions"
		},
		"contentColumn": {
			"name": "Title",
			"actions": "Actions",
			"clear": "Clear"
		},
		"message": {
			"deleteSuccess": "Delete successfully",
			"confirmToDelete": "Do you really want to delete this content?",
			"loadFolderDataFailedPrefix": "Failed to load folder data, error message: "
		},
		"dialog": {
			"positiveText": "Yes",
			"negativeText": "Cancel"
		}
	},
	"zh-CN": {
		"title": "收藏夹管理",
		"button": {
			"addFolder": "新增收藏夹",
			"delete": "删除"
		},
		"column": {
			"name": "收藏夹名称",
			"actions": "操作"
		},
		"contentColumn": {
			"name": "谱名",
			"actions": "操作"
		},
		"message": {
			"deleteSuccess": "删除成功",
			"confirmToDelete": "确定删除吗？",
			"loadFolderDataFailedPrefix": "读取收藏夹信息失败，错误信息: "
		},
		"dialog": {
			"positiveText": "确定",
			"negativeText": "取消"
		}
	}
}</i18n>