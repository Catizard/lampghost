<template>
  <n-flex justify="space-between">
    <n-h1 prefix="bar" style="text-align: start;">
      <n-text type="primary">{{ t('title.tableStatistics') }}</n-text>
    </n-h1>
  </n-flex>
  <n-flex justify="start" style="margin-bottom: 8px;">
    <SelectDifficultTable v-model:value="currentDiffTableID" style="width: 200px;" />
    <SelectRival v-model:value="currentRivalID" width="200px" :placeholder="t('form.placeholderRival')" clearable />
    <SelectRivalTag v-model:value="currentRivalTagID" :rivalId="currentRivalID" width="200px"
      :placeholder="t('form.placeholderRivalTag')" clearable />
  </n-flex>
  <n-data-table :columns="columns" :data="data" :pagination="pagination" :loading="levelTableLoading"
    :row-key="(row: dto.DiffTableHeaderDto) => row.Level" :row-class-name="rowClassName" />
</template>

<script setup lang="ts">
import { FindDiffTableHeaderTreeWithRival } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import { DataTableColumns, NDataTable } from 'naive-ui';
import { h, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import DifficultTableDetail from './DifficultTableDetail.vue';
import { ClearType, ClearTypeDef, DefaultClearTypeColorStyle } from '@/constants/cleartype';
import SelectRivalTag from '@/components/rivals/SelectRivalTag.vue';
import SelectDifficultTable from '@/components/difficult_table/SelectDifficultTable.vue';
import SelectRival from '@/components/rivals/SelectRival.vue';
import { useUserStore } from '@/stores/user';

const i18n = useI18n();
const { t } = i18n;
const userStore = useUserStore();

const currentDiffTableID: Ref<number | null> = ref(null);

const columns: DataTableColumns<dto.DiffTableHeaderDto> = [
  {
    type: "expand", align: "center",
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
  { title: "Level", key: "Name", align: "center" },
  {
    title: "Failed", key: "FailCount", align: "center",
    render(row: dto.DiffTableHeaderDto): number | string {
      return positiveOrEmpty(
        (row.LampCount[ClearType.Failed] ?? 0)
        + (row.LampCount[ClearType.AssistEasy] ?? 0)
        + (row.LampCount[ClearType.LightAssistEasy] ?? 0)
      );
    }
  },
  {
    title: "Easy Clear", key: "ECCount", align: "center",
    render(row: dto.DiffTableHeaderDto): number | string {
      return positiveOrEmpty(row.LampCount[ClearType.Easy] ?? 0);
    }
  },
  {
    title: "Normal Clear", key: "NCCount", align: "center",
    render(row: dto.DiffTableHeaderDto): number | string {
      return positiveOrEmpty(row.LampCount[ClearType.Normal] ?? 0);
    }
  },
  {
    title: "Hard Clear", key: "HCCount", align: "center",
    render(row: dto.DiffTableHeaderDto): number | string {
      return positiveOrEmpty(row.LampCount[ClearType.Hard] ?? 0);
    }
  },
  {
    title: "EX Hard Clear", key: "EXHCCount", align: "center",
    render(row: dto.DiffTableHeaderDto): number | string {
      return positiveOrEmpty(row.LampCount[ClearType.ExHard] ?? 0);
    }
  },
  {
    title: "Full Combo+", key: "FCPlusCount", align: "center",
    render(row: dto.DiffTableHeaderDto): number | string {
      return positiveOrEmpty(
        (row.LampCount[ClearType.FullCombo] ?? 0)
        + (row.LampCount[ClearType.Perfect] ?? 0)
        + (row.LampCount[ClearType.Max] ?? 0)
      );
    }
  },
  { title: "Chart Count", key: "SongCount", align: "center" },
];

// NOTE: we don't build failed lamp marker, that means nothing for user
function rowClassName(row: dto.DiffTableHeaderDto): string {
  let sum = 0;
  for (const [k, v] of Object.entries(ClearType).reverse()) {
    sum += row.LampCount[k] ?? 0;
    if (sum == row.SongCount && v != ClearType.Failed) {
      // I don't have a better idea for managing this
      const def: ClearTypeDef = DefaultClearTypeColorStyle[k];
      return `row-${def.text}`;
    }
  }
  return "";
}
const data: Ref<Array<dto.DiffTableHeaderDto>> = ref([]);
const pagination = false as const;
const levelTableLoading = ref<boolean>(false);
// Level => [DiffTableDataDto]
const levelData = new Map<string, Array<dto.DiffTableDataDto>>();
function loadLevelTableData(difftableID: string | number) {
  levelTableLoading.value = true;
  FindDiffTableHeaderTreeWithRival({ ID: difftableID as number, RivalID: userStore.id } as any)
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

function positiveOrEmpty(v: number): number | string {
  return v == 0 ? "" : v;
}

const currentRivalID: Ref<number | null> = ref(null);
const currentRivalTagID: Ref<number | null> = ref(null);

// Watch 1: Whenever changing current difftable, reload the level table
watch(currentDiffTableID, (newID: string | number) => {
  loadLevelTableData(newID);
});
</script>

<style lang="css" scoped>
/* background color when clearing whole difficult level */
/* NOTE: This color style is actually different with ClearTag, so no code reuse here */
:deep(.row-EASY > td) {
  background-color: rgba(200, 247, 212, 0.7) !important;
}

:deep(.row-NORMAL > td) {
  background-color: rgba(202, 235, 253, 0.7) !important;
}

:deep(.row-HARD > td) {
  background-color: rgba(255, 210, 213, 0.7) !important;
}

:deep(.row-EX-HARD > td) {
  background-color: rgba(255, 230, 212, 0.7) !important;
}

:deep(.row-FULL_COMBO > td) {
  background-color: rgba(255, 241, 202, 0.7) !important;
}

:deep(.row-PERFECT > td) {
  background-color: rgba(255, 241, 202, 0.7) !important;
}

:deep(.row-MAX > td) {
  background-color: rgba(255, 241, 202, 0.7) !important;
}

/*
* header row
* NOTE: This also overrides DifficultTableDetail.vue's header
*/
:deep(th.n-data-table-th) {
  font-weight: bold;
  font-size: 1.1em;
}
</style>
