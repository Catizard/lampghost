<!-- Pop up page for selecting folders -->
<template>
	<n-modal v-model:show="show" preset="dialog" positive-text="提交" negative-text="取消"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" closable
		@close="() => { show = false }">
		<n-data-table :columns="columns" :data="data" :pagination="false" :bordered="false" :row-key="rowKey"
			@update:checked-row-keys="handleCheck" :loading="loading" />
	</n-modal>
</template>

<script setup lang="ts">
import { ref, watch, watchEffect } from 'vue';
import { dto } from '@wailsjs/go/models';
import { FindFolderList } from '@wailsjs/go/controller/FolderController';
import { DataTableColumns, DataTableRowKey } from 'naive-ui';

const show = defineModel<boolean>("show");
const emit = defineEmits<{
	(e: 'submit', selected: Array<any>): void;
}>();

const loading = ref(false);
const data = ref<dto.FolderDto[]>([]);
const checkedRowKeysRef = ref<DataTableRowKey[]>([]);
const columns: DataTableColumns<dto.FolderDto> = [
	{ type: "selection" },
	{ title: "Folder Name", key: "FolderName" }
] as const;

function reload() {
	loading.value = true;
	FindFolderList()
		.then(result => {
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

function rowKey(folder: dto.FolderDto): number {
	return folder.ID
}

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