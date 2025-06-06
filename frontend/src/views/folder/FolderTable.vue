<template>
	<n-data-table :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
		:row-key="(row: dto.FolderDto) => row.ID" />
</template>

<script setup lang="ts">
import { h, Ref, ref, watch } from "vue";
import { dto } from "@wailsjs/go/models";
import {
	DataTableColumns,
	NButton,
	NDataTable,
	useDialog,
} from "naive-ui";
import { useI18n } from "vue-i18n";
import FolderDetail from "./FolderDetail.vue";
import { DelFolder, FindFolderList } from "@wailsjs/go/main/App";

const props = defineProps<{
	customTableId?: number
	type: "table" | "folder"
}>();
defineExpose({ loadData })

const i18n = useI18n();
const { t } = i18n;
const dialog = useDialog();

const loading = ref(false);
const pagination = false as const;
let data: Ref<Array<any>> = ref([]);
const columns: DataTableColumns<dto.FolderDto> = createColumns({
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

function createColumns({
	deleteFolder,
}: {
	deleteFolder: (row: dto.FolderDto) => void;
}): DataTableColumns<dto.FolderDto> {
	return [
		{
			type: "expand",
			renderExpand: (row: dto.FolderDto) => {
				return h(FolderDetail, { folderId: row.ID, type: props.type });
			},
		},
		{
			title: (): string => {
				if (props.type == "table") {
					return (t('column.name.table'));
				} else if (props.type == "folder") {
					return t('column.name.folder');
				}
			},
			key: "FolderName",
			minWidth: "150px"
		},
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
			loadData();
		})
		.catch((err) => {
			window.$notifyError(err);
			loadData();
		});
}

function loadData() {
	loading.value = true;
	FindFolderList({
		RivalID: 1,
		CustomTableID: props.customTableId ?? 1
	} as any)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			data.value = [...result.Rows];
		})
		.catch((err) => window.$notifyError(t('message.loadFolderDataFailed', { msg: err })))
		.finally(() => loading.value = false);
}

watch(() => props.customTableId, () => {
	loadData();
}, { immediate: true });
</script>

<i18n lang="json">{
	"en": {
		"title": "Folder",
		"button": {
			"addFolder": "Add Folder",
			"delete": "Delete"
		},
		"column": {
			"name": {
				"folder": "Folder Name",
				"table": "Difficult Name"
			},
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
			"name": {
				"folder": "收藏夹名称",
				"table": "难度名称"
			},
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
