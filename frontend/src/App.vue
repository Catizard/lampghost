<script setup lang="ts">
import { ReadConfig } from '@wailsjs/go/main/App';
import { Provider, Viewer } from './components';
import { EventsOn } from '@wailsjs/runtime/runtime';
import { useI18n } from 'vue-i18n';

let { locale } = useI18n();

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
})
</script>

<template>
  <Provider>
    <Viewer />
  </Provider>
</template>
