<template>
    <div style="padding: 20px 40px; min-height: 85vh;">
        <n-card>
            <n-space justify="space-between" vertical style="width: 65%;">
                <n-h1 prefix="bar" style="text-align: start">
                    <n-text type="primary">初始化用户信息</n-text>
                </n-h1>
                <n-form ref="formRef" :rules="rules" :model="formData">
                    <n-form-item label="用户名" path="userName">
                        <n-input v-model:value="formData.userName" placeholder="请输入用户名" />
                    </n-form-item>
                    <n-form-item label="songdata.db路径" path="songdataPath">
                        <n-input v-model:value="formData.songdataPath" placeholder="请输入songdata.db路径" />
                    </n-form-item>
                    <n-form-item label="scorelog.db路径" path="scorelogPath">
                        <n-input v-model:value="formData.scorelogPath" placeholder="请输入scorelog.db路径" />
                    </n-form-item>
                    <n-form-item>
                        <n-button attr-type="button" @click="handleSubmit" type="primary">提交</n-button>
                    </n-form-item>
                </n-form>
            </n-space>
        </n-card>
    </div>

</template>

<script setup>
import { useNotification } from 'naive-ui';
import { ref } from 'vue';
import { InitializeMainUser } from '../../wailsjs/go/controller/RivalInfoController';
import router from '../router';

const notification = useNotification();

const formRef = ref(null);
const formData = ref({
    songdataPath: "",
    scorelogPath: "",
    userName: ""
});

const rules = {
    userName: {
        required: true,
        message: "请输入用户名称",
        trigger: ["input", "blur"],
    },
    songdataPath: {
        required: true,
        message: "请输入songdata.db路径",
        trigger: ["input", "blur"],
    },
    scorelogPath: {
        required: true,
        message: "请输入scorelog.db路径",
        trigger: ["input", "blur"],
    }
};

function handleSubmit(e) {
    e.preventDefault();
    formRef.value?.validate(errors => {
        if (errors) {
            console.error(errors);
        } else {
            const rivalInfo = {
                Name: formData.value.userName,
                ScoreLogPath: formData.value.scorelogPath,
                SongDataPath: formData.value.songdataPath,
            };
            InitializeMainUser(rivalInfo).then(result => {
                if (result.Code != 200) {
                    return Promise.reject(result.Msg);
                }
                notification.success({
                    content: "初始化成功",
                    duration: 1500,
                });
                router.push("/home");
            }).catch(err => {
                notification.error({
                    content: err,
                    duration: 3000
                })
            })
        }
    })
}
</script>
