<template>
  <n-config-provider :hljs="hljs">
    <n-notification-provider>
      <n-modal-provider>
        <n-global-style />
        <n-loading-bar-provider>
          <n-message-provider>
            <n-dialog-provider>
              <slot />
            </n-dialog-provider>
          </n-message-provider>
        </n-loading-bar-provider>
      </n-modal-provider>
    </n-notification-provider>
  </n-config-provider>
</template>

<script lang="ts" setup>
import { darkTheme, useOsTheme } from "naive-ui";
import { computed } from "vue";
import hljs from "highlight.js/lib/core";
import json from "highlight.js/lib/languages/json"
import javascript from "highlight.js/lib/languages/javascript";

hljs.registerLanguage('json', json);
hljs.registerLanguage('javascript', javascript)

const osTheme = useOsTheme();
const theme = computed(() => (osTheme.value == "dark" ? darkTheme : null));
</script>
