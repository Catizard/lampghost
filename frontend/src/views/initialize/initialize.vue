<template>
  <div style="padding: 20px 40px; min-height: 85vh">
    <n-card>
      <n-space vertical>
        <n-steps :current="current">
          <n-step :title="t('title.initUser')" :description="t('desc.initUser')" />
          <n-step :title="t('title.initTable')" :description="t('desc.initTable')" />
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
import { useI18n } from "vue-i18n";
import InitUserForm from "./initUserForm.vue";
import InitTableForm from "./initTableForm.vue";
import { ref, watch } from "vue";
import router from "@/router";

const current = ref(1);

function moveOn() {
  current.value += 1;
}

const { t } = useI18n();

watch(current, (newValue) => {
  if (newValue == 3) {
    router.push("/home");
  }
});
</script>

<i18n lang="json">{
  "en": {
    "title": {
      "initUser": "Initialize User",
      "initTable": "Initialize Table"
    },
    "desc": {
      "initUser": "Configure your own save files location",
      "initTable": "Add some most used tables"
    }
  },
  "zh-CN": {
    "title": {
      "initUser": "初始化用户",
      "initTable": "初始化难度表"
    },
    "desc": {
      "initUser": "配置你的存档文件路径",
      "initTable": "添加一些最常用的难度表"
    }
  }
}</i18n>
