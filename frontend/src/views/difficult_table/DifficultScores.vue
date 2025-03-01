<template>
  <perfect-scrollbar>
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">{{ t('title') }}</n-text>
      </n-h1>
    </n-flex>
    <n-flex justify="flex-start">
      <n-select :loading="levelTableLoading" v-model:value="currentDiffTableID" :options="difftableOptions"
        style="width: 200px;" />
      <n-select :loading="loadingRivalData" v-model:value="currentRivalID" :options="rivalOptions" style="width: 200px;"
        :placeholder="t('placeHolderRival')" />
      <n-select :loading="loadingRivalData" v-model:value="currentRivalTagID" :options="rivalTagOptions"
        style="width: 200px;" :placeholder="t('placeHolderRivalTag')" />
    </n-flex>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :loading="levelTableLoading"
      :row-key="(row: dto.DiffTableHeaderDto) => row.Level" />
  </perfect-scrollbar>
</template>

<script setup lang="ts">
import ClearTag from '@/components/ClearTag.vue';
import { FindDiffTableHeaderList, FindDiffTableHeaderTree, QueryDiffTableDataWithRival } from '@wailsjs/go/controller/DiffTableController';
import { FindRivalInfoList } from '@wailsjs/go/controller/RivalInfoController';
import { FindRivalTagList } from '@wailsjs/go/controller/RivalTagController';
import { dto, vo } from '@wailsjs/go/models';
import { DataTableColumns, NButton, NDataTable, SelectOption, useNotification } from 'naive-ui';
import { h, nextTick, Ref, ref, useTemplateRef, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import DifficultTableDetail from './DifficultTableDetail.vue';

const i18n = useI18n();
const { t } = i18n;
const notification = useNotification();

const loadingDifftableOptions = ref<boolean>(false);
const currentDiffTableID: Ref<number | null> = ref(null);
const difftableOptions: Ref<Array<SelectOption>> = ref([]);
function loadDifftableOptions() {
  loadingDifftableOptions.value = true;
  FindDiffTableHeaderList()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg)
      }
      if (result.Rows.length == 0) {
        return Promise.reject(t('message.noTableerror'))
      }
      difftableOptions.value = result.Rows.map((row: dto.DiffTableHeaderDto) => {
        return {
          label: row.Name,
          value: row.ID,
        } as SelectOption
      });
      currentDiffTableID.value = difftableOptions.value[0].value as number;
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      });
    }).finally(() => {
      loadingDifftableOptions.value = false;
    });
}
loadDifftableOptions();

const columns: DataTableColumns<dto.DiffTableHeaderDto> = [
  {
    type: "expand",
    renderExpand: (row: dto.DiffTableHeaderDto) => {
      return h(
        DifficultTableDetail,
        {
          headerId: row.ID,
          level: row.Level,
        }
      )
    }
  },
  { title: "Level Name", key: "Name" },
  { title: "HC+ Count", key: "HCPlusCount" },
  { title: "Song Count", key: "SongCount" },
];
const data: Ref<Array<dto.DiffTableHeaderDto>> = ref([]);
const pagination = false as const;
const levelTableLoading = ref<boolean>(false);
// Level => [DiffTableDataDto]
const levelData = new Map<string, Array<dto.DiffTableDataDto>>();
function loadLevelTableData(difftableID: string | number) {
  levelTableLoading.value = true;
  FindDiffTableHeaderTree({ ID: difftableID as number } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Rows.length == 0) {
        return Promise.reject(t('message.tableNotFoundError'));
      }
      if (result.Rows.length != 1) {
        return Promise.reject(t('message.duplicateTableError'));
      }
      data.value = [...result.Rows[0].Children];
      levelData.clear();
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      });
    }).finally(() => {
      levelTableLoading.value = false;
    });
}

const loadingRivalData = ref(false);
const currentRivalID: Ref<number | null> = ref(null);
const currentRivalTagID: Ref<number | null> = ref(null);
const rivalOptions: Ref<Array<SelectOption>> = ref([]);
const rivalTagOptions: Ref<Array<SelectOption>> = ref([]);
function loadRivalOptions() {
  loadingRivalData.value = true;
  FindRivalInfoList()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      if (result.Rows.length == 0) {
        return Promise.reject(t('message.noRivalError'));
      }
      rivalOptions.value = result.Rows.map((row: dto.RivalInfoDto) => {
        return {
          label: row.Name,
          value: row.ID,
        } as SelectOption
      });
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      })
    }).finally(() => {
      loadingRivalData.value = false;
    });
}
function loadRivalTagOptions(rivalID: number) {
  // TODO: same logic, and should be handled together
  loadingRivalData.value = true;
  FindRivalTagList({ RivalId: rivalID } as any)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      rivalTagOptions.value = result.Rows.map((row: dto.RivalTagDto) => {
        return {
          label: row.TagName,
          value: row.ID,
        } as SelectOption
      });
    }).catch(err => {
      notification.error({
        content: err,
        duration: 3000,
        keepAliveOnHover: true,
      })
    }).finally(() => {
      loadingRivalData.value = false;
    });
}
loadRivalOptions();

// Watch 1: Whenever changing current difftable, reload the level table
watch(currentDiffTableID, (newID: string | number) => {
  loadLevelTableData(newID);
});

// Watch 2: Whenever changing ghost rival, reload its corresponding tags
watch(currentRivalID, (newID: number) => {
  loadRivalTagOptions(newID);
});

// TODO: Watch3: Whenever changing current rival ID or tag, reload the song table
</script>

<i18n lang="json">{
  "en": {
    "title": "Table Scores",
    "column": {
      "songName": "Song Name",
      "artist": "Artist",
      "count": "Play Count",
      "clear": "Clear",
      "ghost": "Ghost",
      "actions": "Actions"
    },
    "message": {
      "noTableError": "Cannot handle no difficult table data currenlty, please add at least one table first",
      "tableNotFoundError": "FATAL ERROR: cannot load table data?",
      "duplicateTableError": "FATAL ERROR: table data is duplicated?",
      "noRivalError": "FATAL ERROR: cannot find at least one rival?"
    },
    "button": {
      "addToFolder": "Add to Folder"
    },
    "placeHolderRival": "Choose Rival",
    "placeHolderRivalTag": "Choose Rival Tag"
  },
  "zh-CN": {
    "title": "难度表得分",
    "column": {
      "songName": "谱名",
      "artist": "作者",
      "count": "游玩次数",
      "clear": "点灯",
      "ghost": "影子灯",
      "actions": "操作"
    },
    "message": {
      "noTableError": "目前无法处理一个难度表都没有的情况，请至少先添加一个难度表",
      "tableNotFoundError": "未知错误: 无法读取难度表数据?",
      "duplicateTableError": "未知错误: 难度表重复?",
      "noRivalError": "未知错误: 无玩家信息?"
    },
    "button": {
      "addToFolder": "添加至收藏夹"
    },
    "placeHolderRival": "选择对比玩家",
    "placeHolderRivalTag": "选择对比玩家的标签"
  }
}</i18n>