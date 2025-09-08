<script setup lang="ts">
import { QueryMainUser, QueryMetaInfo, ReadConfig } from '@wailsjs/go/main/App';
import { Provider, Viewer } from './components';
import { EventsOn } from '@wailsjs/runtime/runtime';
import { useI18n } from 'vue-i18n';
import { entity } from '@wailsjs/go/models';
import { onMounted } from 'vue';
import router from './router';
import { useUserStore } from './stores/user';
import { useConfigStore } from './stores/config';

let { t, locale } = useI18n();

// Global refresh
EventsOn("global:refresh", () => {
  window.location.reload();
});

// Update config
EventsOn("config:update", () => {
  readConfig();
});

const userStore = useUserStore();
const configStore = useConfigStore();
onMounted(async () => {
  readConfig();

  try {
    const result = await QueryMainUser();
    if (result.Code != 200) {
      throw t("message.noMainUserError")
    }
    userStore.setter(result.Data);
  } catch (err) {
    window.$notifyError(err);
    router.push("/initialize");
    // This return is ok since the below checks are not necessary if initialization hasn't been done
    return;
  }

  try {
    const result = await QueryMetaInfo();
    if (result.Code != 200) {
      throw result.Msg;
    }
    const data: entity.MetaInfo = result.Data;
    if (data.ReleaseVersion != data.CurrentVersion) {
      window.$notifyInfo(t("message.releaseVersionDiverage", {
        releaseVersion: data.ReleaseVersion,
        currentVersion: data.CurrentVersion
      }));
    }
    if (!data.ClipboardSetup) {
      window.$notifyError(t("message.clipboardForImagesAreNotSetup"));
    }
  } catch (err) {
    window.$notifyError(err);
  }
});

async function readConfig() {
  try {
    const result = await ReadConfig();
    if (result.Code != 200) {
      throw result.Msg;
    }
    locale.value = result.Data.Locale;
    configStore.setter(result.Data);
  } catch (err) {
    window.$notifyError(err);
    // Something might be broken, but the possiblity is rare
  }
}
</script>

<template>
  <Provider>
    <Viewer />
  </Provider>
</template>
