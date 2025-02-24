<template>
  <n-tag :color="{ color: def.color, textColor: def.textColor}">
    {{ def.text }}
  </n-tag>
</template>

<script setup lang="ts">
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
    color: "#FF0000",
    textColor: "#000000",
  },
  {
    value: true,
    text: "YES",
    color: "#00FF00",
    textColor: "#000000",
  }
];

function getDef(state: boolean): TagDef {
  return state ? table[1] : table[0];
}


const def = getDef(props.state);

</script>