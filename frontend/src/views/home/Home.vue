<template>
  <n-h1 prefix="bar" style="text-align: start">
    <n-text type="primary">
      {{ t("title.playerInfo") }}
    </n-text>
  </n-h1>
  <n-grid :cols="12">
    <n-gi :span="4">
      <n-flex>
        {{ t("column.name") }}: {{ playerData.playerName }}
        <n-divider />
        {{ t("column.playCount") }}: {{ playerData.playCount }}
        <n-divider />
        {{ t("column.lastSyncTime") }}: {{ playerData.lastUpdate }}
        <n-divider />
        <n-button type="info" @click="handleSyncClick" :loading="syncLoading" style="width: 80%; white-space: normal;">
          {{ t("button.reloadSaveFile") }}
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
      {{ t("title.lampStatus") }}
    </n-text>
  </n-h1>
  <LampCountChart :rival-id="currentUser?.ID" />
  <n-divider />
  <n-h1 prefix="bar" style="text-align: start;">
    <n-text type="primary">
      {{ t("title.keyCount") }}
    </n-text>
  </n-h1>
  <KeyCountChart :rival-id="currentUser?.ID" />
  <n-h1 prefix="bar" style="text-align: start;">
    <n-flex justify="space-between">
      <n-text type="primary">
        {{ t("title.timeline") }}
      </n-text>
      <n-button type="info" @click="handleCopyTimeline" size="medium">
        {{ t('button.copyToClipboard') }}
      </n-button>
    </n-flex>
  </n-h1>
  <RecentTimeline :rival-id="currentUser?.ID" />
  <TimelinePreview v-model:show="showPreviewTimeline" :rivalId="currentUser?.ID" />
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ReloadRivalData, QueryMainUser, QueryUserInfoByID } from "@wailsjs/go/main/App";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";
import PlayCountChart from "./PlayCountChart.vue";
import LampCountChart from "./LampCountChart.vue";
import { dto } from '@wailsjs/go/models';
import { useRoute, useRouter } from "vue-router";
import RecentTimeline from "./recentTimeline/RecentTimeline.vue";
import KeyCountChart from "./KeyCountChart.vue";
import TimelinePreview from "./recentTimeline/Preview.vue";
import { useUserStore } from "@/stores/user";

const { t } = useI18n();
const router = useRouter();
const route = useRoute();

const userStore = useUserStore();
const currentUser = ref(null);
const playerData = reactive({
  playerName: "U",
  playCount: 0,
  lastUpdate: "",
});
const currentRivalID = ref(null);

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
      window.$notifyError(t("message.reloadError") + err);
      syncLoading.value = false;
    });
}

const showPreviewTimeline = ref(false);
function handleCopyTimeline() {
  showPreviewTimeline.value = true;
}

onMounted(async () => {
  let user = userStore.$state.user;
  currentRivalID.value = route?.query?.rivalID ?? user.ID;
  if (currentRivalID.value != 1) {
    try {
      const result = await QueryUserInfoByID(parseInt(currentRivalID.value));
      if (result.Code != 200) {
        throw result.Msg;
      }
      user = result.Data;
    } catch (err) {
      window.$notifyError(err);
    }
  }
  currentUser.value = user;
  playerData.playerName = user.Name;
  playerData.playCount = user.PlayCount;
  playerData.lastUpdate = dayjs(user.UpdatedAt).format("YYYY-MM-DD HH:mm:ss");
});
</script>
