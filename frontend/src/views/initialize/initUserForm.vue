<template>
	<n-space justify="space-between" vertical style="width: 65%">
		<!-- <n-h1 prefix="bar" style="text-align: start">
      <n-text type="primary">{{ t('title') }}</n-text>
    </n-h1> -->
		<n-form ref="formRef" :rules="rules" :model="formData" :loading="loading">
			<n-form-item :label="t('form.labelLanguage')">
				<n-select v-model:value="locale" :options="localeOptions" style="width: 150px;" />
			</n-form-item>
			<n-form-item :label="t('form.labelUserName')" path="name">
				<n-input v-model:value="formData.name" :placeholder="t('form.placeholderUserName')" />
			</n-form-item>
			<n-radio-group v-model:value="importStrategy" name="radiogroup">
				<n-radio-button key="directory" value="directory" :label="t('import.labelDirectory')" />
				<n-radio-button key="separate" value="separate" :label="t('import.labelSeparate')" />
			</n-radio-group>
			<template v-if="importStrategy == 'directory'">
				<div style="margin-top: 10px;">
					<n-button type="info" @click="chooseBeatorajaDirectory()">
						{{ t('button.chooseBeatorajaDirectory') }}
					</n-button>
					<n-select :options="playerDirectories" v-model:value="formData.playerDirectory"
						style="margin-top: 10px;width: 250px;">
					</n-select>
				</div>
			</template>
			<template v-if="importStrategy == 'separate'">
				<n-form-item :label="t('form.labelSongdataPath')" path="songdataPath" style="margin-top: 10px;">
					<n-button type="info" @click="chooseFile('Choose songdata.db', 'songdataPath')">
						{{ t('button.chooseFile') }}
					</n-button>
					<n-divider vertical />
					<n-input v-model:value="formData.songdataPath" :placeholder="t('form.placeholderSongdataPath')" />
				</n-form-item>
				<n-form-item :label="t('form.labelScorelogPath')" path="scorelogPath">
					<n-button type="info" @click="chooseFile('Choose scorelog.db', 'scorelogPath')">
						{{ t('button.chooseFile') }}
					</n-button>
					<n-divider vertical />
					<n-input v-model:value="formData.scorelogPath" :placeholder="t('form.placeholderScorelogPath')" />
				</n-form-item>
				<n-form-item :label="t('form.labelScoredatalogPath')" path="scoredatalogPath">
					<n-button type="info" @click="chooseFile('Choose scoredatalog.db', 'scoredatalogPath')">
						{{ t('button.chooseFile') }}
					</n-button>
					<n-divider vertical />
					<n-input v-model:value="formData.scoredatalogPath" :placeholder="t('form.placeholderScoredatalogPath')" />
				</n-form-item>
			</template>
			<n-form-item>
				<n-button :loading="loading" attr-type="button" @click="handleSubmit" type="primary">
					{{ t('button.submit') }}
				</n-button>
			</n-form-item>
		</n-form>
	</n-space>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { ref, watch, reactive, Ref } from "vue";
import { ChooseBeatorajaDirectory, InitializeMainUser, OpenFileDialog } from "@wailsjs/go/main/App";
import { SelectOption } from 'naive-ui';
import { dto } from '@wailsjs/go/models';

const { t } = useI18n();
const globalI18n = useI18n({ useScope: 'global' });

const props = defineProps<{
	moveOn: () => void
}>();

const importStrategy: Ref<"directory" | "separate"> = ref("directory");

// target == "scorelogPath" | "songdataPath" | "scoredatalogPath"
function chooseFile(title: string, target: "scorelogPath" | "songdataPath" | "scoredatalogPath") {
	OpenFileDialog(title)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			if (result.Data != null && result.Data != undefined && result.Data != "") {
				if (target == "scorelogPath") {
					formData.scorelogPath = result.Data;
				} else if (target == "songdataPath") {
					formData.songdataPath = result.Data;
				} else if (target == "scoredatalogPath") {
					formData.scoredatalogPath = result.Data;
				}
			}
		}).catch(err => window.$notifyError(err));
}

