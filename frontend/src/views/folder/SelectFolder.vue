<!-- Pop up page for selecting folders -->
<template>
	<n-modal v-model:show="show" preset="dialog" :positive-text="t('dialog.positiveText')"
		:negative-text="t('dialog.negativeText')" @positive-click="handlePositiveClick"
		@negative-click="handleNegativeClick" closable @close="() => { show = false }">
		<n-data-table :columns="columns" :data="data" :pagination="false" :bordered="false"
			:row-key="(row: dto.FolderDto) => row.ID" @update:checked-row-keys="handleCheck" :loading="loading"
			:checked-row-keys="checkedRowKeysRef" />
	</n-modal>
</template>

<script setup lang="ts">
import { ref, watch, watchEffect } from 'vue';
import { dto } from '@wailsjs/go/models';
import { FindFolderList } from '@wailsjs/go/controller/FolderController';
import { DataTableColumns, DataTableRowKey, useNotification } from 'naive-ui';
import { useI18n } from 'vue-i18n';

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
			"positiveText": "Submit",
			"negativeText": "Cancel"
		},
		"column": {
			"name": "Folder Name"
		}
	},
	"zh-CN": {
		"dialog": {
			"positiveText": "提交",
			"negativeText": "取消"
		},
		"column": {
			"name": "收藏夹名称"
		}
	}
}</i18n>