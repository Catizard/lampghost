<template>
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

	<FolderAddForm v-model:show="showAddModal" @refresh="loadFolderData" />
</template>

<script setup lang="ts">
import { h, Ref, ref } from "vue";
import { dto } from "@wailsjs/go/models";
import {
	DataTableColumns,
	NButton,
	NDataTable,
	useDialog,
} from "naive-ui";
import { DelFolder, FindFolderList, FindFolderTree } from "@wailsjs/go/main/App";
import { useI18n } from "vue-i18n";
import FolderAddForm from "./FolderAddForm.vue";
import FolderDetail from "./FolderDetail.vue";

const i18n = useI18n();
const { t } = i18n;
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

loadFolderData();

function createColumns({
	deleteFolder,
}: {
	deleteFolder: (row: dto.FolderDto) => void;
}): DataTableColumns<dto.FolderDto> {
	return [
		{
			type: "expand",
			renderExpand: (row: dto.FolderDto) => {
				// TODO: remove rivalId = 1 limitation
				return h(FolderDetail, { folderId: row.ID });
			},
		},
		{ title: t('column.name'), key: "FolderName", minWidth: "150px" },
		{
			title: t('column.actions'),
			key: "actions",
			render(row: dto.FolderDto) {
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

function deleteFolder(id: number) {
	DelFolder(id)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			window.$notifySuccess(t('message.deleteSuccess'));
			loadFolderData();
		})
		.catch((err) => {
			window.$notifyError(err);
			loadFolderData();
		});
}

function loadFolderData() {
	// TODO: remove magical 1
	loading.value = true;
	FindFolderList({ RivalID: 1 } as any)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			data.value = [...result.Rows];
		})
		.catch((err) => window.$notifyError(t('message.loadFolderDataFailed', { msg: err })))
		.finally(() => loading.value = false);
}
</script>

<i18n lang="json">{
	"en": {
		"title": "Folder",
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
			"clear": "Clear",
			"tag": "Tag"
		},
		"message": {
			"deleteSuccess": "Delete successfully",
			"confirmToDelete": "Do you really want to delete this content?",
			"loadFolderDataFailed": "Failed to load folder data, error message: {msg}"
		},
		"dialog": {
			"positiveText": "Yes",
			"negativeText": "Cancel"
		}
	},
	"zh-CN": {
		"title": "收藏夹",
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
			"actions": "操作",
			"clear": "通关状态",
			"tag": "难度表标签"
		},
		"message": {
			"deleteSuccess": "删除成功",
			"confirmToDelete": "确定删除吗？",
			"loadFolderDataFailed": "读取收藏夹信息失败，错误信息: {msg}"
		},
		"dialog": {
			"positiveText": "确定",
			"negativeText": "取消"
		}
	}
}</i18n>
