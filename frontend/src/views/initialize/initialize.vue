<template>
  <div style="padding: 20px 40px; min-height: 85vh">
    <n-card>
      <n-space vertical>
        <n-steps :current="current">
          <n-step :title="$t('title.initUser')" :description="$t('hint.initUser')" />
          <n-step :title="$t('title.initTable')" :description="$t('hint.initTable')" />
        </n-steps>
      </n-space>
      <div id="container" style="margin-top: 10px;">
        <InitUserForm v-if="current === 1" :moveOn="moveOn" />
        <InitTableForm v-if="current === 2" :moveOn="moveOn" />
      </div>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
import InitUserForm from "./initUserForm.vue";
import InitTableForm from "./initTableForm.vue";
import { ref, watch } from "vue";
import router from "@/router";

const current = ref(1);

function moveOn() {
  current.value += 1;
}

watch(current, (newValue) => {
  if (newValue == 3) {
    router.push("/home");
  }
});
</script>
