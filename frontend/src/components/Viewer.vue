<template>
    <RouterView />
</template>

<script lang="ts" setup>
import { useLoadingBar } from 'naive-ui';
import { RouterView } from 'vue-router';
import router from '../router';

const loadingBar = useLoadingBar();

router.beforeEach(() => loadingBar?.start());
router.afterEach(() => loadingBar?.finish());
</script>
