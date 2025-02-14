<template>
	<perfect-scrollbar>
		<n-flex justify="space-between">
			<n-h1 prefix="bar" style="text-align: left">
				<n-text type="primary">收藏夹管理</n-text>
			</n-h1>
			<n-flex justify="end">
				<n-button type="primary" @click="showAddModal = true">新增收藏夹</n-button>
				<n-button type="primary" @click="showJsonModal = true">生成JSON</n-button>
				<n-button type="error" @click="handleClickSync">同步数据</n-button>
			</n-flex>
		</n-flex>
		<n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" />
	</perfect-scrollbar>

	<n-modal v-model:show="showAddModal" preset="dialog" title="新增收藏夹" positive-text="新增" negative-text="取消"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
		<n-form ref="formRef" :model="formData" :rules="rules">
			<n-form-item label="名称" path="name">
				<n-input v-model:value="formData.name" placeholder="输入名称" />
			</n-form-item>
		</n-form>
	</n-modal>

	<n-modal v-model:show="showJsonModal">
		<n-card style="width: 80%" :boarded="false" size="huge" role="dialog" aria-modal="true" closable
			@close="() => { showJsonModal = false; }" :loading="jsonGenerateLoading">
			<pre>{{ folderJsonCode }}</pre>
		</n-card>
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

const notification = useNotification();
const dialog = useDialog();

const showAddModal = ref(false);
const showJsonModal = ref(false);
const jsonGenerateLoading = ref(false);
const pagination = false as const;
let data: Ref<Array<any>> = ref([]);
const columns = createColumns({
	deleteFolder(row: any) {
		dialog.warning({
			title: "确定要删除么?",
			positiveText: "确定",
			negativeText: "取消",
			onPositiveClick: () => {
				deleteFolder(row.ID);
			},
		});
	},
});
const folderContentColumns = createContentColumns({
	deleteFolderContent(row: any) {
		dialog.warning({
			title: "确定要删除么?",
			positiveText: "确定",
			negativeText: "取消",
			onPositiveClick: () => {
				deleteFolderContent(row.ID);
			},
		});
	},
});
const folderJsonCode = ref("");

const formRef = ref<FormInst | null>(null);
const formData = ref({
	name: "",
});
const rules = {
	name: {
		required: true,
		message: "请输入名称",
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
		{ title: "Folder Name", key: "FolderName" },
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
						onClick: () => deleteFolder(row),
					},
					{ default: () => "删除" },
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
		{ title: "Name", key: "Title" },
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
						onClick: () => deleteFolderContent(row),
					},
					{ default: () => "删除" },
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
			notifySuccess("新增收藏夹成功: " + result.Msg);
			loadFolderData();
		})
		.catch((err) => {
			notifyError("新增收藏家失败: " + err);
			loadFolderData();
		});
}

function deleteFolder(id: number) {
	DelFolder(id)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			notifySuccess("删除成功");
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
			notifySuccess("删除成功");
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
			data.value = [...result.Rows].map((row) => {
				return {
					key: row.ID,
					...row,
				};
			});
		})
		.catch((err) => {
			notifyError("读取收藏夹数据出错:" + err);
		});
}

function handlePositiveClick(): boolean {
	formRef.value
		?.validate()
		.then(() => {
			addFolder(formData.value.name);
			showAddModal.value = false;
		})
		.catch((err) => { });
	return false;
}

function handleNegativeClick() {
	formData.value.name = "";
}

function generateJson() {
	jsonGenerateLoading.value = true;
	QueryFolderDefinition()
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			jsonGenerateLoading.value = false;
			folderJsonCode.value = JSON.stringify(result.Rows, null, "\t");
		}).catch(err => {
			notifyError("生成json失败: " + err)
			jsonGenerateLoading.value = false;
		})
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

function handleClickSync() {
	dialog.warning({
		title: "警告",
		content: "该操作会直接修改你在本地的songdata.db文件, 移除已有的所有配置(包括你在游戏中自行设置的favorite和invisible)并重新应用当前的收藏夹配置。该操作不可逆, 请在操作前对你的数据进行备份!",
		positiveText: "确定",
		negativeText: "取消",
		onPositiveClick: () => {
			SyncSongData().then(result => {
				if (result.Code != 200) {
					return Promise.reject(result.Msg)
				}
				notifySuccess("更新成功, 后台返回:" + result.Msg)
			}).catch(err => {
				notifyError("更新失败, 后台返回:" + err)
			})
		},
		onNegativeClick: () => {
			// do nothing
		}
	})
}

watch(showJsonModal, (newValue, oldValue) => {
	if (newValue == true) {
		generateJson();
	}
});
</script>
