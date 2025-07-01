<template>
	<n-modal :loading="loading" v-model:show="show" preset="dialog" :title="title"
		:positive-text="t('button.submit')" :negative-text="t('button.cancel')"
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
		return t('title.newDifficult')
	} else if (props.type == "folder") {
		return t('title.newFolder')
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
