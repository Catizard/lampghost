<template>
  <n-radio-group v-model:value="currentScheme">
    <n-radio-button v-for="option in schemeOptions" :key="option.key" :value="option.value" :label="option.label" />
  </n-radio-group>
  <n-data-table :columns="categoryColumns" :data="categoryData" :row-key="(category: Category) => category.key" />
  <n-flex justify="end" style="margin-top: 20px;">
    <n-button :loading="loading" attr-type="button" @click="handleSkip">
      {{ t('button.skip') }}
    </n-button>
    <n-button :loading="loading" attr-type="button" @click="handleSubmit" type="primary">
      {{ t('button.submit') }}
    </n-button>
  </n-flex>
</template>

<script setup lang="ts">
import { AddBatchDiffTableHeader, QueryPredefineTableSchemes } from '@wailsjs/go/main/App';
import { entity, vo } from '@wailsjs/go/models';
import { DataTableColumn, DataTableRowKey, NDataTable, NTag, useModal } from 'naive-ui';
import { TagColor } from 'naive-ui/es/tag/src/common-props';
import { h, reactive, Reactive, Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  moveOn: () => void
}>();

const { t } = useI18n();
const modal = useModal();
const loading = ref(false);

type RadioOption = {
  key: string,
  value: string,
  label: string
};

type Category = {
  key: string,
  name: string,
  headers: entity.PredefineTableHeader[]
};

const currentScheme: Ref<string | null> = ref(null);
const schemeOptions: Ref<Array<RadioOption>> = ref([]);

const categoryColumns: Array<DataTableColumn<Category>> = [
  {
    type: "expand",
    renderExpand(row: Category) {
      console.log('row.key:', row.key);
      return h(NDataTable, {
        columns: headerColumns,
        data: row.headers,
        rowKey: (header: entity.PredefineTableHeader) => header.HeaderUrl,
        "onUpdate:checkedRowKeys": (rowKeys: DataTableRowKey[]) => {
          checkedHeaderKey.set(row.key, rowKeys);
        }
      });
    }
  },
  { title: t('column.name'), key: "name" }
];
let currentHeaders: entity.PredefineTableHeader[] = [];
const categoryData: Ref<Category[]> = ref([]);
// category name => selected header urls
const checkedHeaderKey: Reactive<Map<String, DataTableRowKey[]>> = reactive(new Map<String, DataTableRowKey[]>);

const headerColumns: Array<DataTableColumn<entity.PredefineTableHeader>> = [
  { type: "selection" },
  { title: t('column.name'), key: "Name", width: "200px" },
  {
    title: t('column.tag'), key: "Tag",
    width: "75px",
    render(row: entity.PredefineTableHeader) {
      let tagColorProp: TagColor = {};
      if (row.TagColor != '') {
        tagColorProp.color = row.TagColor;
      }
      if (row.TagTextColor != '') {
        tagColorProp.textColor = row.TagTextColor;
      }
      return h(
        NTag,
        { color: tagColorProp },
        { default: () => row.Symbol == '' ? '/' : row.Symbol },
      )
    }
  },
  { title: t('column.url'), key: "HeaderUrl" },
];
let schemeData: Array<entity.PredefineTableScheme> = [];

function loadSchemes() {
  loading.value = true;
  QueryPredefineTableSchemes()
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      const rows: Array<entity.PredefineTableScheme> = result.Rows;
      schemeData = [...rows];
      schemeOptions.value = [...rows.map(scheme => {
        return {
          key: scheme.Name,
          value: scheme.Name,
          label: scheme.Name
        };
      })];
      currentScheme.value = schemeOptions.value[0].value;
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function handleSubmit() {
  loading.value = true;
  const headerURLMap = new Map<string, entity.PredefineTableHeader>();
  currentHeaders.forEach(header => {
    headerURLMap.set(header.HeaderUrl, header);
  });
  // The row key is the header's url, we need to re-construct the param by using it
  const candidates: vo.DiffTableHeaderVo[] = [];
  checkedHeaderKey.forEach((headers) => {
    headers.forEach((headerURL: string) => {
      const rawHeaderDef = headerURLMap.get(headerURL);
      candidates.push({
        HeaderUrl: headerURL,
        name: rawHeaderDef.Name,
        TagColor: rawHeaderDef.TagColor,
        TagTextColor: rawHeaderDef.TagTextColor
      } as any)
    })
  });
  if (candidates.length == 0) {
    createSkipModal();
    loading.value = false;
    return;
  }
  AddBatchDiffTableHeader(candidates)
    .then(result => {
      if (result.Code != 200) {
        return Promise.reject(result.Msg);
      }
      const failedTables: Array<vo.DiffTableHeaderVo> = result.Rows;
      console.log(failedTables);
      const failedMessage = failedTables.map((header: vo.DiffTableHeaderVo) => {
        return header.name;
      }).join(";");
      console.log(failedMessage);
      if (failedTables.length > 0) {
        modal.create({
          type: "error",
          title: t('failedModal.title'),
          preset: "dialog",
          content: t('failedModal.content', { tables: failedMessage }),
          closable: true,
        });
      }
      props.moveOn();
    })
    .catch(err => window.$notifyError(err))
    .finally(() => loading.value = false);
}

function handleSkip() {
  createSkipModal();
}

function createSkipModal() {
  modal.create({
    title: t('skipModal.title'),
    preset: "dialog",
    content: t('skipModal.content'),
    positiveText: t('button.positive'),
    negativeText: t('button.negative'),
    onPositiveClick: () => {
      props.moveOn();
    },
    onNegativeClick: () => {
      // Do nothing
    }
  })
}

watch(currentScheme, (newScheme) => {
  const scheme = schemeData.find(scheme => scheme.Name == newScheme);
  // Each scheme is an array of headers along with some meta data
  // We have to split the headers into multiple groups
  currentHeaders = scheme.Headers;
  const groupedHeaders = new Map<string, Array<entity.PredefineTableHeader>>();
  currentHeaders.forEach(header => {
    if (!groupedHeaders.has(header.Category)) {
      groupedHeaders.set(header.Category, []);
    }
    groupedHeaders.get(header.Category).push(header);
  });
  // Convert grouped map into table rows
  categoryData.value = [];
  groupedHeaders.forEach((headers, categoryName) => {
    categoryData.value.push({
      key: categoryName,
      name: categoryName,
      headers: headers
    });
  });
});

watch(checkedHeaderKey, (newValue) => {
  console.log(newValue);
});

loadSchemes();
</script>

<i18n lang="json">{
  "en": {
    "column": {
      "name": "Name",
      "tag": "Tag",
      "url": "URL"
    },
    "button": {
      "skip": "Skip",
      "submit": "Submit",
      "positive": "Yes",
      "negative": "Cancel"
    },
    "skipModal": {
      "title": "Skip Table Initialization",
      "content": "Lampghost might work abnormally with no tables. Do you really want to skip this step?"
    },
    "failedModal": {
      "title": "Failed",
      "content": "Some tables are failed to add, you could try adding them manually: {tables}"
    }
  },
  "zh-CN": {
    "column": {
      "name": "名称",
      "tag": "难度表标签",
      "url": "地址"
    },
    "button": {
      "skip": "跳过",
      "submit": "提交",
      "positive": "确定",
      "negative": "取消"
    },
    "skipModal": {
      "title": "跳过难度表新增",
      "content": "Lampghost在一个难度表都没有的时候可能产生一些奇怪的现象。你确定要跳过新增步骤吗?"
    },
    "failedModal": {
      "title": "失败",
      "content": "部分难度表新增失败,你可以稍后手动添加他们: {tables}"
    }
  }
}</i18n>
