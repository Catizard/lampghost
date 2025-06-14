<template>
	<n-infinite-scroll :loading="loading" style="height: 90%;" :distance="10" @load="handleLoad">
		<n-timeline>
			<n-timeline-item v-for="node in timelineNodes" :key="node.key" v-bind="dynamicNodeProps(node)">
				<template #header>
					<n-text class="timeline-node" :style="{ 'color': node.color }">
						{{ node.title }}
					</n-text>
				</template>
				<!-- <template v-for="([clearType, value]) in dailyLog.summary" :key="clearType">
            <ClearTag :clear="clearType" @click="handleFocusClearTag(dailyLog, clearType)" /> {{ value }}
          </template> -->
				<template #default>
					<div v-for="log in node.logs" :key="log.ID" style="margin-top: 10px;">
						<template v-if="log.TableTags.length > 0">
							<n-tag v-for="dTag in log.TableTags" :key="dTag.TableName + dTag.TableLevel" size="small"
								style="margin-right: 5px" v-bind="dynamicDTagProps(dTag)">
								{{ dTag.TableSymbol + dTag.TableLevel }}
							</n-tag>
						</template>
						<template v-if="log.Title == ''">
							<n-icon :component="WarningOutline" color="red"/>
							Missing Song
						</template>
						<template v-else>
							{{ log.Title }}
						</template>
					</div>
				</template>
			</n-timeline-item>
		</n-timeline>
	</n-infinite-scroll>
</template>

<script lang="ts" setup>
import { ClearType, ClearTypeDef, queryClearTypeColorStyle } from '@/constants/cleartype';
import { QueryPrevDayScoreLogList } from '@wailsjs/go/main/App';
import { dto } from '@wailsjs/go/models';
import dayjs from 'dayjs';
import { ref, Ref, watch, nextTick } from 'vue';
import { WarningOutline } from '@vicons/ionicons5';

const props = defineProps<{
	rivalId?: number;
}>();


type Node = {
	key: number, // unique id, required by vue's v-for
	type: string | null,
	title: string,
	color: string, // hex color
	logs: Array<dto.RivalScoreLogDto>,
}

const loading = ref(false);
let nodeId = 0;
const timelineNodes: Ref<Array<Node>> = ref([]);
const noMore = ref(false);
let queryDate = dayjs().add(1, 'day');

async function handleLoad() {
	if (loading.value || noMore.value) {
		return;
	}
	loading.value = true;
	try {
		const result = await QueryPrevDayScoreLogList({
			EndRecordTimestamp: queryDate.unix(),
			RivalId: props.rivalId,
		} as any);
		if (result.Code != 200) {
			throw new Error(result.Msg);
		}
		const nextLogs: Array<dto.RivalScoreLogDto> = result.Rows;
		if (nextLogs.length == 0) {
			console.log('no more values')
			noMore.value = true;
		} else {
			console.log('next logs: ', nextLogs);
			queryDate = dayjs(nextLogs[0].RecordTime);
			timelineNodes.value.push(buildTimelineDateNode(queryDate.format('YYYY-MM-DD')));
			for (let i = 0; i < nextLogs.length; ++i) {
				let j = i;
				while (j + 1 < nextLogs.length && nextLogs[j + 1].Clear == nextLogs[i].Clear) j++;
				const clearTypeDef: ClearTypeDef = queryClearTypeColorStyle(nextLogs[i].Clear);
				let clearSummaryNode: Node = buildTimelineSummaryNode(j - i + 1, clearTypeDef);
				for (let k = i; k <= j; ++k) {
					clearSummaryNode.logs.push(nextLogs[k]);
				}
				i = j;
				timelineNodes.value.push(clearSummaryNode);
			}
			loading.value = false;
			// Hack: To prevent having a very long white space, keep load data if log count < 15
			// There might exist a rare race condition...but I hardly think this is damaging
			if (timelineNodes.value.length < 15) {
				nextTick(() => {
					handleLoad();
				});
			}
		}
	} catch (err) {
		window.$notifyError(err);
	} 
	loading.value = false;
}

// Generate a date node
function buildTimelineDateNode(date: string): Node {
	return {
		title: date,
		type: "info",
		key: nodeId++,
		logs: [],
		color: "#000000",
	}
}

// Generate a lamp clear summary node
function buildTimelineSummaryNode(count: number, clearTypeDef: ClearTypeDef): Node {
	const title = clearTypeDef.value == ClearType.FullCombo
		? `New ${count} ${clearTypeDef.text}`
		: `New ${count} ${clearTypeDef.text} Clear`
	return {
		title: title,
		key: nodeId++,
		type: "success",
		logs: [],
		color: clearTypeDef.color,
	}
}

function handleFocusClearTag(node: Node, clearType: number) {
	console.log('node: ', node, 'clearType: ', clearType);
	// TODO: implement me!
}

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

watch(() => props.rivalId, (_, newValue) => {
	handleLoad();
}, { once: true });

</script>

<style lang="css" scoped>
.timeline-node {
	font-weight: bold;
	font-size: 20px;
}
</style>
