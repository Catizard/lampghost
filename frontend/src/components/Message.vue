<script lang="ts">
import { createDiscreteApi } from 'naive-ui';
import './message.d.ts';
import { defineComponent } from 'vue';
import { EventsOn } from '@wailsjs/runtime/runtime.js';
import { dto } from '@wailsjs/go/models.js';

const { message, dialog, notification, loadingBar, modal } = createDiscreteApi(
  ['message', 'dialog', 'notification', 'loadingBar', 'modal']
);


window.$message = message;
window.$dialog = dialog;
window.$notification = notification;
window.$loadingBar = loadingBar;
window.$modal = modal;
function createNotification(type: "info" | "default" | "error" | "warning" | "success", content: string, duration: number, keepAliveOnHover: boolean) {
  window.$notification.create({
    type: type,
    content: content,
    duration: duration,
    keepAliveOnHover: keepAliveOnHover,
  })
}
window.$notifyDefault = (content: string, duration: number = 3000, keepAliveOnHover: boolean = true) => {
  createNotification("default", content, duration, keepAliveOnHover);
};
window.$notifyInfo = (content: string, duration: number = 3000, keepAliveOnHover: boolean = true) => {
  createNotification("info", content, duration, keepAliveOnHover);
};
window.$notifyWarning = (content: string, duration: number = 3000, keepAliveOnHover: boolean = true) => {
  createNotification("warning", content, duration, keepAliveOnHover);
};
window.$notifySuccess = (content: string, duration: number = 3000, keepAliveOnHover: boolean = true) => {
  createNotification("success", content, duration, keepAliveOnHover);
};
window.$notifyError = (content: string | Error, duration: number = 3000, keepAliveOnHover: boolean = true) => {
  if (content instanceof Error) {
    content = String(content);
  }
  createNotification("error", content, duration, keepAliveOnHover);
};


// Global notification
EventsOn("global:notify", (notify: dto.NotificationDto) => {
  createNotification(notify.type as any, notify.content, 3000, true);
});

export default defineComponent({});
</script>


<template></template>
