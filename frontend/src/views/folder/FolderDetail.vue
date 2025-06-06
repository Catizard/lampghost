<template>
	<n-data-table remote :loading="loading" :columns="columns" :data="data" :pagination="pagination" :bordered="false"
		:row-key="(row: dto.FolderDto) => row.ID" />
</template>

<script lang="ts" setup>
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NButton } from 'naive-ui';
import { h, reactive, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import ClearTag from '@/components/ClearTag.vue';
import { DelFolderContent, QueryFolderContentWithRival } from '@wailsjs/go/main/App';
import TableTags from '@/components/TableTags.vue';

const loading = ref(false);
const { t } = useI18n();

const props = defineProps<{
	folderId: number,
	rivalId?: number,
	type: "table" | "folder"
}>();

let data: Ref<dto.FolderContentDto[]> = ref([]);
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
const columns: DataTableColumns<dto.FolderContentDto> = [
	{ title: t('column.name'), key: "Title" },
	{
		title: (): string => {
			if (props.type == "folder") {
				return t('column.tag.folder');
			} else if (props.type == "table") {
				return t('column.tag.table');
			}
		},
		key: "Tag",
		width: "200px",
		render(row: dto.FolderContentDto) {
			return h(TableTags, { tableTags: row.TableTags });
		}
	},
	{
		title: t('column.clear'), key: "Clear",
		width: "125px",
		render: (row: dto.FolderContentDto) => {
			return h(ClearTag, { clear: row.Lamp },)
		}
	},
	{
		title: t('column.actions'),
		key: "actions",
		width: "150px",
		render(row: dto.FolderContentDto) {
			return h(
				NButton,
				{
					strong: true,
					tertiary: true,
					size: "small",
					type: "error",
					onClick: () => deleteFolderContent(row.ID),
				},
				{ default: () => t('button.delete') },
			);
		},
	},
];

function loadData() {
	loading.value = true;
	// TODO: remove magic 1
	QueryFolderContentWithRival({
		RivalID: 1,
		FolderID: props.folderId,
		Pagination: pagination,
	} as any).then(result => {
		if (result.Code != 200) {
			return Promise.reject(result.Msg);
		}
		data.value = [...result.Rows];
		pagination.pageCount = result.Pagination.pageCount;
	})
		.catch(err => window.$notifyError(err))
		.finally(() => loading.value = false)
}

function deleteFolderContent(id: number) {
	DelFolderContent(id)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			window.$notifySuccess(t('message.deleteSuccess'));
			loadData();
		})
		.catch((err) => {
			window.$notifyError(err);
			loadData();
		});
}

watch(props, () => loadData());
loadData();
</script>

<i18n lang="json">{
	"en": {
		"column": {
			"name": "Name",
			"tag": {
				"folder": "Tag",
				"table": "External Tag"
			},
			"clear": "Clear",
			"actions": "Actions"
		},
		"button": {
			"delete": "Delete"
		}
	},
	"zh-CN": {
		"column": {
			"name": "谱面名称",
			"tag": {
				"folder": "难度表标签",
				"table": "外部难度表标签"
			},
			"clear": "通关状态",
			"actions": "操作"
		},
		"button": {
			"delete": "删除"
		}
	}
}</i18n>
