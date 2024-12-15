<template>
    <perfect-scrollbar>
        <n-space justify="space-between">
            <n-h1 prefix="bar" style="text-align: start;">
                <n-text type="primary">
                    是的，这是难度表信息!!
                </n-text>
            </n-h1>
            <n-button>
                新增一个难度表
            </n-button>
        </n-space>
        <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" />
    </perfect-scrollbar>
</template>

<script lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, useMessage } from 'naive-ui'
import { useNotification } from 'naive-ui'
import { defineComponent, h, Ref, ref } from 'vue'
import { AddDiffTableHeader, FindDiffTableHeader } from '../../wailsjs/go/controller/DiffTableController'
import { entity } from '../../wailsjs/go/models'

const notification = useNotification();

function createColumns({ play }: { play: (row: entity.DiffTableHeader) => void }): DataTableColumns<entity.DiffTableHeader> {
    return [
        {
            title: "Name",
            key: "name"
        },
        {
            title: "url",
            key: "HeaderUrl"
        },
        {
            title: 'Action',
            key: 'actions',
            render(row) {
                return h(
                    NButton,
                    {
                        strong: true,
                        tertiary: true,
                        size: 'small',
                        onClick: () => play(row)
                    },
                    { default: () => 'Play' }
                )
            }
        }
    ]
}

let data: Ref<Array<entity.DiffTableHeader>> = ref([]);

function addDiffTableHeader(url: string) {
    AddDiffTableHeader(url).then(result => {
        notification['success']({
            content: "新增难度表似乎成功了，后台返回结果: " + result,
            meta: "成功?",
            duration: 2500,
            keepAliveOnHover: true
        });
    }).catch(err => {
        notification['error']({
            content: "新增难度表似乎失败了，后台返回错误: " + err,
            meta: "失败?",
            duration: 2500,
            keepAliveOnHover: true
        });
    })
}

function loadDiffTableData() {
    FindDiffTableHeader().then(result => {
        data.value = [...result]
        console.log(result)
    }).catch(err => {
        console.error(err)
    })
}

export default defineComponent({
    setup() {
        const message = useMessage()
        loadDiffTableData();
        return {
            data,
            columns: createColumns({
                play(row: entity.DiffTableHeader) {
                    message.info(`Play ${row.name}`)
                }
            }),
            pagination: false as const,
            loadDiffTableData,
            notification,
        }
    }
})
</script>