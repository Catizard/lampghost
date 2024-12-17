<template>
  <n-config-provider>
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

const osTheme = useOsTheme();
const theme = computed(() => (osTheme.value == "dark" ? darkTheme : null));
</script>
