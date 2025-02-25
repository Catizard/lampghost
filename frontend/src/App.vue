<script setup>
import { ReadConfig } from '@wailsjs/go/controller/ConfigController';
import { Provider, Viewer } from './components';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();

ReadConfig()
  .then(result => {
    if (result.Code != 200) {
      // Should we report this?
      return;
    }
    i18n.locale = result.Data.Locale
    console.log('i18n.locale = ', i18n.locale);
  })
</script>

<template>
  <Provider>
    <Viewer />
  </Provider>
</template>