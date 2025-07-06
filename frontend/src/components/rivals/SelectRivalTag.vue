<template>
  <n-select :loading="loading" v-model:value="value" :options="options" :style="{ width: width }"
    :placeholder="placeholder" :render-option="renderRivalTagOption" :clearable="clearable ?? false"/>
</template>

<script lang="ts" setup>
import { FindRivalTagList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { NTooltip, SelectOption } from 'naive-ui';
import { h, Ref, ref, VNode, watch } from 'vue';

const value = defineModel<number | null>("value", { default: null });
const props = defineProps<{
  rivalId: number | null,
  clearable?: boolean,
  placeholder?: string
  width: string
}>();

const loading = ref(false);
const options: Ref<SelectOption[]> = ref([]);

function renderRivalTagOption({ node, option }: { node: VNode, option: SelectOption }) {
  return h(NTooltip, {
    style: `max-width: ${props.width}; font-color: white`,
  }, {
    trigger: () => node,
    default: () => option.label
  });
}

watch(() => props.rivalId, (rivalId: number | null) => {
  if (rivalId == null) {
    return;
  }
  loading.value = true;
  FindRivalTagList({ RivalId: rivalId } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      options.value = result.Rows.map((row: dto.RivalTagDto) => {
        return {
          label: row.TagName,
          value: row.ID,
        } as SelectOption
      });
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}, { immediate: true});
</script>