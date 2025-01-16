<template>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" />
</template>

<script setup lang="ts">
import { DataTableColumns, useNotification } from 'naive-ui';
import { dto } from '../../../wailsjs/go/models';
import { ref, Ref, watch } from 'vue';
import { QueryDiffTableDataWithRival } from '../../../wailsjs/go/controller/DiffTableController';

const notification = useNotification();

const props = defineProps<{
    headerId?: number
    level?: string
}>()

function createColumns(): DataTableColumns<dto.DiffTableDataDto> {
    return [
        { title: "Song Name", key: "Title" },
        { title: "Artist", key: "Artist" },
        { title: "Play Count", key: "PlayCount" },
        { title: "Clear", key: "Lamp" },
    ];
}
const columns = createColumns();

let data: Ref<Array<any>> = ref([]);
function loadData(headerId: number, level: string) {
    // TODO: remove magic 1
    QueryDiffTableDataWithRival(headerId, level, 1)
    .then(result => {
        if (result.Code != 200) {
            return Promise.reject(result.Msg)
        }
        data.value = [...result.Rows].map(row => {
            return {
                key: row.ID,
                ...row
            }
        })
    }).catch(err => {
        notification.error({
            content: err,
            duration: 3000,
            keepAliveOnHover: true
        })
    })
}

const pagination = false as const;

watch(props, (newProps) => {
    loadData(newProps.headerId, newProps.level)
}, { immediate: true })

</script>