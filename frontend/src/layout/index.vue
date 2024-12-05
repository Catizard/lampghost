<template>
    <n-layout position="absolute" has-sider>
        <Sidebar />
        <Content />
    </n-layout>
</template>

<script lang="ts" setup>
import Content from './content.vue';
import Sidebar from './sidebar.vue'
</script>
