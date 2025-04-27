<!-- Pop up page for selecting folders -->
<template>
	<n-modal :loading="loading" v-model:show="show" :title="t('dialog.title')" preset="dialog"
		:positive-text="t('dialog.positiveText')" :negative-text="t('dialog.negativeText')"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" closable
		@close="() => { show = false }">
		<n-button type="primary" @click="handleClickAddFolder">{{ t('button.addFolder') }}</n-button>
		<n-data-table :columns="columns" :data="data" :pagination="false" :bordered="false"
			:row-key="(row: dto.FolderDto) => row.ID" @update:checked-row-keys="handleCheck" :loading="loading"
			:checked-row-keys="checkedRowKeysRef" />
	</n-modal>

	<FolderAddForm v-model:show="showAddModal" @refresh="reload" />
</template>

<script setup lang="ts">
import { ref, watch, watchEffect } from 'vue';
import { dto } from '@wailsjs/go/models';
import { FindFolderList } from '@wailsjs/go/main/App';
import { DataTableColumns, DataTableRowKey, useNotification } from 'naive-ui';
import { useI18n } from 'vue-i18n';
import FolderAddForm from './FolderAddForm.vue';

const { t } = useI18n();
const show = defineModel<boolean>("show");
const props = defineProps<{
	songDataId?: number,
	sha256?: string
}>();
const emit = defineEmits<{
	(e: 'submit', selected: Array<any>): void
}>();

const loading = ref(false);
const showAddModal = ref(false);
const data = ref<dto.FolderDto[]>([]);
const checkedRowKeysRef = ref<DataTableRowKey[]>([]);
const columns: DataTableColumns<dto.FolderDto> = [
	{ type: "selection" },
	{ title: t('column.name'), key: "FolderName" }
] as const;

function reload() {
	loading.value = true;
	FindFolderList({
		IgnoreRivalSongDataID: props.songDataId,
		IgnoreSha256: props.sha256
	} as any).then(result => {
		if (result.Code != 200) {
			return Promise.reject(result.Msg);
		}
		data.value = [...result.Rows];
	}).catch(err => {
		console.error(err);
	}).finally(() => {
		loading.value = false;
	})
}

function handleClickAddFolder() {
	console.log('here');
	showAddModal.value = true;
}

watch(show, (newValue, oldValue) => {
	if (newValue == true) {
		reload();
	}
}, { deep: true });

function handleCheck(rowKeys: DataTableRowKey[]) {
	checkedRowKeysRef.value = rowKeys
}

function handlePositiveClick() {
	emit('submit', checkedRowKeysRef.value);
}

function handleNegativeClick() {
	data.value = [];
	show.value = false;
}
</script>

<i18n lang="json">{
	"en": {
		"dialog": {
			"title": "Bind to Folder",
			"positiveText": "Submit",
			"negativeText": "Cancel"
		},
		"column": {
			"name": "Folder Name"
		},
		"button": {
			"addFolder": "Add Folder"
		}
	},
	"zh-CN": {
		"dialog": {
			"title": "加入收藏夹",
			"positiveText": "提交",
			"negativeText": "取消"
		},
		"column": {
			"name": "收藏夹名称"
		},
		"button": {
			"addFolder": "添加收藏夹"
		}
	}
}</i18n>