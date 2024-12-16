<template>
    <perfect-scrollbar>
        <n-space justify="space-between">
            <n-h1 prefix="bar" style="text-align: start;">
                <n-text type="primary">
                    是的，这是难度表信息!!
                </n-text>
            </n-h1>
            <n-button @click="showAddModal = true">
                新增一个难度表
            </n-button>
        </n-space>
        <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" />
    </perfect-scrollbar>

    <n-modal v-model:show="showAddModal" preset="dialog" title="新增难度表信息" positive-text="新增" negative-text="取消"
        @positive-click="handlePositiveClick" @negative-click="handleNegativeClick" :mask-closable="false">
        <n-form ref="formRef" :model="formData" :rules="rules">
            <n-form-item label="地址" path="url">
                <n-input v-model:value="formData.url" placeholder="输入地址" />
            </n-form-item>
        </n-form>
    </n-modal>
</template>

<script lang="ts">
import type { DataTableColumns, FormInst } from 'naive-ui'
import { NButton, useMessage } from 'naive-ui'
import { useNotification } from 'naive-ui'
import { defineComponent, h, Ref, ref } from 'vue'
import { AddDiffTableHeader, FindDiffTableHeader } from '../../wailsjs/go/controller/DiffTableController'
import { entity } from '../../wailsjs/go/models'

const notification = useNotification();
const showAddModal = ref(false);

const formRef = ref<FormInst | null>(null);
const formData = ref({
    url: ""
});
const rules = {
    url: {
        required: true,
        message: "请输入地址",
        trigger: ["input", "blur"]
    }
}

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
    console.log('trying add ', url)
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

function handlePositiveClick(): boolean {
    formRef.value?.validate().then(() => {
        addDiffTableHeader(formData.value.url)
        showAddModal.value = false
    }).catch((err) => {
    })
    return false
}

function handleNegativeClick() {
    formData.value.url = ""
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
            showAddModal,
            formRef,
            formData,
            rules,
            handlePositiveClick,
            handleNegativeClick
        }
    }
})
</script>