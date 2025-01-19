<template>
	<n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" max-height="400px"
		:loading="loading" />
	<select-folder v-model:show="showFolderSelection" @submit="handleSubmit" />
</template>

<script setup lang="ts">
import { DataTableColumns, NButton, useNotification } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { h, ref, Ref, watch } from 'vue';
import { QueryDiffTableDataWithRival } from '@wailsjs/go/controller/DiffTableController';
import SelectFolder from '@/views/folder/SelectFolder.vue';
import { BindSongToFolder } from '@wailsjs/go/controller/FolderController';

const notification = useNotification();
const loading = ref<boolean>(false);

const props = defineProps<{
	headerId?: number
	level?: string
}>()

function createColumns(): DataTableColumns<dto.DiffTableDataDto> {
	return [
		{ title: "Song Name", key: "Title", resizable: true },
		{ title: "Artist", key: "Artist", resizable: true },
		{ title: "Play Count", key: "PlayCount", minWidth: "100px", resizable: true },
		{ title: "Clear", key: "Lamp", minWidth: "100px", resizable: true },
		{
			title: "Action",
			key: "actions",
			resizable: true,
			minWidth: "150px",
			render(row) {
				return h(
					NButton,
					{
						strong: true,
						tertiary: true,
						size: "small",
						onClick: () => handleAddToFolder(row.ID),
					},
					{ default: () => "添加至收藏夹" }
				)
			}
		}
	];
}
const columns = createColumns();

let data: Ref<Array<any>> = ref([]);
function loadData(headerId: number, level: string) {
	loading.value = true;
	// TODO: remove magic 1
	QueryDiffTableDataWithRival(headerId, level, 1)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg)
			}
			data.value = [...result.Rows].map(row => {
				return {
					key: row.ID,
					...row
				}
			})
		}).catch(err => {
			notifyError("读取数据失败:" + err);
		}).finally(() => {
			loading.value = false;
		})
}

const pagination = false as const;

watch(props, (newProps) => {
	loadData(newProps.headerId, newProps.level)
}, { immediate: true })

const candidateSongDataID = ref<number>(null);
const showFolderSelection = ref<boolean>(false);
function handleAddToFolder(ID: number) {
	candidateSongDataID.value = ID;
	showFolderSelection.value = true;
}

function handleSubmit(selected: [any]) {
	BindSongToFolder(candidateSongDataID.value, selected)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			notifySuccess("添加成功");
		}).catch((err) => {
			notifyError("添加至收藏夹失败:" + err);
		});
}

function notifySuccess(msg: string) {
	notification.success({
		content: msg,
		duration: 5000,
		keepAliveOnHover: true
	})
}

function notifyError(msg: string) {
	notification.error({
		content: msg,
		duration: 5000,
		keepAliveOnHover: true
	})
}

</script>