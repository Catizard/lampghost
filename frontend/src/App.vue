<script setup>
import { RouterView } from 'vue-router';
</script>

<template>
  <div class="content">
    <RouterView />
  </div>
</template>

<style>
.content {
  height: 100%;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
