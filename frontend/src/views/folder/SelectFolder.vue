<!-- Pop up page for selecting folders -->
<template>
	<n-modal :loading="loading" v-model:show="show" :title="t('title.bindToFavoriteFolder')" preset="dialog"
		:positive-text="t('button.submit')" :negative-text="t('button.cancel')"
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
