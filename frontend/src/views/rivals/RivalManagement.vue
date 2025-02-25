<template>
	<perfect-scrollbar>
		<n-flex justify="space-between">
			<n-h1 prefix="bar" style="text-align: start;">
				<n-text type="primary">{{ t('title') }}</n-text>
			</n-h1>
			<!-- <n-flex justify="end">
				<n-button>Reload All</n-button>
			</n-flex> -->
			<n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" :row-key="row => row.ID"
				:loading="loading" />
		</n-flex>
	</perfect-scrollbar>
</template>

<script lang="ts" setup>
import { FindRivalInfoList, SyncRivalScoreLog } from '@wailsjs/go/controller/RivalInfoController';
import { dto } from '@wailsjs/go/models';
import * as dayjs from 'dayjs';
import { DataTableColumns, NButton, useNotification } from 'naive-ui';
import { h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const { t } = i18n;
const notification = useNotification();

const pagination = false as const;
const columns = createColumns();
const data = ref<Array<dto.RivalInfoDto>>([]);
const loading = ref<boolean>(false);

function createColumns(): DataTableColumns<dto.RivalInfoDto> {
	return [
		{ title: t('column.name'), key: "Name" },
		{ title: t('column.count'), key: "PlayCount" },
		{ title: t('column.scoreLogFilePath'), key: "ScoreLogPath", ellipsis: true, maxWidth: "150px" },
		{ title: t('column.songdataFilePath'), key: "SongDataPath", ellipsis: true, maxWidth: "150px" },
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
				return h(
					NButton,
					{
						strong: true,
						tertiary: true,
						onClick: () => { handleSyncClick(row.ID) }
					},
					{ default: () => t('button.reload') }
				)
			}
		}
	]
}

function loadData() {
	loading.value = true;
	FindRivalInfoList()
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

function handleSyncClick(id) {
	loading.value = true;
	SyncRivalScoreLog(id)
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

<i18n>
{
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
		},
		"message": {
			"reloadSuccess": "Reload successfully",
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
		},
		"message": {
			"reloadSuccess": "同步成功",
		}
	},
}
</i18n>