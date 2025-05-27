<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: left;">
      <n-text type="primary"> {{ t('title') }}</n-text>
    </n-h1>
  </n-flex>
  <n-data-table :loading="loading" :columns="columns" :data="data" :row-key="(row: entity.DownloadTask) => row.ID" />
</template>

<script lang="ts" setup>
import { entity } from '@wailsjs/go/models';
import { EventsOn } from '@wailsjs/runtime/runtime';
import { DataTableColumns, NEllipsis, tooltipProps } from 'naive-ui';
import { h, onMounted, onUnmounted, ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import TaskStatusTag from './TaskStatusTag.vue';
import { DownloadTaskStatus } from '@/constants/downloadTaskStatus';

const { t } = useI18n();

const loading = ref(false);

let cancel: () => void = null;
const data: Ref<entity.DownloadTask[]> = ref([]);

onMounted(() => {
  cancel = EventsOn("DownloadTask:pushup", ((tasks: entity.DownloadTask[]) => {
    data.value = [...tasks];
  }));
});

onUnmounted(() => {
  if (cancel != null) {
    cancel();
  }
});


const columns: DataTableColumns<entity.DownloadTask> = [
  {
    title: t('column.taskName'),
    key: "TaskName",
    render(row: entity.DownloadTask) {
      return h(
        'div',
        {
          onClick: () => {
            navigator.clipboard.writeText(row.URL);
            console.log('wow, your url is: ', row.URL)
          },
        },
        [
          h(
            NEllipsis,
            { tooltip: true },
            { default: () => row.TaskName }
          )
        ]
      )
    }
  },
  {
    title: t('column.progress'), key: "progress",
    render(row: entity.DownloadTask) {
      return `${humanFileSize(row.DownloadSize, true)}/${humanFileSize(row.ContentLength, true)}(${(row.DownloadSize / row.ContentLength * 100).toFixed(2)}%)`;
    }
  },
  {
    title: t('column.status'), key: "status",
    render(row: entity.DownloadTask) {
      return h(
        TaskStatusTag,
        { status: DownloadTaskStatus.from(row.Status), errorMsg: row.ErrorMessage },
      )
    }
  }
];

/**
 * Format bytes as human-readable text.
 * https://stackoverflow.com/questions/10420352/converting-file-size-in-bytes-to-human-readable-string
 * 
 * @param bytes Number of bytes.
 * @param si True to use metric (SI) units, aka powers of 1000. False to use 
 *           binary (IEC), aka powers of 1024.
 * @param dp Number of decimal places to display.
 * 
 * @return Formatted string.
 */
function humanFileSize(bytes, si = false, dp = 1) {
  const thresh = si ? 1000 : 1024;

  if (Math.abs(bytes) < thresh) {
    return bytes + ' B';
  }

  const units = si
    ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
    : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  let u = -1;
  const r = 10 ** dp;

  do {
    bytes /= thresh;
    ++u;
  } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


  return bytes.toFixed(dp) + ' ' + units[u];
}
</script>

<i18n lang="json">{
  "en": {
    "title": "Download Tasks",
    "column": {
      "url": "URL",
      "taskName": "Task Name(Click to Copy)",
      "progress": "Progress"
    }
  },
  "zh-CN": {
    "title": "下载任务(点击复制)",
    "column": {
      "url": "URL",
      "taskName": "任务名称",
      "progress": "当前进度"
    }
  }
}</i18n>