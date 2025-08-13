<template>
  <template v-if="noplay">/</template>
  <template v-else>
    <n-flex vertical align="center">
      <n-text strong>{{ rank }}</n-text>
      <n-text depth="2" style="font-size: 11;">{{ accuracy.toFixed(2) }}%</n-text>
      <n-text depth="3" style="font-size: 9;">[{{ exscore }}]</n-text>
    </n-flex>
  </template>

</template>

<script setup lang="ts">
import { dto } from '@wailsjs/go/models';
import { computed } from 'vue';

const { data } = defineProps<{
  data: dto.DiffTableDataDto | dto.RivalScoreDataLogDto
}>();

const noplay = computed<boolean>(() => {
  return data.Notes == 0;
})

const exscore = computed<number>(() => {
  return (data.Epg + data.Lpg) * 2 + data.Egr + data.Lgr;
});

const accuracy = computed<number>(() => {
  return exscore.value * 50 / data.Notes;
});

const rank = computed<string>(() => {
  const acc = accuracy.value;
  if (acc >= 88.88) {
    return "AAA";
  } else if (acc >= 77.77) {
    return "AA";
  } else if (acc >= 66.66) {
    return "A";
  } else if (acc >= 55.55) {
    return "B";
  } else if (acc >= 44.44) {
    return "C";
  } else if (acc >= 33.33) {
    return "D";
  } else if (acc >= 22.22) {
    return "E";
  }
  return "F";
});

</script>
