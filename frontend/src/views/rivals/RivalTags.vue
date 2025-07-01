<template>
	<perfect-scrollbar>
		<n-flex justify="start">
			<n-h1 prefix="bar" style="text-align: start;">
				<n-text type="primary">{{ t('title.playerTags') }}</n-text>
			</n-h1>
		</n-flex>
		<n-flex justify="start" style="margin-bottom: 15px;">
			<n-select :loading="tableLoading" v-model:value="currentRivalID" :options="rivalOptions" style="width: 200px;" />
			<n-button style="margin-left: auto;" type="primary" @click="showAddModal = true">{{ t('button.addCustomTag') }}</n-button>
		</n-flex>
		<n-data-table remote :columns="columns" :data="data" :pagination="pagination" :loading="tableLoading"
			:row-key="(row: dto.RivalTagDto) => row.ID" />
	</perfect-scrollbar>
	<n-modal v-model:show="showAddModal" preset="dialog" :title="t('title.addPlayerTag')"
		:positive-text="t('button.submit')" :negative-text="t('button.submit')"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
		<n-form ref="formRef" :model="formData" :rules="rules">
			<n-form-item :label="t('form.labelTagName')" path="TagName">
				<n-input v-model:value="formData.TagName" :placeholder="t('form.placeholderTagName')" />
			</n-form-item>
			<n-form-item :label="t('form.labelTagRecordTime')" path="RecordTimestamp">
				<n-date-picker clearable v-model:value="formData.RecordTimestamp" type="datetime" />
			</n-form-item>
		</n-form>
	</n-modal>
</template>

<script setup lang="ts">
import { FindRivalInfoList } from '@wailsjs/go/main/App';
import { AddRivalTag, DeleteRivalTagByID, QueryRivalTagPageList, RevertRivalTagEnabledState } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, FormInst, FormRules, NButton, SelectOption, useDialog } from 'naive-ui';
import { h, reactive, Ref, ref, watch } from 'vue';
import dayjs from 'dayjs';
import { useI18n } from 'vue-i18n';
import YesNotTag from '@/components/YesNotTag.vue';

const i18n = useI18n();
const { t } = i18n;
const dialog = useDialog();

const tableLoading = ref(false);
const currentRivalID: Ref<number | null> = ref(null);
const rivalOptions: Ref<Array<SelectOption>> = ref([]);
function loadRivalOptions() {
	tableLoading.value = true;
	FindRivalInfoList()
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			if (result.Rows.length == 0) {
				return Promise.reject(t('message.noRivalError'));
			}
			rivalOptions.value = result.Rows.map((row: dto.RivalInfoDto) => {
				return {
					label: row.Name,
					value: row.ID,
				} as SelectOption
			});
			currentRivalID.value = rivalOptions.value[0].value as number;
		})
		.catch(err => window.$notifyError(err))
		.finally(() => tableLoading.value = false);
}
loadRivalOptions();

const columns: DataTableColumns<dto.RivalTagDto> = [
	{ title: t('column.tagName'), key: "TagName", width: "200px", ellipsis: { tooltip: true } },
	{
		title: t('column.tagTime'),
		width: "125px",
		key: "RecordTime",
		render: (row: dto.RivalTagDto) => dayjs(row.RecordTime).format('YYYY-MM-DD HH:mm:ss')
	},
	{
		title: t('column.enabled'), key: "Enabled",
		width: "75px",
		render: (row: dto.RivalTagDto) => {
			return h(
				YesNotTag,
				{ state: row.Enabled, onClick: () => {} }
			);
		}
	},
	{
		title: t('column.actions'), key: "Actions",
		width: "100px",
		render: (row: dto.RivalTagDto) => {
			const deleteTagButton = row.Generated == false ? h(
				NButton,
				{
					strong: true,
					tertiary: true,
					size: 'small',
					type: "error",
					style: {
						"margin-left": "5px",
					},
					onClick: () => deleteTag(row),
				},
				{ default: () => t('button.delete') }
			) : null;
			const revertEnabledStateButton = h(
				NButton,
				{
					strong: true,
					tertiary: true,
					size: "small",
					onClick: () => revertTagEnabledState(row),
				},
				{ default: () => row.Enabled ? t('button.disable') : t('button.enable') }
			);
			if (deleteTagButton != null) {
				return [revertEnabledStateButton, deleteTagButton];
			}
			return revertEnabledStateButton;
		}
	}
];
const data: Ref<Array<dto.RivalTagDto>> = ref([]);
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
	},
	reset: () => {
		pagination.page = 1;
		loadData();
	}
});
function loadData() {
	tableLoading.value = true;
	QueryRivalTagPageList({
		RivalId: currentRivalID.value,
		Pagination: pagination,
		NoIgnoreEnabled: true,
	} as any)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			data.value = [...result.Rows];
			pagination.pageCount = result.Pagination.pageCount;
		})
		.catch(err => window.$notifyError(err))
		.finally(() => tableLoading.value = false);
}

const showAddModal = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = ref({
	RivalID: null,
	TagName: null,
	RecordTimestamp: null
});
const rules: FormRules = {
	RecordTimestamp: {
		type: 'number',
		required: true,
		message: t('rule.missingTagRecordTime'),
		trigger: ["input", "blur"]
	}
}
function handlePositiveClick(): boolean {
	formRef.value
		?.validate()
		.then(() => {
			AddRivalTag({
				RivalId: currentRivalID.value,
				...formData.value
			} as any)
				.then(result => {
					if (result.Code != 200) {
						return Promise.reject(result.Msg);
					}
					showAddModal.value = false;
				})
				.catch(err => window.$notifyError(err))
				.finally(() => loadData());
		})
		.catch((err) => { });
	return false;
}
function handleNegativeClick() {
	formData.value.TagName = null;
	formData.value.RecordTimestamp = null;
}

function deleteTag(tag: dto.RivalTagDto) {
	dialog.warning({
		title: t('deleteDialog.title'),
		positiveText: t('deleteDialog.positiveText'),
		negativeText: t('deleteDialog.negativeText'),
		onPositiveClick: () => {
			DeleteRivalTagByID(tag.ID)
				.then(result => {
					if (result.Code != 200) {
						return Promise.reject(result.Msg);
					}
					loadData();
				}).catch(err => window.$notifyError(err));
		}
	})
}

function revertTagEnabledState(tag: dto.RivalTagDto) {
	RevertRivalTagEnabledState(tag.ID)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			loadData();
		}).catch(err => window.$notifyError(err));
	;
}

// Watch: Whenever changing current rival, reset current page to the first one
// NOTE: There is no need to call loadData(), reset method would call it
watch(currentRivalID, () => {
	pagination.reset();
})
</script>
