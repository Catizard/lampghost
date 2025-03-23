<template>
  <n-tag :color="{ color: def.color, textColor: def.textColor }">
    {{ def.text }}
  </n-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  clear: number
}>();

interface ClearTypeDef {
  value: number;
  text: string;
  color: string;
  textColor: string;
}

const table: Array<ClearTypeDef> = [
  {
    value: 0,
    text: "NO_PLAY",
    color: "#8C8C8C",
    textColor: "#ffffff",
  },
  {
    value: 1,
    text: "FAIELD",
    color: "#CC5C76",
    textColor: "#ffffff",
  },
  {
    value: 2,
    text: "A-Easy",
    color: "#FF9FF9",
    textColor: "#ffffff",
  },
  {
    value: 3,
    text: "LA-Easy",
    color: "#FF9FF9",
    textColor: "#ffffff",
  },
  {
    value: 4,
    text: "EASY",
    color: "#49E670",
    textColor: "#ffffff",
  },
  {
    value: 5,
    text: "NORMAL",
    color: "#4FBCF7",
    textColor: "#ffffff",
  },
  {
    value: 6,
    text: "HARD",
    color: "#FF6B74",
    textColor: "#ffffff",
  },
  {
    value: 7,
    text: "EX-HARD",
    color: "#FFAD70",
    textColor: "#ffffff",
  },
  {
    value: 8,
    text: "FULL-COMBO",
    color: "#FFD251",
    textColor: "#ffffff",
  },
  {
    value: 9,
    text: "MAX",
    color: "#FFD251",
    textColor: "#ffffff",
  }
]

function getDef(clearType: number): ClearTypeDef {
  return table[clearType];
}

const def = computed(() => getDef(props.clear));

</script>