<!-- Similar to SelectFolder, but candidates are difficult folders from custom tables -->
<template>
	<n-modal :loading="loading" v-model:show="show" :title="t('dialog.title')" preset="dialog"
		:positive-text="t('dialog.positiveText')" :negative-text="t('dialog.negativeText')"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" closable
		@close="() => { show = false }">
		<n-flex justify="space-between">
			<n-select v-model:value="currentCustomTableID" :options="customTableOptions" style="width: 200px;" />
			<n-button type="primary" @click="handleClickAddFolder">{{ t('button.addFolder') }}</n-button>
		</n-flex>
		<SelectUnboundFolder ref="selectUnboundFolderRef" type="folder" v-model:checkedFolderIds="checkedFolderIds"
			:sha256="sha256" :customTableId="currentCustomTableID" />
	</n-modal>

	<FolderAddForm type="table" v-model:show="showAddModal" @refresh="reload" />
</template>

<script setup lang="ts">
import { nextTick, Ref, ref, watch } from 'vue';
import { dto } from '@wailsjs/go/models';
import { useI18n } from 'vue-i18n';
import FolderAddForm from '../folder/FolderAddForm.vue';
import SelectUnboundFolder from '../folder/SelectUnboundFolder.vue';
import { SelectOption } from 'naive-ui';
import { FindCustomDiffTableList } from '@wailsjs/go/main/App';
import { ReloadDiffTableHeader } from '../../../wailsjs/go/main/App';

const { t } = useI18n();
const show = defineModel<boolean>("show");
const props = defineProps<{
	sha256?: string
}>();
const emit = defineEmits<{
	(e: 'submit', selected: Array<any>): void
}>();

const loading = ref(false);
const showAddModal = ref(false);
const data = ref<dto.FolderDto[]>([]);
let checkedFolderIds: Ref<number[]> = ref([]);
const selectUnboundFolderRef: Ref<InstanceType<typeof SelectUnboundFolder>> = ref(null);

const currentCustomTableID: Ref<number | null> = ref(null);
const customTableOptions: Ref<SelectOption[]> = ref([]);
async function loadCustomTableOptions() {
	loading.value = true;
	try {
		const result = await FindCustomDiffTableList({
			IgnoreDefaultTable: true
		} as any)
		if (result.Code != 200) {
			throw result.Msg;
		}
		if (result.Rows.length == 0) {
			throw t('message.noTableError');
		}
		customTableOptions.value = result.Rows.map((row: dto.CustomDiffTableDto): SelectOption => {
			return {
				label: row.Name,
				value: row.ID
			}
		});
		currentCustomTableID.value = customTableOptions.value[0].value as number;
		console.log('setting custom table id to ', currentCustomTableID.value);
	} catch (e) {
		window.$notifyError(e);
	}
	loading.value = false;
}

function handleClickAddFolder() {
	showAddModal.value = true;
}

function reload() {
	selectUnboundFolderRef.value.reload();
}

watch(show, (newValue, oldValue) => {
	if (newValue == true) {
		loadCustomTableOptions();
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
			"title": "Bind to Custom Table",
			"positiveText": "Submit",
			"negativeText": "Cancel"
		},
		"button": {
			"addFolder": "Add Difficult"
		},
		"message": {
			"noTableError": "Cannot handle no custom difficult table data currenlty, please add at least one custom table first"
		}
	},
	"zh-CN": {
		"dialog": {
			"title": "加入自定义难度表",
			"positiveText": "提交",
			"negativeText": "取消"
		},
		"button": {
			"addFolder": "添加难度"
		},
		"message": {
			"noTableError": "目前无法处理一个自定义难度表都没有的情况，请至少先添加一个难度表"
		}
	}
}</i18n>
