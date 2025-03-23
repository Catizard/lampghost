<template>
	<perfect-scrollbar>
		<n-flex justify="space-between">
			<n-h1 prefix="bar" style="text-align: start;">
				<n-text type="primary">{{ t('title') }}</n-text>
			</n-h1>
			<!-- <n-flex justify="end">
				<n-button>Reload All</n-button>
			</n-flex> -->
			<n-button type="primary" @click="showAddModal = true">{{ t('button.add') }}</n-button>
		</n-flex>
		<n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" :row-key="row => row.ID"
			:loading="loading" />
	</perfect-scrollbar>

	<n-modal v-model:show="showAddModal" preset="dialog" :title="t('modal.title')"
		:positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
		<n-form ref="formRef" :model="formData" :rules="rules">
			<n-form-item :label="t('modal.labelRivalName')" path="Name">
				<n-input v-model:value="formData.Name" :placeholder="t('modal.placeholderRivalName')" />
			</n-form-item>
			<n-form-item :label="t('modal.labelScoreLogPath')" path="ScoreLogPath">
				<n-input v-model:value="formData.ScoreLogPath" :placeholder="t('modal.placeholderScoreLogPath')" />
			</n-form-item>
			<n-form-item :label="t('modal.labelSongDataPath')" path="SongDataPath">
				<n-input disabled v-model:value="formData.SongDataPath" :placeholder="t('modal.placeholderSongDataPath')" />
			</n-form-item>
		</n-form>
	</n-modal>
</template>

<script lang="ts" setup>
import router from '@/router';
import { AddRivalInfo, DelRivalInfo, QueryRivalInfoPageList, SyncRivalDataByID } from '@wailsjs/go/controller/RivalInfoController';
import { dto, entity } from '@wailsjs/go/models';
import dayjs from 'dayjs';
import { DataTableColumns, FormInst, NAnchorLink, NButton, useDialog, useNotification } from 'naive-ui';
import { h, reactive, ref, VNode } from 'vue';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const { t } = i18n;
const notification = useNotification();
const dialog = useDialog();

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
const columns = createColumns();
const data = ref<Array<dto.RivalInfoDto>>([]);
const loading = ref<boolean>(false);

function createColumns(): DataTableColumns<dto.RivalInfoDto> {
	return [
		{
			title: t('column.name'), key: "Name", width: "100px", ellipsis: { tooltip: true },
			render: (row: dto.RivalInfoDto) => {
				return h(
					NButton,
					{
						text: true,
						type: "info",
						onClick: () => router.push({ path: "/home", query: { rivalID: row.ID } })
					},
					{ default: () => row.Name }
				)
			}
		},
		{ title: t('column.count'), key: "PlayCount", width: "100px", ellipsis: { tooltip: true } },
		{ title: t('column.scoreLogFilePath'), key: "ScoreLogPath", maxWidth: "150px", ellipsis: { tooltip: true } },
		{ title: t('column.songdataFilePath'), key: "SongDataPath", maxWidth: "150px", ellipsis: { tooltip: true } },
		{
			title: t('column.lastSyncTime'),
			key: "UpdateAt",
			render: (row: dto.RivalInfoDto) => {
				return dayjs(row.UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
			}
		},
		{
			title: t('column.actions'),
			key: "action",
			render: (row: dto.RivalInfoDto) => {
				let vnodes: VNode[] = [];
				const reloadBubutton = h(
					NButton,
					{
						strong: true,
						tertiary: true,
						type: 'primary',
						onClick: () => { handleSyncClick(row.ID) }
					},
					{ default: () => t('button.reload') }
				);
				vnodes.push(reloadBubutton);
				if (!row.MainUser) {
					const deleteButton = h(
						NButton,
						{
							strong: true,
							tertiary: true,
							type: "error",
							onClick: () => { handleDeleteClick(row.ID) }
						},
						{ default: () => t('button.delete') }
					);
					vnodes.push(deleteButton);
				}
				return vnodes;
			}
		}
	]
}

function loadData() {
	loading.value = true;
	QueryRivalInfoPageList({
		Pagination: pagination,
	} as any)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg)
			}
			data.value = [...result.Rows];
		}).catch(err => {
			notification.error({
				content: err,
				duration: 3000,
				keepAliveOnHover: true
			})
		}).finally(() => {
			loading.value = false;
		});
}

function handleSyncClick(id: number) {
	loading.value = true;
	SyncRivalDataByID(id)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg)
			}
			notification.success({
				content: t('message.reloadSuccess'),
				duration: 3000,
				keepAliveOnHover: false,
			});
		}).catch(err => {
			notification.error({
				content: err,
				duration: 3000,
				keepAliveOnHover: true,
			});
		}).finally(() => {
			loading.value = false;
		})
}

