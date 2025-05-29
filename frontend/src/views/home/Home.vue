<template>
	<n-h1 prefix="bar" style="text-align: start">
		<n-text type="primary">
			{{ t("infoTitle") }}
		</n-text>
	</n-h1>
	<n-grid :cols="12">
		<n-gi :span="4">
			<n-flex>
				{{ t("playerInfo.name") }}: {{ playerData.playerName }}
				<n-divider />
				{{ t("playerInfo.count") }}: {{ playerData.playCount }}
				<n-divider />
				{{ t("playerInfo.lastSyncTime") }}: {{ playerData.lastUpdate }}
				<n-divider />
				<n-button @click="handleSyncClick" :loading="syncLoading">
					{{ t("button.sync") }}
				</n-button>
			</n-flex>
		</n-gi>
		<n-gi :span="8">
			<PlayCountChart :rival-id="currentUser?.ID" />
		</n-gi>
	</n-grid>
	<n-divider />
	<n-h1 prefix="bar" style="text-align: start">
		<n-text type="primary">
			{{ t("lampStatusTitle") }}
		</n-text>
	</n-h1>
	<LampCountChart :rival-id="currentUser?.ID" />
	<n-divider />
	<n-h1 prefix="bar" style="text-align: start;">
		<n-text type="primary">
			{{ t("keyCountTitle") }}
		</n-text>
	</n-h1>
	<KeyCountChart :rival-id="currentUser?.ID" />
	<n-h1 prefix="bar" style="text-align: start;">
		<n-text type="primary">
			{{ t("timelineTitle") }}
		</n-text>
	</n-h1>
	<RecentTimeline :rival-id="currentUser?.ID" />
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { ReloadRivalData, QueryMainUser, QueryUserInfoByID } from "@wailsjs/go/main/App";
import { useNotification } from "naive-ui";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";
import PlayCountChart from "./PlayCountChart.vue";
import LampCountChart from "./LampCountChart.vue";
import { dto } from '@wailsjs/go/models';
import { useRoute, useRouter } from "vue-router";
import RecentTimeline from "./RecentTimeline.vue";
import KeyCountChart from "./KeyCountChart.vue";

const i18n = useI18n();
const { t } = i18n;
const router = useRouter();
const route = useRoute();

const currentUser = ref(null);
const playerData = reactive({
	playerName: "U",
	playCount: 0,
	lastUpdate: "",
});
const currentRivalID = ref(null);

function initUser() {
	QueryMainUser()
		.then((result) => {
			if (result.Code != 200) {
				router.push("/initialize");
				return Promise.reject(t("message.noMainUserError"));
			}
			return result.Data;
		}).then((mainUserData: dto.RivalInfoDto) => {
			currentRivalID.value = route?.query?.rivalID ?? 1;
			if (currentRivalID.value == 1) {
				return mainUserData;
			} else {
				return QueryUserInfoByID(parseInt(currentRivalID.value))
					.then(result => {
						if (result.Code != 200) {
							return Promise.reject(result.Msg);
						}
						return result.Data;
					});
			}
		}).then((data: dto.RivalInfoDto) => {
			currentUser.value = data;
			playerData.playerName = data.Name;
			playerData.playCount = data.PlayCount;
			playerData.lastUpdate = dayjs(data.UpdatedAt).format("YYYY-MM-DD HH:mm:ss");
		}).catch(err => {
			window.$notifyError(err);
		});
}

const syncLoading = ref(false);
function handleSyncClick() {
	syncLoading.value = true;
	ReloadRivalData(currentUser.value.ID, false)
		.then(result => {
			if (result.Code != 200) {
				return Promise.reject();
			}
			syncLoading.value = false;
			window.$notifySuccess(t("message.reloadSuccess"));
			// Reload data by simply reloading this page
			window.location.reload();
		})
		.catch((err) => {
			window.$notifyError(t("message.reloadFailedPrefix") + err);
			syncLoading.value = false;
		});
}

initUser();
</script>

<i18n lang="json">{
	"en": {
		"infoTitle": "Player Info",
		"playerInfo": {
			"name": "Player Name",
			"count": "Player Count",
			"lastSyncTime": "Last Sync Time"
		},
		"button": {
			"sync": "Reload Save File",
			"chooseYear": "Choose Year"
		},
		"lampStatusTitle": "Lamp Status",
		"keyCountTitle": "Key Count",
		"timelineTitle": "Recent Activity",
		"message": {
			"noMainUserError": "Found no main user, please first load your save file in",
			"loadUserDataErrorPrefix": "Cannot load user data: ",
			"reloadSuccess": "Successfully reloaded",
			"reloadFailedPrefix": "Failed to load save file, error message: "
		}
	},
	"zh-CN": {
		"infoTitle": "玩家信息",
		"playerInfo": {
			"name": "玩家名称",
			"count": "游玩次数",
			"lastSyncTime": "最后同步时间"
		},
		"button": {
			"sync": "同步最新存档",
			"chooseYear": "选择年份"
		},
		"lampStatusTitle": "点灯情况",
		"keyCountTitle": "按键次数",
		"timelineTitle": "最近游玩",
		"message": {
			"noMainUserError": "找不到主用户信息，请先导入你自己的存档",
			"loadUserDataErrorPrefix": "获取用户信息失败: ",
			"reloadSuccess": "同步成功",
			"reloadFailedPrefix": "同步失败，返回结果: "
		}
	}
}</i18n>

<style scoped>
.n-button {
	width: 80%;
	white-space: normal;
}
</style>
