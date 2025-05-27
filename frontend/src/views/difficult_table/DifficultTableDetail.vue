<template>
	<n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" min-height="500px"
		:loading="loading" :row-key="(row: dto.DiffTableDataDto) => row.ID" @update-sorter="handleUpdateSorter" />
	<select-folder v-model:show="showFolderSelection" :sha256="candidateSongSha256" @submit="handleSubmit" />
</template>

<script setup lang="ts">
import { DataTableColumns, DataTableSortState, DropdownOption, NButton, NDropdown, NIcon, NText, NTooltip, useNotification } from 'naive-ui';
import { dto } from '@wailsjs/go/models';
import { h, reactive, ref, Ref, VNode, watch } from 'vue';
import { BindDiffTableDataToFolder, QueryDiffTableDataWithRival, SubmitSingleMD5DownloadTask } from '@wailsjs/go/main/App';
import SelectFolder from '@/views/folder/SelectFolder.vue';
import ClearTag from '@/components/ClearTag.vue';
import { useI18n } from 'vue-i18n';
import { WarningOutline } from '@vicons/ionicons5';
import { BrowserOpenURL } from '@wailsjs/runtime/runtime';

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

const sorter: Ref<Sorter> = ref({
	SortBy: null,
	SortOrder: null,
});
const columns: DataTableColumns<dto.DiffTableDataDto> = [
	{
		title: t('column.songName'), key: "Title", ellipsis: { tooltip: true }, resizable: true, sorter: true,
		render: (row: dto.DiffTableDataDto) => {
			let vnodes = [];
			// if (row.DataLost) {
			// 	vnodes.push(h(NIcon, { component: WarningOutline, color: "red" }));
			// }
			vnodes.push(h(NText, {}, { default: () => row.Title }));
			return vnodes;
		}
	},
	{ title: t('column.artist'), key: "Artist", ellipsis: { tooltip: true }, sorter: true, width: "125px" },
	{ title: t('column.count'), key: "PlayCount", width: "100px" },
	{
		title: t('column.clear'), key: "Lamp", width: "125px", resizable: true, sorter: true,
		render(row: dto.DiffTableDataDto) {
			return h(ClearTag, { clear: row.Lamp },)
		}
	},
	{
		// TODO: Change sorter from false to true when user choosed ghost rival
		title: t('column.ghost'), key: "GhostLamp", width: "125px", resizable: true, sorter: true,
		render(row: dto.DiffTableDataDto) {
			return h(ClearTag, { clear: row.GhostLamp, },)
		}
	},
	{
		title: t('column.actions'),
		key: "actions",
		resizable: true,
		width: "90px",
		render(row: dto.DiffTableDataDto) {
			return h(
				NDropdown,
				{
					trigger: "hover",
					options: [
						{ label: t('button.addToFolder'), key: "AddToFolder" },
						{ label: t('button.download'), key: "Download" },
						{ label: t('button.gotoURL'), key: "GotoURL", disabled: row.Url == "" },
						{ label: t('button.gotoURLDiff'), key: "GotoURLDiff", disabled: row.UrlDiff == "" },
						{ label: t('button.gotoLR2IR'), key: "GotoLR2IR", disabled: row.Md5 == "" },
						{ label: t('button.gotoMochaIR'), key: "GotoMochaIR", disabled: row.Sha256 == "" },
					],
					onSelect: (key: string) => {
						const md5 = row.Md5;
						const sha256 = row.Sha256;
						switch (key) {
							case 'AddToFolder': handleAddToFolder(row); break;
							case 'Download': handleSubmitSingleMD5DownloadTask(row); break;
							case "GotoURL": BrowserOpenURL(row.Url); break;
							case "GotoURLDiff": BrowserOpenURL(row.UrlDiff); break;
							case "GotoLR2IR": BrowserOpenURL(`https://www.dream-pro.info/~lavalse/LR2IR/search.cgi?mode=ranking&bmsmd5=${md5}`); break;
							case "GotoMochaIR": BrowserOpenURL(`https://mocha-repository.info/song.php?sha256=${sha256}`); break;
						}
					}
				},
				{ default: () => h(NButton, null, { default: () => '...' }) }
			);
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
		SortBy: sorter.value.SortBy,
		SortOrder: sorter.value.SortOrder,
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

function handleUpdateSorter(option: DataTableSortState | null) {
	// TODO: This is a pain in the a**
	switch (option.columnKey) {
		case "Title": sorter.value.SortBy = "title"; break;
		case "Artist": sorter.value.SortBy = "artist"; break;
		default: sorter.value.SortBy = option.columnKey as string; break;
	}
	if (option.order != false) {
		sorter.value.SortOrder = option.order;
	}
	loadData();
}

function handleSubmitSingleMD5DownloadTask(row: dto.DiffTableDataDto) {
	loading.value = true;
	SubmitSingleMD5DownloadTask(row.Md5, row.Title)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			notification.success({
				content: t('message.submitSuccess'),
				duration: 3000,
				keepAliveOnHover: true
			});
		}).catch(err => {
			notification.error({
				content: err,
				duration: 3000,
				keepAliveOnHover: true
			});
		}).finally(() => loading.value = false);
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
			"addToFolder": "Add to folder",
			"download": "Download",
			"gotoURL": "Open song url in browser",
			"gotoURLDiff": "Open sabun url in browser",
			"gotoLR2IR": "Open LR2IR in broswer",
			"gotoMochaIR": "Open MochaIR in broswer",
			"gotoBMSScoreViewer": "Open BMS Score Viewer in browser"
		},
		"message": {
			"bindSuccess": "Bind successfully",
			"bindFailedPrefix": "Failed to bind song to folder, error message: ",
			"submitSuccess": "Submit successfully"
		}
	},
	"zh-CN": {
		"column": {
			"songName": "谱面名称",
			"artist": "作者",
			"count": "游玩次数",
			"clear": "通关状态",
			"ghost": "对比通关状态",
			"actions": "操作"
		},
		"button": {
			"addToFolder": "添加至收藏夹",
			"download": "下载",
			"gotoURL": "在浏览器中打开单曲URL",
			"gotoURLDiff": "在浏览器中打开差分URL",
			"gotoLR2IR": "Open LR2IR in broswer",
			"gotoMochaIR": "Open MochaIR in broswer",
			"gotoBMSScoreViewer": "Open BMS Score Viewer in browser"
		},
		"message": {
			"bindSucess": "绑定成功",
			"bindFailedPrefix": "绑定失败, 错误信息: ",
			"submitSuccess": "提交成功"
		}
	}
}</i18n>
