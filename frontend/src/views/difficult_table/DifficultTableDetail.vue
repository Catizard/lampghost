<template>
	<n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" min-height="500px"
		:loading="loading" :row-key="(row: dto.DiffTableDataDto) => row.ID" />
	<select-folder v-model:show="showFolderSelection" :sha256="candidateSongSha256" @submit="handleSubmit" />
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
	{ title: t('column.songName'), key: "Title", width: "300px", ellipsis: { tooltip: true }, resizable: true },
	{ title: t('column.artist'), key: "Artist", ellipsis: { tooltip: true }, },
	{ title: t('column.count'), key: "PlayCount", width: "100px", },
	{
		title: t('column.clear'), key: "Lamp", width: "100px", resizable: true,
		render(row: dto.DiffTableDataDto) {
			return h(ClearTag, { clear: row.Lamp },)
		}
	},
	{
		title: t('column.ghost'), key: "GhostLamp", width: "100px", resizable: true,
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
		render(row: dto.DiffTableDataDto) {
			return h(
				NButton,
				{
					strong: true,
					tertiary: true,
					size: "small",
					onClick: () => handleAddToFolder(row),
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

const candidateDiffDataId = ref<number | null>(null);
const candidateSongSha256 = ref<string | null>(null);
const showFolderSelection = ref<boolean>(false);
function handleAddToFolder(row: dto.DiffTableDataDto) {
	candidateDiffDataId.value = row.ID;
	candidateSongSha256.value = row.Sha256;
	showFolderSelection.value = true;
}

function handleSubmit(selected: [any]) {
	BindDiffTableDataToFolder(candidateDiffDataId.value, selected)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			notification.success({
				content: t('message.bindSuccess'),
				duration: 3000,
				keepAliveOnHover: true
			});
		}).catch((err) => {
			notification.error({
				content: t('message.bindFailedPrefix') + err,
				duration: 3000,
				keepAliveOnHover: true
			});
		});
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
		},
		"message": {
			"bindSuccess": "Bind successfully",
			"bindFailedPrefix": "Failed to bind song to folder, error message: "
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
		},
		"message": {
			"bindSucess": "绑定成功",
			"bindFailedPrefix": "绑定失败, 错误信息: "
		}
	}
}</i18n>