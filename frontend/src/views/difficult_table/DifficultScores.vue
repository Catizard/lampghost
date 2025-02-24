<template>
  <perfect-scrollbar>
    <n-flex justify="space-between">
      <n-h1 prefix="bar" style="text-align: start;">
        <n-text type="primary">难度表得分情况</n-text>
      </n-h1>
    </n-flex>
    <n-flex justify="flex-start">
      <n-select :loading="levelTableLoading" v-model:value="currentDiffTableID" :options="difftableOptions"
        style="width: 200px;" />
      <n-select :loading="loadingRivalData" v-model:value="currentRivalID" :options="rivalOptions" style="width: 200px;"
        placeholder="Ghost Rival" />
      <n-select :loading="loadingRivalData" v-model:value="currentRivalTagID" :options="rivalTagOptions"
        style="width: 200px;" placeholder="Rival Tag" />
    </n-flex>
    <n-data-table :columns="columns" :data="data" :pagination="pagination" :loading="levelTableLoading"
      :row-key="(row: dto.DiffTableHeaderDto) => row.Level"
      @update-expanded-row-keys="handleUpdateLevelTableExpandedRowKeys" />
  </perfect-scrollbar>
</template>

<script setup lang="ts">
import ClearTag from '@/components/ClearTag.vue';
import { FindDiffTableHeaderList, FindDiffTableHeaderTree, QueryDiffTableDataWithRival } from '@wailsjs/go/controller/DiffTableController';
import { FindRivalInfoList } from '@wailsjs/go/controller/RivalInfoController';
import { FindRivalTagList } from '@wailsjs/go/controller/RivalTagController';
import { dto, vo } from '@wailsjs/go/models';
import { DataTableColumns, NButton, NDataTable, SelectOption, useNotification } from 'naive-ui';
import { h, Ref, ref, watch } from 'vue';

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
        return Promise.reject("目前无法处理一个难度表都没有的情况, 请至少添加一个")
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
        NDataTable,
        {
          columns: songDataColumns,
          data: levelData[row.Level],
          pagination: false,
          bordered: false,
          rowKey: (row: dto.DiffTableDataDto) => row.ID,
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
        return Promise.reject("该难度表不存在，发生什么事了?");
      }
      if (result.Rows.length != 1) {
        return Promise.reject("这个难度表不只一个数据，发生什么事了?");
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
        return Promise.reject("一个rival都没有, 发生什么事了?");
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

function handleUpdateLevelTableExpandedRowKeys(keys: Array<string>) {
  // TODO: do we need to handle multiple loading state?
  keys.forEach((level) => {
    if (!levelData.has(level)) {
      // TODO: remove magic 1
      loadSongTableData(currentDiffTableID.value, level, 1)
    }
  });
}

const songDataColumns: DataTableColumns<dto.DiffTableDataDto> = [
  { title: "Song Name", key: "Title", ellipsis: true, resizable: true },
  { title: "Artist", key: "Artist", resizable: true },
  { title: "Play Count", key: "PlayCount", minWidth: "100px", resizable: true },
  {
    title: "Clear", key: "Lamp", minWidth: "100px", resizable: true,
    render(row) {
      return h(
        ClearTag,
        {
          clear: row.Lamp
        },
      )
    }
  },
  {
    title: "Action",
    key: "actions",
    resizable: true,
    minWidth: "150px",
    render(row) {
      return h(
        NButton,
        {
          strong: true,
          tertiary: true,
          size: "small",
          onClick: () => handleAddToFolder(row.ID),
        },
        { default: () => "添加至收藏夹" }
      )
    }
  }
];

// TODO: pass ghost rival's ID & tag here
function loadSongTableData(headerID: number, level: string, rivalID: number) {
  levelTableLoading.value = true;
  QueryDiffTableDataWithRival(headerID, level, rivalID)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      levelData[level] = [...result.Rows];
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

// TODO: implement me!
function handleAddToFolder(songDataID: number) {
}

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