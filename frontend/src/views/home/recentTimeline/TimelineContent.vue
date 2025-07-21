<template>
  <n-timeline>
    <n-timeline-item v-for="node in nodes" :key="node.key" v-bind="dynamicNodeProps(node)">
      <template #header>
        <n-text class="timeline-node" :style="{ 'color': node.color }">
          {{ node.title }}
        </n-text>
      </template>
      <template #default>
        <div v-for="log in node.logs" :key="log.ID" style="margin-top: 10px;">
          <template v-if="log.TableTags.length > 0">
            <n-tag v-for="dTag in log.TableTags" :key="dTag.TableName + dTag.TableLevel" size="small"
              style="margin-right: 5px" v-bind="dynamicDTagProps(dTag)">
              {{ dTag.TableSymbol + dTag.TableLevel }}
            </n-tag>
          </template>
          <template v-if="log.Title == ''">
            <n-icon :component="WarningOutline" color="red" />
            Missing Song
          </template>
          <template v-else>
            {{ log.Title }}
          </template>
        </div>
      </template>
    </n-timeline-item>
  </n-timeline>
</template>

<script lang="ts" setup>
import { dto } from '@wailsjs/go/models';
import { Node } from './node';
import { WarningOutline } from "@vicons/ionicons5";
import { NTimeline } from 'naive-ui';

const props = defineProps<{
  nodes: Node[]
}>();

// dynamic v-bind
function dynamicNodeProps(node: Node) {
  if (node.type != null) {
    return {
      type: node.type,
    } as any // make typescript chill...
  }
  return {
  }
}

// dynamic v-bind
function dynamicDTagProps(dTag: dto.DiffTableTagDto) {
  let props: any = {}
  let hasProp = false;
  if (dTag.TableTagColor.length > 0) {
    props.color = dTag.TableTagColor;
    hasProp = true;
  }
  if (dTag.TableTagTextColor.length > 0) {
    props.textColor = dTag.TableTagTextColor
    hasProp = true;
  }
  if (!hasProp) {
    return {}
  }
  return {
    color: props
  }
}

</script>

<style scoped>
.timeline-node {
  font-weight: bold;
  font-size: 20px;
}
</style>
