<template>
	<n-flex justify="space-between">
		<n-h1 prefix="bar" style="text-align: start;">
			<n-text type="primary">{{ t('title.playerManagement') }}</n-text>
		</n-h1>
		<!-- <n-flex justify="end">
				<n-button>Reload All</n-button>
			</n-flex> -->
		<n-button type="primary" @click="showAddModal = true">{{ t('button.addPlayer') }}</n-button>
	</n-flex>
	<n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" :row-key="row => row.ID"
		:loading="loading" />

	<RivalAddForm v-model:show="showAddModal" @refresh="loadData()" />
	<RivalEditForm ref="editFormRef" @refresh="loadData()" />
</template>

<script lang="ts" setup>
import router from '@/router';
import { DelRivalInfo, QueryRivalInfoPageList, ReloadRivalData } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import dayjs from 'dayjs';
import { DataTableColumns, NButton, NFlex, useDialog } from 'naive-ui';
import { h, reactive, ref, VNode } from 'vue';
import { useI18n } from 'vue-i18n';
import RivalAddForm from './RivalAddForm.vue';
import RivalEditForm from './RivalEditForm.vue';

const showAddModal = ref(false);
const editFormRef = ref<InstanceType<typeof RivalEditForm>>(null);

const i18n = useI18n();
const { t } = i18n;
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
			title: t('column.name'), key: "Name", ellipsis: { tooltip: true },
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
		{ title: t('column.playCount'), key: "PlayCount", width: "100px", ellipsis: { tooltip: true } },
		{ title: t('column.scoreLogFilePath'), key: "ScoreLogPath", width: "75px", ellipsis: { tooltip: true } },
		{ title: t('column.songdataFilePath'), key: "SongDataPath", width: "75px", ellipsis: { tooltip: true } },
		{
			title: t('column.lastSyncTime'),
			key: "UpdateAt",
			width: "125px",
			render: (row: dto.RivalInfoDto) => {
				return dayjs(row.UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
			}
		},
		{
			title: t('column.actions'),
			key: "action",
			width: "300px",
			render(row: dto.RivalInfoDto): VNode {
				return h(
					NFlex,
					{},
					{
						default: () => {
							return [
								h(NButton, { type: 'primary', size: 'small', onClick: () => handleSyncClick(row.ID) }, { default: () => t('button.reload') }),
								h(NButton, { type: 'info', size: 'small', onClick: () => handleSelectOtherAction(row, "Edit") }, { default: () => t('button.edit') }),
								h(NButton, { type: 'error', size: 'small', onClick: () => handleSelectOtherAction(row, "Delete") }, { default: () => t('button.delete') }),
							]
						}
					}
				);
			}
		}
	]
}

function handleSelectOtherAction(row: dto.RivalInfoDto, key: string) {
	if ("Delete" === key) {
		if (row.MainUser) {
			window.$notifyError(t('message.cannotDeleteMainUser'));
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
					})
					.catch(err => window.$notifyError(err))
					.finally(() => {
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
		})
		.catch(err => window.$notifyError(err))
		.finally(() => loading.value = false);
}

function handleSyncClick(id: number) {
	loading.value = true;
	ReloadRivalData(id, true)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg)
			}
			window.$notifySuccess(t('message.reloadSuccess'));
		})
		.catch(err => window.$notifyError(err))
		.finally(() => loading.value = false);
}

loadData();
</script>
