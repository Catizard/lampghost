<template>
  <n-tag :color="{ color: def.color, textColor: def.textColor }">
    {{ def.text }}
  </n-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  state: boolean
}>();

interface TagDef {
  value: boolean;
  text: string;
  color: string;
  textColor: string;
}

const table: Array<TagDef> = [
  {
    value: false,
    text: "NO",
    color: "#FF6B74",
    textColor: "#ffffff",
  },
  {
    value: true,
    text: "YES",
    color: "#17964E",
    textColor: "#ffffff",
  }
];

function getDef(state: boolean): TagDef {
  return state ? table[1] : table[0];
}


const def = computed(() => getDef(props.state));

</script>