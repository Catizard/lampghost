<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: left;">
      <n-text type="primary"> {{ t('title.downloadTask') }}</n-text>
    </n-h1>
  </n-flex>
  <n-flex justify="space-between">
    <n-radio-group v-model:value="status" name="downloadStatusRadioGroup">
      <n-radio-button key="progressing" value="progressing" :label="t('button.progressing')" />
      <n-radio-button key="completed" value="completed" :label="t('button.completed')" />
    </n-radio-group>
    {{ downloadProgress }}
    <n-button type="warning" size="small" :disabled="!hasErrorTasks" @click="handleQuickRetry" style="margin-left: 8px;">
      {{ t('button.quickRetry') }}
    </n-button>
  </n-flex>
  <n-data-table :loading="loading" :columns="columns" :data='status == "progressing" ? progressing : completed'
    :row-key="(row: entity.DownloadTask) => row.ID" :pagination="pagination" />
</template>

<script lang="ts" setup>
import { entity } from '@wailsjs/go/models';
import { ClipboardSetText, EventsOn } from '@wailsjs/runtime/runtime';
import { DataTableColumns, NButton, NDropdown, NEllipsis } from 'naive-ui';
import { computed, h, onMounted, onUnmounted, reactive, ref, Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import TaskStatusTag from './TaskStatusTag.vue';
import { DownloadTaskStatus } from '@/constants/downloadTaskStatus';
import { CancelDownloadTask, RestartDownloadTask } from '@wailsjs/go/main/App';
import { Download } from '@vicons/ionicons5';

const { t } = useI18n();

const loading = ref(false);

let cancel: () => void = null;
const progressing: Ref<entity.DownloadTask[]> = ref([]);
const completed: Ref<entity.DownloadTask[]> = ref([]);
const status = ref<"progressing" | "completed">("progressing");

onMounted(() => {
  cancel = EventsOn("DownloadTask:pushup", ((tasks: entity.DownloadTask[]) => {
    console.log(tasks);
    tasks.sort((a: entity.DownloadTask, b: entity.DownloadTask): number => {
      console.log('a: ', a, 'b: ', b);
      return DownloadTaskStatus.compare(
        DownloadTaskStatus.from(a.Status),
        DownloadTaskStatus.from(b.Status)
      );
    });
    const isSuccess = (task: entity.DownloadTask) => {
      console.log(task, task.Status);
      return DownloadTaskStatus.from(task.Status) == DownloadTaskStatus.SUCCESS;
    };
    console.log('1');
    completed.value = [...tasks.filter(task => isSuccess(task))];
    console.log('2');
    progressing.value = [...tasks.filter(task => !isSuccess(task))];
  }));
});

onUnmounted(() => {
  if (cancel != null) {
    cancel();
  }
});

const downloadProgress = computed(() => {
  const c = completed.value.length, p = progressing.value.length;
  return t('message.downloadProgress', {
    success: c,
    count: c + p,
    progress: ((c + p) == 0 ? 100 : c / (c + p) * 100).toFixed(2)
  });;
});

const columns: DataTableColumns<entity.DownloadTask> = [
  {
    title: t('column.taskName'), key: "TaskName", ellipsis: { tooltip: true },
    render(row: entity.DownloadTask) {
      return h(
        'div',
        {
          onClick: () => {
            ClipboardSetText(row.URL);
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
    title: t('column.progress'), key: "progress", width: "220px",
    render(row: entity.DownloadTask) {
      return `${humanFileSize(row.DownloadSize, true)}/${humanFileSize(row.ContentLength, true)}(${(row.DownloadSize / row.ContentLength * 100).toFixed(2)}%)`;
    }
  },
  {
    title: t('column.status'), key: "status", width: "125px",
    render(row: entity.DownloadTask) {
      console.log(row.Status);
      return h(
        TaskStatusTag,
        { status: DownloadTaskStatus.from(row.Status), errorMsg: row.ErrorMessage },
      )
    }
  },
  {
    title: t('column.actions'), key: "actions", width: "75px",
    render(row: entity.DownloadTask) {
      return h(
        NDropdown,
        {
          trigger: "hover",
          options: [
            { label: t('button.cancel'), key: "cancel", disabled: !DownloadTaskStatus.cancelable(row.Status) },
            { label: t('button.restart'), key: "restart", disabled: !DownloadTaskStatus.restartable(row.Status) },
          ],
          onSelect: (key: "cancel" | "restart") => {
            switch (key) {
              case "restart": restartDownloadTask(row.ID); break;
              case 'cancel': cancelDownloadTask(row.ID); break;
            }
          }
        },
        { default: () => h(NButton, null, { default: () => "..." }) }
      )
    }
  }
];
const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
  }
});

function restartDownloadTask(ID: number) {
  RestartDownloadTask(ID)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
    }).catch(err => window.$notifyError(err));
}

function cancelDownloadTask(ID: number) {
  CancelDownloadTask(ID)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
    }).catch(err => window.$notifyError(err));
}

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

const hasErrorTasks = computed(() => {
  return progressing.value.some(task => DownloadTaskStatus.from(task.Status) == DownloadTaskStatus.ERROR);
});

async function handleQuickRetry() {
  const errorTasks = progressing.value.filter(task => DownloadTaskStatus.from(task.Status) == DownloadTaskStatus.ERROR);
  if (errorTasks.length === 0) {
    window.$notifySuccess(t('message.noErrorTasks'));
    return;
  }
  try {
    const results = await Promise.all(errorTasks.map(task => RestartDownloadTask(task.ID)));
    const failed = results.filter(result => result.Code != 200);
    if (failed.length > 0) {
      window.$notifyError(failed.map(f => f.Msg).join('; '));
    } else {
      window.$notifySuccess(t('message.quickRetryStarted'));
    }
  } catch (err) {
    window.$notifyError(err);
  }
}
</script>
