<!-- Pop up page for selecting folders -->
<template>
	<n-modal :loading="loading" v-model:show="show" :title="t('dialog.title')" preset="dialog"
		:positive-text="t('dialog.positiveText')" :negative-text="t('dialog.negativeText')"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" closable
		@close="() => { show = false }">
		<n-button type="primary" @click="handleClickAddFolder">{{ t('button.addFavoriteFolder') }}</n-button>
		<SelectUnboundFolder ref="selectUnboundFolderRef" type="folder" v-model:checkedFolderIds="checkedFolderIds"
			:sha256="sha256" :customTableId="customTableId" />
	</n-modal>

	<FolderAddForm type="folder" v-model:show="showAddModal" :customTableId="customTableId" @refresh="reload" />
</template>

<script setup lang="ts">
import { nextTick, Ref, ref, watch } from 'vue';
import { dto } from '@wailsjs/go/models';
import { useI18n } from 'vue-i18n';
import FolderAddForm from './FolderAddForm.vue';
import SelectUnboundFolder from './SelectUnboundFolder.vue';

const { t } = useI18n();
const show = defineModel<boolean>("show");
const props = defineProps<{
	sha256?: string
	customTableId?: number
}>();
const emit = defineEmits<{
	(e: 'submit', selected: Array<any>): void
}>();

const loading = ref(false);
const showAddModal = ref(false);
const data = ref<dto.FolderDto[]>([]);
let checkedFolderIds: Ref<number[]> = ref([]);
const selectUnboundFolderRef: Ref<InstanceType<typeof SelectUnboundFolder>> = ref(null);

function handleClickAddFolder() {
	showAddModal.value = true;
}

function reload() {
	selectUnboundFolderRef.value.reload();
}

watch(show, (newValue, oldValue) => {
	if (newValue == true) {
		// NOTE: SelectUnboundFolder component is under NModel
		// Therefore, we have to wait until it mounted
		const spinUntilMounted = () => {
			if (selectUnboundFolderRef == null) {
				nextTick(() => spinUntilMounted());
			} else {
				selectUnboundFolderRef.value.reload();
			}
		}
		nextTick(() => spinUntilMounted());
	}
}, { deep: true });

function handlePositiveClick() {
	emit('submit', checkedFolderIds.value);
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
		}
	},
	"zh-CN": {
		"dialog": {
			"title": "加入收藏夹",
			"positiveText": "提交",
			"negativeText": "取消"
		}
	}
}</i18n>
