<script setup>
import HelloWorld from './components/HelloWorld.vue'
import wailsLogo from './/assets/images/logo-universal.png'
import naiveLogo from './assets/images/naive-logo.svg'
</script>

<template>
  <div class="content">
    <n-message-provider >
      <n-grid x-gap="12" :cols="2">
        <n-grid-item>
          <n-image height="350" :src="wailsLogo" />
        </n-grid-item>
        <n-grid-item>
          <n-image height="350" :src="naiveLogo" />
        </n-grid-item>
      </n-grid>
      <HelloWorld/>
    </n-message-provider>
  </div>

</template>

<style>
.content {
  display: block;
  width: 90%;
  height: 60%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