const playerDirectories: Ref<SelectOption[]> = ref([]);
function chooseBeatorajaDirectory() {
	ChooseBeatorajaDirectory()
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg)
			}
			const data: dto.BeatorajaDirectoryMeta = result.Data;
			playerDirectories.value = data.PlayerDirectories.map(row => {
				return {
					label: row,
					value: row
				} as SelectOption
			});
			formData.beatorajaDirectoryPath = data.BeatorajaDirectoryPath;
			formData.playerDirectory = playerDirectories.value[0].value as string;
		}).catch(err => window.$notifyError(err));
}

const loading = ref(false);
const formRef = ref(null);
const formData = reactive({
	locale: null,
	songdataPath: "",
	scorelogPath: "",
	scoredatalogPath: "",
	name: "",
	beatorajaDirectoryPath: "",
	playerDirectory: "",
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
	(globalI18n as any).locale = ref(newLocale);
})

const rules = {
	name: {
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
	scoredatalogPath: {
		required: true,
		message: t('rule.missingScoredatalogPath'),
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
			// TODO: This was a mistake...
			const rivalInfo = {
				Name: formData.name,
				ScoreLogPath: formData.scorelogPath,
				SongDataPath: formData.songdataPath,
				ScoreDataLogPath: formData.scoredatalogPath,
				BeatorajaDirectoryPath: formData.beatorajaDirectoryPath,
				PlayerDirectory: formData.playerDirectory,
				ImportStrategy: importStrategy.value,
				Locale: globalI18n.locale.value,
			};
			InitializeMainUser(rivalInfo as any)
				.then((result) => {
					if (result.Code != 200) {
						return Promise.reject(result.Msg);
					}
					props.moveOn();
				})
				.catch(err => window.$notifyError(err))
				.finally(() => loading.value = false);
		}
	});
}
</script>

<i18n lang="json">{
	"en": {
		"title": "Initialize User Data",
		"import": {
			"labelDirectory": "Beatoraja Directory",
			"labelSeparate": "Separate Files"
		},
		"form": {
			"labelLanguage": "Language",
			"labelUserName": "Your user name",
			"labelSongdataPath": "songdata.db file path",
			"labelScorelogPath": "scorelog.db file path",
			"labelScoredatalogPath": "scoredatalog.db file path",
			"placeholderUserName": "Please input your name",
			"placeholderSongdataPath": "Please input /beatoraja path/songdata.db file path",
			"placeholderScorelogPath": "Please input /beatoraja path/*Player Name*/scorelog.db file path",
			"placeholderScoredatalogPath": "Please input /beatoraja path/*Player Name*/scoredatalog.db file path"
		},
		"button": {
			"submit": "submit",
			"chooseFile": "Choose File",
			"chooseBeatorajaDirectory": "Choose Beatoraja Directory"
		},
		"rule": {
			"missingUserName": "Name cannot be empty",
			"missingSongdataPath": "songdata.db file path cannot be empty",
			"missingScorelogPath": "scorelog.db file path cannot be empty",
			"missingScoredatalogPath": "scoredatalog.db file path cannot be empty"
		}
	},
	"zh-CN": {
		"title": "初始化用户信息",
		"import": {
			"labelDirectory": "Beatoraja目录",
			"labelSeparate": "文件导入"
		},
		"form": {
			"labelLanguage": "语言",
			"labelUserName": "用户名称",
			"labelSongdataPath": "songdata.db文件路径",
			"labelScorelogPath": "scorelog.db文件路径",
			"labelScoredatalogPath": "scoredatalog.db文件路径",
			"placeholderUserName": "请输入用户名",
			"placeholderSongdataPath": "请输入/beatoraja path/songdata.db文件路径",
			"placeholderScorelogPath": "请输入/beatoraja path/*Player Name*/scorelog.db文件路径",
			"placeholderScoredatalogPath": "请输入/beatoraja path/*Player Name*/scoredatalog.db文件路径"
		},
		"button": {
			"submit": "提交",
			"chooseFile": "选择文件",
			"chooseBeatorajaDirectory": "选择Beatoraja路径"
		},
		"rule": {
			"missingUserName": "用户名不可为空",
			"missingSongdataPath": "songdata.db文件路径不可为空",
			"missingScorelogPath": "scorelog.db文件路径不可为空",
			"missingScoredatalogPath": "scoredatalog.db文件路径不可为空"
		}
	}
}</i18n>
