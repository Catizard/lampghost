<template>
	<perfect-scrollbar>
		<n-flex justify="space-between">
			<n-h1 prefix="bar" style="text-align: start;">
				<n-text type="primary">{{ t('title') }}</n-text>
			</n-h1>
		</n-flex>
		<n-flex justify="start">
			<n-input :placeholder="t('searchNamePlaceholder')" v-model:value="searchNameLike" @keyup.enter="loadData()"
				style="width: 350px;" />
			<!-- <n-button>{{ t('button.chooseClearType') }}</n-button>
      <n-button>{{ t('button.minimumClearType') }}</n-button> -->
		</n-flex>
		<n-data-table remote :columns="columns" :data="data" :pagination="pagination" :bordered="false" :loading="loading"
			:row-key="row => row.ID" />
	</perfect-scrollbar>
	<select-folder v-model:show="showFolderSelection" :songDataId="candidateSongDataID" @submit="handleSubmit" />
</template>

<script lang="ts" setup>
import ClearTag from '@/components/ClearTag.vue';
import { QueryRivalScoreLogPageList } from '@wailsjs/go/main/App';
import { DataTableColumns, DataTableRowKey, NButton, NTag, useNotification } from 'naive-ui';
import { h, onMounted, reactive, Ref, ref, VNode } from 'vue';
import SelectFolder from '../folder/SelectFolder.vue';
import { BindRivalSongDataToFolder } from '@wailsjs/go/main/App';
import { dto, vo } from '@wailsjs/go/models';
import { useI18n } from 'vue-i18n';
import dayjs from 'dayjs';

const i18n = useI18n();
const { t } = i18n;
const notification = useNotification();
const loading = ref<boolean>(false);
const showFolderSelection = ref<boolean>(false);
const candidateSongDataID = ref<number>(null);
const searchNameLike: Ref<string | null> = ref(null);

function handleAddToFolder(ID: number) {
	candidateSongDataID.value = ID;
	showFolderSelection.value = true;
}

function handleSubmit(folderIds: number[]) {
	BindRivalSongDataToFolder(candidateSongDataID.value, folderIds)
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

function createColumns(): DataTableColumns<dto.RivalScoreLogDto> {
	return [
		{ title: t('column.name'), key: "Title", resizable: true },
		{
			title: t('column.tag'), key: "Tag", minWidth: "100px", resizable: true,
			render(row: dto.RivalScoreLogDto) {
				const nodes: Array<VNode> = [];
				if (row.TableTags.length == 0) {
					return nodes;
				}
				row.TableTags.forEach(tag => {
					const props = {
						size: "small",
						style: {
							"margin-right": "5px"
						},
						color: {

						},
					};
					if (tag.TableTagColor.length > 0) {
						(props as any).color.color = tag.TableTagColor;
					}
					if (tag.TableTagTextColor.length > 0) {
						(props as any).color.textColor = tag.TableTagTextColor;
					}
					const node = h(
						NTag,
						props as any,
						{ default: () => tag.TableSymbol + tag.TableLevel }
					)
					nodes.push(node);
				});
				return nodes;
			}
		},
		{
			title: t('column.clear'), key: "Clear", minWidth: "100px", resizable: true,
			render(row: dto.RivalScoreLogDto) {
				return h(ClearTag, { clear: row.Clear }, {});
			}
		},
		{
			title: t('column.time'), key: "RecordTime", minWidth: "100px", resizable: true,
			render(row: dto.RivalScoreLogDto) {
				return dayjs(row.RecordTime).format('YYYY-MM-DD HH:mm:ss');
			}
		},
		{
			title: t('column.minbp'), key: "MinBP", minWidth: "100px", resizable: true,
			render(row: dto.RivalScoreLogDto) {
				return row.Minbp;
			}
		},
		{
			title: t('column.actions'), key: "actions", resizable: true, minWidth: "150px",
			render(row: dto.RivalScoreLogDto) {
				return h(
					NButton,
					{
						strong: true,
						tertiary: true,
						size: "small",
						onClick: () => handleAddToFolder(row.RivalSongDataID),
					},
					{ default: () => t('button.addToFolder') }
				)
			}
		}
	]
}

function loadData() {
	loading.value = true;
	// TODO: remove magical 1 here
	let arg: vo.RivalScoreLogVo = {
		RivalID: 1,
		Pagination: pagination,
		SongNameLike: searchNameLike.value,
		NoCourseLog: true,
	} as any;
	QueryRivalScoreLogPageList(arg)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			console.log(result);
			data.value = [...result.Rows];
			pagination.pageCount = result.Pagination.pageCount;
		})
		.catch(err => {
			notification.error({
				content: t('message.loadRecentRecordFailedPrefix') + err,
				duration: 3000,
				keepAliveOnHover: true
			});
		}).finally(() => {
			loading.value = false;
		});
}

const columns = createColumns();
let data: Ref<Array<any>> = ref([]);
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

onMounted(() => {
	loadData();
})
</script>

<i18n lang="json">{
	"en": {
		"title": "Recent Play",
		"button": {
			"chooseClearType": "Choose Clear Type",
			"minimumClearType": "Minimum Clear Type",
			"addToFolder": "Add to Folder"
		},
		"column": {
			"name": "Song Name",
			"tag": "Tag",
			"clear": "Clear",
			"time": "Record Time",
			"minbp": "Min BP",
			"actions": "Actions"
		},
		"message": {
			"loadRecentRecordFailedPrefix": "Load recent records failed, error message: ",
			"bindSuccess": "Bind successfully",
			"bindFailedPrefix": "Failed to bind song to folder, error message: "
		},
		"searchNamePlaceholder": "Search Song/Sabun Name"
	},
	"zh-CN": {
		"title": "最近游玩",
		"button": {
			"chooseClearType": "筛选点灯记录",
			"minimumClearType": "筛选最小灯记录",
			"addToFolder": "添加至收藏夹"
		},
		"column": {
			"name": "谱面名称",
			"tag": "难度表标签",
			"clear": "通关状态",
			"time": "记录时间",
			"minbp": "最小BP",
			"actions": "操作"
		},
		"message": {
			"loadRecentRecordFailedPrefix": "加载最近游玩记录失败，错误信息: ",
			"bindSucess": "绑定成功",
			"bindFailedPrefix": "绑定失败, 错误信息: "
		},
		"searchNamePlaceholder": "搜索歌曲/差分名称"
	}
}</i18n>
