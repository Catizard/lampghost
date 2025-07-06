<template>
  <n-select :loading="loading" v-model:value="value" :options="options" :placeholder="placeholder"
    :render-option="renderOption" :clearable="clearable ?? false" :style="{ width: width }" />
</template>

<script setup lang="ts">
import { FindRivalInfoList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { NTooltip, SelectOption } from 'naive-ui';
import { h, onMounted, ref, Ref, VNode } from 'vue';

const loading = ref(false);
const value = defineModel<number | null>("value");
const props = defineProps<{
  clearable?: boolean,
  placeholder?: string,
  width: string,
  defaultSelect?: boolean
}>();
const options: Ref<SelectOption[]> = ref([]);

function renderOption({ node, option }: { node: VNode, option: SelectOption }) {
  return h(NTooltip, {
    style: `max-width: ${props.width}; font-color: white`,
  }, {
    trigger: () => node,
    default: () => option.label
  });
};

function loadData() {
  loading.value = true;
  FindRivalInfoList()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      options.value = result.Rows.map((rival: dto.RivalInfoDto) => {
        return {
          label: rival.Name,
          value: rival.ID,
        }
      });
      if (props.defaultSelect) {
        value.value = options.value[0].value as number;
      }
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false)
}

onMounted(() => {
  loadData();
});

</script>
