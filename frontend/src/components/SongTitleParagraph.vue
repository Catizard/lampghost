<template>
  <n-flex vertical>
    <n-text strong>{{ data.Title }}</n-text>
    <n-text v-if="data.SubTitle != ''" depth="2" style="font-size: 11;">{{ data.SubTitle }}</n-text>
    <n-text v-if="data.Artist != ''" depth="3" style="font-size: 9;">{{ data.Artist }}</n-text>
  </n-flex>
</template>

<script setup lang="ts">
import { dto } from '@wailsjs/go/models';

const { data } = defineProps<{
  data: dto.RivalScoreLogDto | dto.RivalScoreDataLogDto | dto.DiffTableDataDto | dto.RivalSongDataDto
}>();
</script>
