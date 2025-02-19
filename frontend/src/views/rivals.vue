<template>
	<perfect-scrollbar>
		<n-flex justify="space-between">
			<n-h1 prefix="bar" style="text-align: start;">
				<n-text type="primary">Rivals</n-text>
			</n-h1>
			<n-flex justify="end">
				<n-button>Reload All</n-button>
			</n-flex>

			<n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" :row-key="row => row.ID" />
		</n-flex>
	</perfect-scrollbar>
</template>

<script lang="ts" setup>
import { FindRivalInfoList } from '@wailsjs/go/controller/RivalInfoController';
import { dto } from '@wailsjs/go/models';
import * as dayjs from 'dayjs';
import { DataTableColumns, useNotification } from 'naive-ui';
import { ref } from 'vue';

const notification = useNotification();

const pagination = false as const;
const columns = createColumns();
const data = ref<Array<dto.RivalInfoDto>>([]);
const loading = ref<boolean>(false);

function createColumns(): DataTableColumns<dto.RivalInfoDto> {
	return [
		{ title: "Name", key: "Name" },
		{ title: "Play Count", key: "PlayCount" },
		{ title: "Scorelog.db File Path", key: "ScoreLogPath", ellipsis: true, maxWidth: "150px" },
		{ title: "SongData.db File Path", key: "SongDataPath", ellipsis: true, maxWidth: "150px" },
		{ 
			title: "Last Sync Time",
			key: "UpdateAt",
			render: row => {
				return dayjs(row.UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
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

loadData();
</script>