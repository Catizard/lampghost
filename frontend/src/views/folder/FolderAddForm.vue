<template>
	<n-modal :loading="loading" v-model:show="show" preset="dialog" :title="title"
		:positive-text="t('modal.positiveText')" :negative-text="t('modal.negativeText')"
		@positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
		<n-form ref="formRef" :model="formData" :rules="rules">
			<n-form-item :label="t('form.labelName')" path="FolderName">
				<n-input v-model:value="formData.FolderName" :placeholder="t('form.placeholderName')" />
			</n-form-item>
		</n-form>
	</n-modal>
</template>

<script lang="ts" setup>
import { AddFolder } from '@wailsjs/go/main/App';
import { FormInst } from 'naive-ui';
import { computed, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const show = defineModel<boolean>("show");
// This component is both used as the `custom difficult table` and 
// `favorite folder` module. Therefore, it requires a `type` field
// to display different contents 
const props = defineProps<{
	customTableId?: number,
	type: "table" | "folder"
}>();
const emit = defineEmits<{
	(e: 'refresh'): void
}>();

// Dynamic contents
const title = computed((): string => {
	if (props.type == "table") {
		return t('modal.title.table')
	} else if (props.type == "folder") {
		return t('modal.title.folder')
	}
})

const { t } = useI18n();

const loading = ref(false);
const formRef = ref<FormInst | null>(null);
const formData = reactive({
	FolderName: "",
});
const rules = {
	name: {
		required: true,
		message: t('message.missingName'),
		trigger: ["input", "blur"],
	},
};

function handlePositiveClick(): boolean {
	loading.value = true;
	formRef.value
		?.validate()
		.then(async () => {
			const result = await AddFolder({
				...formData,
				CustomTableID: props.customTableId ?? 1
			} as any);
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			show.value = false;
			formData.FolderName = null;
			emit('refresh');
		})
		.catch(err => window.$notifyError(err))
		.finally(() => loading.value = false);
	return false;
}

function handleNegativeClick() {
	formData.FolderName = "";
	show.value = false;
}
</script>

<i18n lang="json">{
	"en": {
		"modal": {
			"title": {
				"folder": "New Folder",
				"table": "New Difficult"
			},
			"positiveText": "Submit",
			"negativeText": "Cancel"
		},
		"form": {
			"labelName": "Name",
			"placeholderName": "Input name"
		},
		"message": {
			"missingName": "Please input name"
		}
	},
	"zh-CN": {
		"modal": {
			"title": {
				"folder": "新增收藏夹",
				"table": "新增难度"
			},
			"positiveText": "提交",
			"negativeText": "取消"
		},
		"form": {
			"labelName": "名称",
			"placeholderName": "请输入名称"
		},
		"message": {
			"missingName": "请输入名称"
		}
	}
}</i18n>
