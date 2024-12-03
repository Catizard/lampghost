<script setup>
import {reactive} from 'vue'
import {Greet} from '../../wailsjs/go/main/App'
import { useMessage } from 'naive-ui'

const message = useMessage()

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
})

function greet() {
  Greet(data.name).then(result => {
    message.info(result)
  })
}


</script>

<template>
  <main>
    <div class="result">{{ data.resultText }}</div>
    <div class="input-box">
      <n-input v-model:value="data.name" type="text" placeholder="input name" class="input"/>
      <n-button type="info" @click="greet"> Greet </n-button>
    </div>
  </main>
</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}
.input-box {
  text-align: center;
}
.input {
  width: 200px;
  margin-right: 10px;
}
</style>
