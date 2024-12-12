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
import { defineComponent, h } from 'vue'
import { AddDiffTableHeader } from '../../wailsjs/go/controller/DiffTableController'

const notification = useNotification();

interface Song {
    no: number
    title: string
    length: string
}

function createColumns({
    play
}: {
    play: (row: Song) => void
}): DataTableColumns<Song> {
    return [
        {
            title: 'No',
            key: 'no'
        },
        {
            title: 'Title',
            key: 'title'
        },
        {
            title: 'Length',
            key: 'length'
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

const data: Song[] = [
    { no: 3, title: 'Wonderwall', length: '4:18' },
    { no: 4, title: 'Don\'t Look Back in Anger', length: '4:48' },
    { no: 12, title: 'Champagne Supernova', length: '7:27' }
]

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

export default defineComponent({
    setup() {
        const message = useMessage()
        return {
            data,
            columns: createColumns({
                play(row: Song) {
                    message.info(`Play ${row.title}`)
                }
            }),
            pagination: false as const
        }
    }
})
</script>