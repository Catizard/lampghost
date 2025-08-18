<template>
  <n-flex vertical :size="4">
    <n-flex>
      <n-icon v-if="lost" :component="WarningOutline" color="red" />
      <strong style="font-size: 1.1em;">{{ data.Title }}</strong>
    </n-flex>
    <n-text v-if="data.SubTitle != ''" depth="2" style="font-size: 0.9em;">{{ data.SubTitle }}</n-text>
    <n-text v-if="data.Artist != ''" depth="3" style="font-size: 0.9em;">{{ data.Artist }}</n-text>
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
