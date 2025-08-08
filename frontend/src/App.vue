<script setup lang="ts">
import { QueryMainUser, ReadConfig } from '@wailsjs/go/main/App';
import { Provider, Viewer } from './components';
import { EventsOn } from '@wailsjs/runtime/runtime';
import { useI18n } from 'vue-i18n';
import { dto } from '@wailsjs/go/models';
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
    return;
  }
});
</script>

<template>
  <Provider>
    <Viewer />
  </Provider>
</template>
