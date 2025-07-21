<template>
  <n-modal :loading="loading" v-model:show="show" preset="dialog" :title="t('title.preview')"
    :positive-text="t('button.copy')" :negative-text="t('button.cancel')" @positive-click="handlePositiveClick"
    @negative-click="handleNegativeClick" :mask-closable="false" style="width: 85vw;">
    <TimelineContent ref="timelineContentRef" :nodes="timelineNodes" />
  </n-modal>
</template>

<script setup lang="ts">
import { QueryPrevDayScoreLogList, CopyImageToClipboard } from '@wailsjs/go/main/App';
import html2canvas from 'html2canvas';
import { Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import TimelineContent from './TimelineContent.vue';
import { Node, parseOneDayLogs } from './node';
import dayjs from 'dayjs';

const show = defineModel<boolean>("show");
const props = defineProps<{
  rivalId?: number
}>()

const { t } = useI18n();
const loading = ref(false);
const timelineNodes: Ref<Node[]> = ref([]);
const timelineContentRef = ref<InstanceType<typeof TimelineContent>>(null);

async function loadData() {
  loading.value = true;
  try {
    const result = await QueryPrevDayScoreLogList({
      EndRecordTimestamp: dayjs().add(1, 'day').unix(),
      RivalId: props.rivalId,
    } as any);
    if (result.Code != 200) {
      throw result.Msg;
    }
    if (result.Rows.length == 0) {
      window.$notifyInfo("message.noRecordError");
    } else {
      const [_, nextNodes] = parseOneDayLogs(result.Rows);
      timelineNodes.value = [...nextNodes];
    }
  } catch (err) {
    window.$notifyError(err)
  } finally {
    loading.value = false;
  }
}

async function handlePositiveClick() {
  const el: HTMLElement = timelineContentRef.value.$el;
  try {
    const canvas = await html2canvas(el);
    try {
      const result = await CopyImageToClipboard(canvas.toDataURL("image/png"));
      if (result.Code != 200) {
        throw result.Msg;
      }
      window.$notifyInfo(t("message.setClipboardSuccess"));
    } catch (err) {
      window.$notifyError(t("message.setClipboardError", err));
    }
  } catch (err) {
    window.$notifyError("cannot convert html to canvas: " + err);
  } finally {
    show.value = false;
  }
}

function handleNegativeClick() {
  show.value = false;
}

watch(show, (newValue) => {
  if (newValue != true) {
    return;
  }
  loadData();
});

</script>
