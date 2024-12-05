<template>
  <n-grid :cols="12">
    <n-gi :span="8">
      <div class="light-green" />
    </n-gi>
    <n-gi :span="4">
      <div class="green" />
    </n-gi>
  </n-grid>
</template>

<script lang="ts" setup>
</script>

<style scoped>
.light-green {
  height: 108px;
  background-color: rgba(0, 128, 0, 0.12);
}

.green {
  height: 108px;
  background-color: rgba(0, 128, 0, 0.24);
}
</style>