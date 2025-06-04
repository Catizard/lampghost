<template>
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
      style="width: 200px;" :placeholder="t('placeHolderRivalTag')" :render-option="renderRivalTagOption" />
  </n-flex>
  <n-data-table :columns="columns" :data="data" :pagination="pagination" :loading="levelTableLoading"
    :row-key="(row: dto.DiffTableHeaderDto) => row.Level" :row-class-name="rowClassName" />
</template>

<script setup lang="ts">
import { FindDiffTableHeaderList, FindDiffTableHeaderTreeWithRival, FindRivalInfoList, FindRivalTagList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NDataTable, NTooltip, SelectOption } from 'naive-ui';
import { h, Ref, ref, VNode, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import DifficultTableDetail from './DifficultTableDetail.vue';
import { ClearType, ClearTypeDef, DefaultClearTypeColorStyle } from '@/constants/cleartype';

const i18n = useI18n();
const { t } = i18n;

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
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loadingDifftableOptions.value = false);
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
          ghostRivalId: currentRivalID.value,
          ghostRivalTagId: currentRivalTagID.value
        }
      )
    }
  },
  { title: "Level", key: "Name" },
  {
    title: "Failed", key: "FailCount",
    render: (row: dto.DiffTableHeaderDto) => {
      return (row.LampCount[ClearType.Failed] ?? 0)
        + (row.LampCount[ClearType.AssistEasy] ?? 0)
        + (row.LampCount[ClearType.LightAssistEasy] ?? 0);
    }
  },
  {
    title: "Easy Clear", key: "ECCount",
    render: (row: dto.DiffTableHeaderDto) => {
      return row.LampCount[ClearType.Easy] ?? 0
    }
  },
  {
    title: "Normal Clear", key: "NCCount",
    render: (row: dto.DiffTableHeaderDto) => {
      return row.LampCount[ClearType.Normal] ?? 0;
    }
  },
  {
    title: "Hard Clear", key: "HCCount",
    render: (row: dto.DiffTableHeaderDto) => {
      return row.LampCount[ClearType.Hard] ?? 0;
    }
  },
  {
    title: "EX Hard Clear", key: "EXHCCount",
    render: (row: dto.DiffTableHeaderDto) => {
      return row.LampCount[ClearType.ExHard] ?? 0;
    }
  },
  {
    title: "Full Combo+", key: "FCPlusCount",
    render: (row: dto.DiffTableHeaderDto) => {
      return (row.LampCount[ClearType.FullCombo] ?? 0)
        + (row.LampCount[ClearType.Perfect] ?? 0)
        + (row.LampCount[ClearType.Max] ?? 0);
    }
  },
  { title: "Chart Count", key: "SongCount" },
];

// NOTE: we don't build failed lamp marker, that means nothing for user
function rowClassName(row: dto.DiffTableHeaderDto): string {
  let sum = 0;
  for (const [k, v] of Object.entries(ClearType).reverse()) {
    sum += row.LampCount[k] ?? 0;
    if (sum == row.SongCount && v != ClearType.Failed) {
      // I don't have a better idea for managing this  
      const def: ClearTypeDef = DefaultClearTypeColorStyle[k];
      console.log('row: ', row.Level, 'picking', def.text);
      return def.text;
    }
  }
  return ""
}
const data: Ref<Array<dto.DiffTableHeaderDto>> = ref([]);
const pagination = false as const;
const levelTableLoading = ref<boolean>(false);
// Level => [DiffTableDataDto]
const levelData = new Map<string, Array<dto.DiffTableDataDto>>();
function loadLevelTableData(difftableID: string | number) {
  levelTableLoading.value = true;
  // TODO: remove magic one
  FindDiffTableHeaderTreeWithRival({ ID: difftableID as number, RivalID: 1 } as any)
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
    })
    .catch(err => window.$notifyError(err))
    .finally(() => levelTableLoading.value = false);
}

const loadingRivalData = ref(false);
const currentRivalID: Ref<number | null> = ref(null);
const currentRivalTagID: Ref<number | null> = ref(null);
const rivalOptions: Ref<Array<SelectOption>> = ref([]);
const rivalTagOptions: Ref<Array<SelectOption>> = ref([]);
function renderRivalTagOption({ node, option }: { node: VNode, option: SelectOption }) {
  return h(NTooltip, {
    style: "max-width: 200px; font-color: white",
  }, {
    trigger: () => node,
    default: () => option.label
  });
}
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
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loadingRivalData.value = false);
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
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loadingRivalData.value = false);
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
</script>

<i18n lang="json">{
  "en": {
    "title": "Table Statistics",
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
    "title": "难度表统计信息",
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

<style lang="css" scoped>
/* background color when clearing whole difficult level */
/* NOTE: This color style is actually different with ClearTag, so no code reuse here */
:deep(.EASY > td) {
  background-color: rgba(200, 247, 212, 0.7) !important;
}

:deep(.NORMAL > td) {
  background-color: rgba(202, 235, 253, 0.7) !important;
}

:deep(.HARD > td) {
  background-color: rgba(255, 210, 213, 0.7) !important;
}

:deep(.EX-HARD > td) {
  background-color: rgba(255, 230, 212, 0.7) !important;
}

:deep(.FULL_COMBO> td) {
  background-color: rgba(255, 241, 202, 0.7) !important;
}

:deep(.PERFECT > td) {
  background-color: rgba(255, 241, 202, 0.7) !important;
}

:deep(.MAX > td) {
  background-color: rgba(255, 241, 202, 0.7) !important;
}
</style>