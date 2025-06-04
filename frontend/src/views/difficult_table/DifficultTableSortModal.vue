<template>
	<n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('modal.title')" style="width: 75%;"
		:positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
		<n-data-table ref="tableRef" :loading="loading" :columns="columns" :data="data" :bordered="false"
			:row-key="(row: dto.DiffTableHeaderDto) => row.ID" />
	</n-modal>
</template>

<script lang="ts" setup>
import { nextTick, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { dto } from '@wailsjs/go/models';
import { FindDiffTableHeaderTree, UpdateHeaderOrder } from '@wailsjs/go/main/App';
import { NDataTable } from 'naive-ui';
import Sortable from "sortablejs";

const show = defineModel<boolean>("show");
const emit = defineEmits<{
	(e: 'refresh'): void
}>();

const { t } = useI18n();
const loading = ref(false);

const columns = [
	{ title: t('column.name'), key: "Name" },
	{
		title: t('column.rowIndex'), key: "RowIndex",
		render: (_, rowIndex: number) => {
			return rowIndex + 1
		}
	}
];

const tableRef = ref<InstanceType<typeof NDataTable>>();
let data = reactive([]);
function loadDiffTableData() {
	loading.value = true;
	FindDiffTableHeaderTree(null)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			data = [...result.Rows]
			// when data.length == 0, hookSortable would make whole application stall
			// because the body of the table can never be rendered successfully
			if (data.length > 0) {
				// Here goes the magic
				hookSortable();
			}
		})
		.catch(err => window.$notifyError(err))
		.finally(() => loading.value = false);
}

function hookSortable() {
	const el: HTMLDivElement = tableRef.value?.$el;
	const tbody: HTMLElement = el.querySelector('.n-data-table-tbody');
	// Naive-UI's data table could not completely render in one tick
	// Therefore we need to use this hack to hook sortable correctly
	if (tbody == null) {
		nextTick(() => {
			hookSortable();
		});
		return;
	}
	Sortable.create(tbody, {
		onEnd(evt) {
			const oldElem = data[evt.oldIndex!];
			data.splice(evt.oldIndex!, 1);
			data.splice(evt.newIndex!, 0, oldElem);
		}
	});
}

function handlePositiveClick(): boolean {
	loading.value = true;
	UpdateHeaderOrder(data.map((u: dto.DiffTableHeaderDto) => u.ID))
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			show.value = false;
			emit('refresh');
		}).catch(err => window.$notifyError(err)).finally(() => loading.value = false);
	return false
}

function handleNegativeClick() {
	show.value = false;
}

watch(show, (newValue: boolean) => {
	console.log(show);
	if (newValue == true) {
		loadDiffTableData();
	}
});
</script>

<i18n lang="json">{
	"en": {
		"modal": {
			"title": "Customize Table Order",
			"positiveText": "Submit",
			"negativeText": "Cancel"
		},
		"column": {
			"name": "Name",
			"rowIndex": "Index"
		}
	},
	"zh-CN": {
		"modal": {
			"title": "自定义难度表排序",
			"positiveText": "提交",
			"negativeText": "取消"
		},
		"column": {
			"name": "名称",
			"rowIndex": "序号"
		}
	}
}</i18n>