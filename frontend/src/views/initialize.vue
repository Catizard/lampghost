<template>
	<div style="padding: 20px 40px; min-height: 85vh">
		<n-card>
			<n-space justify="space-between" vertical style="width: 65%">
				<n-h1 prefix="bar" style="text-align: start">
					<n-text type="primary">{{ t('title') }}</n-text>
				</n-h1>
				<n-form ref="formRef" :rules="rules" :model="formData">
					<n-form-item :label="t('form.labelLanguage')">
						<n-select v-model:value="locale" :options="localeOptions" style="width: 150px;" />
					</n-form-item>
					<n-form-item :label="t('form.labelUserName')" path="userName">
						<n-input v-model:value="formData.userName" :placeholder="t('form.placeholderUserName')" />
					</n-form-item>
					<n-form-item :label="t('form.labelSongdataPath')" path="songdataPath">
						<n-input v-model:value="formData.songdataPath" :placeholder="t('form.placeholderSongdataPath')" />
					</n-form-item>
					<n-form-item :label="t('form.labelScorelogPath')" path="scorelogPath">
						<n-input v-model:value="formData.scorelogPath" :placeholder="t('form.placeholderScorelogPath')" />
					</n-form-item>
					<n-form-item>
						<n-button @loading="loading" attr-type="button" @click="handleSubmit" type="primary">
							{{ t('button.submit') }}
						</n-button>
					</n-form-item>
				</n-form>
			</n-space>
		</n-card>
	</div>
</template>

<script setup>
import { useNotification } from "naive-ui";
import { ref, watch } from "vue";
import { InitializeMainUser } from "@wailsjs/go/controller/RivalInfoController";
import router from "../router";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const globalI18n = useI18n({ useScope: 'global' });
const notification = useNotification();

const loading = ref(false);
const formRef = ref(null);
const formData = ref({
	locale: null,
	songdataPath: "",
	scorelogPath: "",
	userName: "",
});

const localeOptions = [
	{
		label: "English",
		value: "en",
	},
	{
		label: "中文",
		value: "zh-CN"
	}
];
// I have no idea why need this hack
const locale = globalI18n.locale;
watch(locale, (newLocale) => {
	globalI18n.locale = ref(newLocale);
})

const rules = {
	userName: {
		required: true,
		message: t('rule.missingUserName'),
		trigger: ["input", "blur"],
	},
	songdataPath: {
		required: true,
		message: t('rule.missingSongdataPath'),
		trigger: ["input", "blur"],
	},
	scorelogPath: {
		required: true,
		message: t('rule.missingScorelogPath'),
		trigger: ["input", "blur"],
	},
};

function handleSubmit(e) {
	e.preventDefault();
	formRef.value?.validate((errors) => {
		if (errors) {
			console.error(errors);
		} else {
			loading.value = true;
			const rivalInfo = {
				Name: formData.value.userName,
				ScoreLogPath: formData.value.scorelogPath,
				SongDataPath: formData.value.songdataPath,
				Locale: globalI18n.locale.value,
			};
			InitializeMainUser(rivalInfo)
				.then((result) => {
					if (result.Code != 200) {
						return Promise.reject(result.Msg);
					}
					router.push("/");
				})
				.catch((err) => {
					notification.error({
						content: err,
						duration: 3000,
					});
				}).finally(() => loading.value = false);
		}
	});
}
</script>

<i18n lang="json">{
	"en": {
		"title": "Initialize User Data",
		"form": {
			"labelLanguage": "Language",
			"labelUserName": "Your user name",
			"labelSongdataPath": "songdata.db file path",
			"labelScorelogPath": "scorelog.db file path",
			"placeholderUserName": "Please input your name",
			"placeholderSongdataPath": "Please input songdata.db file path",
			"placeholderScorelogPath": "Please input scorelog.db file path"
		},
		"button": {
			"submit": "submit"
		},
		"rule": {
			"missingUserName": "Name cannot be empty",
			"missingSongdataPath": "songdata.db file path cannot be empty",
			"missingScorelogPath": "scorelog.db file path cannot be empty"
		}
	},
	"zh-CN": {
		"title": "初始化用户信息",
		"form": {
			"labelLanguage": "语言",
			"labelUserName": "用户名称",
			"labelSongdataPath": "songdata.db文件路径",
			"labelScorelogPath": "scorelog.db文件路径",
			"placeholderUserName": "请输入用户名",
			"placeholderSongdataPath": "请输入songdata.db文件路径",
			"placeholderScorelogPath": "请输入scorelog.db文件路径"
		},
		"button": {
			"submit": "提交"
		},
		"rule": {
			"missingUserName": "用户名不可为空",
			"missingSongdataPath": "songdata.db文件路径不可为空",
			"missingScorelogPath": "scorelog.db文件路径不可为空"
		}
	}
}</i18n>
