<script setup lang="ts">
import { QueryMainUser, QueryMetaInfo, ReadConfig } from '@wailsjs/go/main/App';
import { Provider, Viewer } from './components';
import { EventsOn } from '@wailsjs/runtime/runtime';
import { useI18n } from 'vue-i18n';
import { dto, entity } from '@wailsjs/go/models';
import { onMounted } from 'vue';
import router from './router';
import { useUserStore } from './stores/user';

let { t, locale } = useI18n();

ReadConfig()
  .then(result => {
    if (result.Code != 200) {
      // Should we report this?
      return;
    }
    locale.value = result.Data.Locale;
  });

// Global refresh
EventsOn("global:refresh", () => {
  window.location.reload();
});

const userStore = useUserStore();
onMounted(async () => {
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
</script>

<template>
  <Provider>
    <Viewer />
  </Provider>
</template>
