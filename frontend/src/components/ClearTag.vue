<template>
  <n-tag :color="{ color: def.color, textColor: def.textColor }">
    {{ def.text }}
  </n-tag>
</template>

<script setup lang="ts">
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
    color: "#ff0000",
    textColor: "#ffffff",
  },
  {
    value: 2,
    text: "A-Easy",
    color: "#7f00bb",
    textColor: "#ffffff",
  },
  {
    value: 3,
    text: "LA-Easy",
    color: "#d78df7",
    textColor: "#ffffff",
  },
  {
    value: 4,
    text: "EASY",
    color: "#00ff00",
    textColor: "#ffffff",
  },
  {
    value: 5,
    text: "NORMAL",
    color: "#0000ff",
    textColor: "#ffffff",
  },
  {
    value: 6,
    text: "HARD",
    color: "#ffbb00",
    textColor: "#ffffff",
  },
  {
    value: 7,
    text: "EX-HARD",
    color: "#ffee00",
    textColor: "#ffffff",
  },
  {
    value: 8,
    text: "FULL-COMBO",
    color: "#00ffcc",
    textColor: "#ffffff",
  },
  {
    value: 9,
    text: "MAX",
    color: "#000000",
    textColor: "#ffffff",
  }
]

function getDef(clearType: number): ClearTypeDef {
  return table[clearType];
}

const def = getDef(props.clear);

</script>