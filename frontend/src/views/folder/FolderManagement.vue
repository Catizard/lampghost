<template>
	<n-flex justify="space-between">
		<n-h1 prefix="bar" style="text-align: left">
			<n-text type="primary">{{ t('title') }}</n-text>
		</n-h1>
		<n-flex justify="end">
			<n-button type="primary" @click="showAddModal = true">{{ t('button.addFolder') }}</n-button>
		</n-flex>
	</n-flex>
	<FolderTable ref="folderTableRef"/>
	<FolderAddForm v-model:show="showAddModal" @refresh="reload" />
</template>

<script setup lang="ts">
import { ref } from "vue";
import { NButton, } from "naive-ui";
import { useI18n } from "vue-i18n";
import FolderAddForm from "./FolderAddForm.vue";
import FolderTable from "./FolderTable.vue";

const i18n = useI18n();
const { t } = i18n;

const showAddModal = ref(false);
const folderTableRef = ref<InstanceType<typeof FolderTable>>(null);

function reload() {
	folderTableRef.value.loadData();
}
</script>

<i18n lang="json">{
	"en": {
		"title": "Folder",
		"button": {
			"addFolder": "Add Folder"
		}
	},
	"zh-CN": {
		"title": "收藏夹",
		"button": {
			"addFolder": "新增收藏夹"
		}
	}
}</i18n>
