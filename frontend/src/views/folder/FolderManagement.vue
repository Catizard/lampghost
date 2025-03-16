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
		<n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false"
			:row-key="(row: dto.FolderDto) => row.ID" />
	</perfect-scrollbar>

	<n-modal v-model:show="showAddModal" preset="dialog" :title="t('modal.title')"
		:positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
		<n-form ref="formRef" :model="formData" :rules="rules">
			<n-form-item :label="t('form.labelName')" path="name">
				<n-input v-model:value="formData.name" :placeholder="t('form.placeholderName')" />
			</n-form-item>
		</n-form>
	</n-modal>
</template>

<script setup lang="ts">
import { h, Ref, ref, watch } from "vue";
import { dto } from "@wailsjs/go/models";
import {
	DataTableColumns,
	FormInst,
	NButton,
	NDataTable,
	useDialog,
	useNotification,
} from "naive-ui";
import {
	AddFolder,
	DelFolder,
	DelFolderContent,
	FindFolderTree,
	QueryFolderDefinition,
	SyncSongData,
} from "@wailsjs/go/controller/FolderController";
import { useI18n } from "vue-i18n";

const i18n = useI18n();
const { t } = i18n;
const notification = useNotification();
const dialog = useDialog();

const showAddModal = ref(false);
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

const formRef = ref<FormInst | null>(null);
const formData = ref({
	name: "",
});
const rules = {
	name: {
		required: true,
		message: t('message.missingName'),
		trigger: ["input", "blur"],
	},
};

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
			title: t('contentColumn.actions'),
			key: "actions",
			render(row) {
				return h(
					NButton,
					{
						strong: true,
						tertiary: true,
						size: "small",
						onClick: () => deleteFolderContent(row),
					},
					{ default: () => t('button.delete') },
				);
			},
		},
	];
}

function addFolder(name: string) {
	AddFolder(name)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			loadFolderData();
		})
		.catch((err) => {
			notifyError(t('message.addFolderFailedPrefix') + err);
			loadFolderData();
		});
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
	FindFolderTree()
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			data.value = [...result.Rows];
		})
		.catch((err) => {
			notifyError(t('message.loadFolderDataFailedPrefix') + err);
		});
}

function handlePositiveClick(): boolean {
	formRef.value
		?.validate()
		.then(() => {
			addFolder(formData.value.name);
			showAddModal.value = false;
			formData.value.name = null;
		})
		.catch((err) => { });
	return false;
}

function handleNegativeClick() {
	formData.value.name = "";
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
		"modal": {
			"title": "New folder",
			"positiveText": "Submit",
			"negativeText": "Cancel"
		},
		"form": {
			"labelName": "Name",
			"placeholderName": "Input name"
		},
		"column": {
			"name": "Folder Name",
			"actions": "Actions"
		},
		"contentColumn": {
			"name": "Title",
			"actions": "Actions"
		},
		"message": {
			"addFolderFailedPrefix": "Failed to add folder, error message: ",
			"deleteSuccess": "Delete successfully",
			"confirmToDelete": "Do you really want to delete this content?",
			"missingName": "Please input name",
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
		"modal": {
			"title": "新增收藏夹",
			"positiveText": "提交",
			"negativeText": "取消"
		},
		"form": {
			"labelName": "名称",
			"placeholderName": "请输入名称"
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
			"addFolderFailedPrefix": "新增收藏夹失败，错误信息: ",
			"deleteSuccess": "删除成功",
			"confirmToDelete": "确定删除吗？",
			"missingName": "请输入名称",
			"loadFolderDataFailedPrefix": "读取收藏夹信息失败，错误信息: "
		},
		"dialog": {
			"positiveText": "确定",
			"negativeText": "取消"
		}
	}
}</i18n>