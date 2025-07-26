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
import { Close } from '@vicons/ionicons5';

let directories = defineModel<string[]>("directories");

function deletePath(rev: string) {
  directories.value = [...directories.value.filter(directory => directory != rev)];
}
</script>
