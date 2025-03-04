<script setup>
import { ReadConfig } from '@wailsjs/go/controller/ConfigController';
import { Provider, Viewer } from './components';
import { useI18n } from 'vue-i18n';
import { ref } from 'vue';

const i18n = useI18n();

ReadConfig()
  .then(result => {
    if (result.Code != 200) {
      // Should we report this?
      return;
    }
    console.log(result);
    i18n.locale = ref(result.Data.Locale);
  })
</script>

<template>
  <Provider>
    <Viewer />
  </Provider>
</template>