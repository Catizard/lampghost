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

	<RivalAddForm v-model:show="showAddModal" @refresh="loadData()" />
	<RivalEditForm ref="editFormRef" @refresh="loadData()" />
</template>

<script lang="ts" setup>
import router from '@/router';
import { AddRivalInfo, DelRivalInfo, QueryRivalInfoPageList, ReloadRivalData } from '@wailsjs/go/main/App';
import { dto, entity } from '@wailsjs/go/models';
import dayjs from 'dayjs';
import { DataTableColumns, DropdownOption, FormInst, NAnchorLink, NButton, NDropdown, useDialog, useNotification } from 'naive-ui';
import { h, reactive, ref, VNode } from 'vue';
import { useI18n } from 'vue-i18n';
import RivalAddForm from './RivalAddForm.vue';
import RivalEditForm from './RivalEditForm.vue';

const showAddModal = ref(false);
const editFormRef = ref<InstanceType<typeof RivalEditForm>>(null);

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
				const reloadBubutton = h(
					NButton,
					{
						strong: true,
						tertiary: true,
						type: 'primary',
						size: "small",
						style: "margin-right: 5px",
						onClick: () => { handleSyncClick(row.ID) }
					},
					{ default: () => t('button.reload') }
				);
				const otherActions: VNode = h(
					NDropdown,
					{
						trigger: "hover",
						options: otherActionOptions,
						size: "small",
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
				return [reloadBubutton, otherActions];
			}
		}
	]
}

const otherActionOptions: Array<DropdownOption> = [
	{
		label: t('button.edit'),
		key: "Edit",
	},
	{
		label: t('button.delete'),
		key: "Delete",
		props: {
			style: "color: red"
		}
	}
];
function handleSelectOtherAction(row: dto.RivalInfoDto, key: string) {
	if ("Delete" === key) {
		if (row.MainUser) {
			notification.error({
				content: t('message.cannotDeleteMainUser'),
				duration: 3000,
				keepAliveOnHover: true
			});
			return;
		}
		dialog.warning({
			title: t('deleteDialog.title'),
			positiveText: t('deleteDialog.positiveText'),
			negativeText: t('deleteDialog.negativeText'),
			onPositiveClick: () => {
				loading.value = true;
				DelRivalInfo(row.ID)
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
	if ("Edit" === key) {
		editFormRef.value.open(row.ID);
	}
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
	ReloadRivalData(id, true)
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

loadData();
</script>

<i18n lang="json">{
	"en": {
		"title": "Player Management",
		"column": {
			"name": "Name",
			"count": "Play Count",
			"scoreLogFilePath": "scorelog.db File Path",
			"songdataFilePath": "songdata.db File Path",
			"lastSyncTime": "Last Sync Time",
			"actions": "Actions"
		},
		"button": {
			"reload": "Fully Reload",
			"add": "Add Player",
			"delete": "Delete",
			"edit": "Edit"
		},
		"message": {
			"reloadSuccess": "Reload successfully",
			"cannotDeleteMainUser": "Cannot delete main user"
		},
		"deleteDialog": {
			"title": "Confirm to delete?",
			"positiveText": "Yes",
			"negativeText": "No"
		}
	},
	"zh-CN": {
		"title": "玩家管理",
		"column": {
			"name": "名称",
			"count": "游玩次数",
			"scoreLogFilePath": "scorelog.db文件路径",
			"songdataFilePath": "songdata.db文件路径",
			"lastSyncTime": "最后更新时间",
			"actions": "操作"
		},
		"button": {
			"reload": "全量同步",
			"add": "添加玩家",
			"delete": "删除",
			"edit": "修改"
		},
		"message": {
			"reloadSuccess": "同步成功",
			"cannotDeleteMainUser": "不能删除主用户"
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
		"deleteDialog": {
			"title": "确定要删除吗？",
			"positiveText": "是",
			"negativeText": "否"
		}
	}
}</i18n>
