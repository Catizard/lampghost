<script setup>
import HelloWorld from './components/HelloWorld.vue'
import wailsLogo from './/assets/images/logo-universal.png'
import naiveLogo from './assets/images/naive-logo.svg'
</script>

<template>
  <div class="content">
    <HelloWorld />
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
