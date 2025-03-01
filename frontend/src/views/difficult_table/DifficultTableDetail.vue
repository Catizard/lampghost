<template>
	<n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" min-height="500px"
		:loading="loading" :row-key="(row: dto.DiffTableDataDto) => row.ID" />
	<select-folder v-model:show="showFolderSelection" @submit="handleSubmit" />
</template>

<script setup lang="ts">
import { DataTableColumns, NButton, useNotification } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { h, reactive, ref, Ref, watch } from 'vue';
import { BindDiffTableDataToFolder, QueryDiffTableDataWithRival } from '@wailsjs/go/controller/DiffTableController';
import SelectFolder from '@/views/folder/SelectFolder.vue';
import ClearTag from '@/components/ClearTag.vue';
import { useI18n } from 'vue-i18n';

defineExpose({
	loadData,
});

const i18n = useI18n();
const { t } = i18n;
const notification = useNotification();
const loading = ref<boolean>(false);

const props = defineProps<{
	headerId?: number
	level?: string
	ghostRivalId?: number
	ghostRivalTagId?: number
}>()

const columns: DataTableColumns<dto.DiffTableDataDto> = [
	{ title: t('column.songName'), key: "Title", ellipsis: true, resizable: true },
	{ title: t('column.artist'), key: "Artist", resizable: true },
	{ title: t('column.count'), key: "PlayCount", minWidth: "100px", resizable: true },
	{
		title: t('column.clear'), key: "Lamp", minWidth: "100px", resizable: true,
		render(row) {
			return h(
				ClearTag,
				{
					clear: row.Lamp
				},
			)
		}
	},
	{
		title: t('column.ghost'), key: "GhostLamp", minWidth: "100px", resizable: true,
		render(row: dto.DiffTableDataDto) {
			return h(
				ClearTag,
				{
					clear: row.GhostLamp,
				},
			)
		}
	},
	{
		title: t('column.actions'),
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
				{ default: () => t('button.addToFolder') }
			)
		}
	}
];

let data: Ref<Array<dto.DiffTableDataDto>> = ref([]);
function loadData() {
	loading.value = true;
	// TODO: remove magic 1
	QueryDiffTableDataWithRival({
		ID: props.headerId,
		Level: props.level,
		RivalID: 1,
		GhostRivalID: props.ghostRivalId ?? 0,
		GhostRivalTagID: props.ghostRivalTagId ?? 0,
		Pagination: pagination,
	} as any)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			data.value = [...result.Rows];
			pagination.pageCount = result.Pagination.pageCount;
		}).catch(err => {
			notification.error({
				content: err,
				duration: 3000,
				keepAliveOnHover: true,
			});
		}).finally(() => {
			loading.value = false;
		});
}

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

watch(props, () => {
	loadData()
});

const candidateSongDataID = ref<number>(null);
const showFolderSelection = ref<boolean>(false);
function handleAddToFolder(ID: number) {
	candidateSongDataID.value = ID;
	showFolderSelection.value = true;
}

function handleSubmit(selected: [any]) {
	BindDiffTableDataToFolder(candidateSongDataID.value, selected)
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

loadData();
</script>

<i18n lang="json">{
	"en": {
		"column": {
			"songName": "Song Name",
			"artist": "Artist",
			"count": "Play Count",
			"clear": "Clear",
			"ghost": "Ghost",
			"actions": "Actions"
		},
		"button": {
			"addToFolder": "Add to Folder"
		}
	},
	"zh-CN": {
		"column": {
			"songName": "谱名",
			"artist": "作者",
			"count": "游玩次数",
			"clear": "点灯",
			"ghost": "影子灯",
			"actions": "操作"
		},
		"button": {
			"addToFolder": "添加至收藏夹"
		}
	}
}</i18n>