<template>
  <n-flex vertical :size="2">
    <n-flex>
      <n-icon v-if="lost" :component="WarningOutline" color="red" />
      <n-text strong>{{ data.Title }}</n-text>
    </n-flex>
    <n-text v-if="data.SubTitle != ''" depth="2" style="font-size: 11;">{{ data.SubTitle }}</n-text>
    <n-text v-if="data.Artist != ''" depth="3" style="font-size: 9;">{{ data.Artist }}</n-text>
  </n-flex>
</template>

<script setup lang="ts">
import { WarningOutline } from '@vicons/ionicons5';
import { dto } from '@wailsjs/go/models';

const { data, lost } = defineProps<{
  data: dto.RivalScoreLogDto | dto.RivalScoreDataLogDto | dto.DiffTableDataDto | dto.RivalSongDataDto
  lost?: boolean
}>();
</script>