function handleDeleteClick(id: number) {
	dialog.warning({
		title: t('deleteDialog.title'),
		positiveText: t('deleteDialog.positiveText'),
		negativeText: t('deleteDialog.negativeText'),
		onPositiveClick: () => {
			loading.value = true;
			DelRivalInfo(id)
				.then(result => {
					if (result.Code != 200) {
						return Promise.reject(result.Msg);
					}
				}).catch(err => {
					notification.error({
						content: err,
						duration: 3000,
						keepAliveOnHover: true
					});
				}).finally(() => {
					loadData();
					loading.value = false;
				});
		}
	})
}

loadData();

const showAddModal = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = ref({
	Name: null,
	ScoreLogPath: null,
	SongDataPath: null,
});
const rules = {
	Name: {
		required: true,
		message: t('rules.missingRivalName'),
		trigger: ["input", "blur"],
	},
	ScoreLogPath: {
		required: true,
		message: t('rules.missingScoreLogPath'),
		trigger: ["input", "blur"],
	},
	// SongDataPath: {
	// 	required: true,
	// 	message: t('rules.missingSongDataPath'),
	// 	trigger: ["input", "blur"],
	// }
};

function handlePositiveClick(): boolean {
	formRef.value
		?.validate()
		.then(() => {
			return AddRivalInfo(formData.value as any as entity.RivalInfo)
				.then(result => {
					if (result.Code != 200) {
						return Promise.reject(result.Msg);
					}
					showAddModal.value = false;
					loadData();
				});
		})
		.catch((err) => {
			notification.error({
				content: err,
				duration: 3000,
				keepAliveOnHover: true
			});
		});
	return false;
}

function handleNegativeClick() {
	formData.value.Name = null;
	formData.value.ScoreLogPath = null;
	formData.value.SongDataPath = null;
}
</script>

<i18n lang="json">{
	"en": {
		"title": "Rivals",
		"column": {
			"name": "Name",
			"count": "Play Count",
			"scoreLogFilePath": "scorelog.db File Path",
			"songdataFilePath": "songdata.db File Path",
			"lastSyncTime": "Last Sync Time",
			"actions": "Actions"
		},
		"button": {
			"reload": "Reload",
			"add": "Add Rival",
			"delete": "Delete"
		},
		"message": {
			"reloadSuccess": "Reload successfully"
		},
		"modal": {
			"title": "Add Rival",
			"positiveText": "Submit",
			"negativeText": "Cancel",
			"labelRivalName": "Name",
			"labelScoreLogPath": "scorelog.db file path",
			"labelSongDataPath": "songdata.db file path",
			"placeholderRivalName": "Please input rival's name",
			"placeholderScoreLogPath": "Please input scorelog.db file path",
			"placeholderSongDataPath": "Please input songdata.db file path"
		},
		"rules": {
			"missingRivalName": "Rival's name cannot be empty",
			"missingScoreLogPath": "scorelog.db file path cannot be empty",
			"missingSongDataPath": "songdata.db file path cannot be empty"
		},
		"deleteDialog": {
			"title": "Confirm to delete?",
			"positiveText": "Yes",
			"negativeText": "No"
		}
	},
	"zh-CN": {
		"title": "好友管理",
		"column": {
			"name": "名称",
			"count": "游玩次数",
			"scoreLogFilePath": "scorelog.db文件路径",
			"songdataFilePath": "songdata.db文件路径",
			"lastSyncTime": "最后更新时间",
			"actions": "操作"
		},
		"button": {
			"reload": "同步",
			"add": "添加好友",
			"delete": "删除"
		},
		"message": {
			"reloadSuccess": "同步成功"
		},
		"modal": {
			"title": "新增好友",
			"positiveText": "提交",
			"negativeText": "取消",
			"labelRivalName": "名称",
			"labelScoreLogPath": "scorelog.db文件路径",
			"labelSongDataPath": "songdata.db文件路径",
			"placeholderRivalName": "请输入好友名称",
			"placeholderScoreLogPath": "请输入scorelog.db文件路径",
			"placeholderSongDataPath": "请输入songdata.db文件路径"
		},
		"rules": {
			"missingRivalName": "好友名称不可为空",
			"missingScoreLogPath": "scorelog.db文件路径不可为空",
			"missingSongDataPath": "songdata.db文件路径不可为空"
		},
		"deleteDialog": {
			"title": "确定要删除吗？",
			"positiveText": "是",
			"negativeText": "否"
		}
	}
}</i18n>