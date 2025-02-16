<template>
  <div class="container">
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: left;">
        <n-text type="primary">
          Settings
        </n-text>
      </n-h1>
      <n-button type="primary" @click="handleSaveSettings">保存</n-button>
    </n-flex>

    <n-form ref="formRef" :model="model">
      <n-h2>
        存档设置
        <n-p>
          <n-form-item label="用户名, 如果你想换个用户名的话" path="userName">
            <n-input show-count v-model:value="model.userName" placeholder="请输入用户名" style="width: 50%;" />
          </n-form-item>
        </n-p>
        <n-p class="alert-p">
          注意: 如果你修改了文件路径地址, 保存设置时会重新再次加载你的存档数据, 因此你不需要再手动同步一次存档
        </n-p>
        <n-p>
          <n-form-item label="scorelog.db文件路径" path="scorelogFilePath">
            <n-input clearable v-model:value="model.scorelogFilePath" placeholder="请输入scorelog.db文件路径"
              style="width: 50%;" />
          </n-form-item>
        </n-p>
        <n-p>
          <n-form-item label="songdata.db文件路径" path="songdataFilePath">
            <n-input clearable v-model:value="model.songdataFilePath" placeholder="请输入songdata.db文件路径"
              style="width: 50%;" />
          </n-form-item>
        </n-p>
      </n-h2>
      <n-h2>
        <n-text>
          自定义收藏夹
        </n-text>
        <n-p class="alert-p">
          注意: 如果你修改了下列设置之前已经导入过难度表, 需要修改表定义之后重新刷新所有难度表信息来移出此前导入的表信息！
        </n-p>
        <n-form-item label="内部服务器端口号" path="internalServerPort">
          <n-input-number :show-button="false" v-model:value="model.internalServerPort" placeholder="请输入端口号"
            :maxlength="5" style="width: 150px;" />
        </n-form-item>
        <n-form-item label="难度表标志" path="folderSymbol">
          <n-input v-model:value="model.folderSymbol" placeholder="默认为空" :maxlength="5" style="width: 150px;">
            <template #prefix>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-icon :component="HintIcon" />
                </template>
                难度表标志即发狂表的🌟或satellite的sl, 建议留空。
              </n-tooltip>
            </template>
          </n-input>
        </n-form-item>
      </n-h2>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { FormInst } from 'naive-ui';
import { ref } from 'vue';
import {
  ChatboxEllipsesOutline as HintIcon,
} from '@vicons/ionicons5';

const formRef = ref<FormInst | null>(null);
const model = ref({
  internalServerPort: null,
  userName: null,
  scorelogFilePath: null,
  songdataFilePath: null,
  folderSymbol: null,
})

function handleSaveSettings() {
  console.log('TODO')
}
</script>

<style scoped>
.container {
  text-align: left;
}

.alert-p {
  color: #ff1100;
}
</style>