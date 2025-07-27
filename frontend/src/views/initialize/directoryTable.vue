<template>
  <n-table>
    <tbody>
      <tr v-for="path in directories">
        <n-button type="error" size="small" @click="deletePath(path)">
          <template #icon>
            <Close />
          </template>
        </n-button>
        {{ path }}
      </tr>
    </tbody>
  </n-table>
</template>

<script setup lang="ts">
import { DataTableColumns } from 'naive-ui';
import { Close } from '@vicons/ionicons5';
import { watch } from 'vue';

let directories = defineModel<string[]>("directories");

function deletePath(rev: string) {
  console.log('rev: ', rev);
  directories.value = [...directories.value.filter(directory => directory != rev)];
}

watch(directories, (newValue) => {
  console.log(directories.value);
})
</script>
