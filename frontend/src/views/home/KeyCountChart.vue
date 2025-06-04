<template>
	<!-- Hack: using `:key=...` for updating heatmap forcefully -->
	<!-- Honestly, I have no idea why `endDate` is not reactive-->
	<calendar-heatmap :values="data" :endDate="endDate" :tooltip-formatter="tooltipFormatter" :no-data-text="'No Play'"
		style="width: 95%" :key="endDate" />
</template>

<script setup lang="ts">
import { QueryUserKeyCountInYear } from "@wailsjs/go/main/App";
import { dto } from "@wailsjs/go/models";
import { computed, ref, Ref, watch } from "vue";
import dayjs from "dayjs";
import "vue3-calendar-heatmap/dist/style.css";
import { useSpecifyYearStore } from "@/stores/home";

const props = defineProps<{
	rivalId?: number;
}>();

const specifyYearStore = useSpecifyYearStore();
const endDate = computed(() => {
	const current = dayjs();
	const currentYear = current.get('year').toString();
	const specifyYear = specifyYearStore.specifyYear;
	if (currentYear == specifyYear) {
		return current;
	}
	return dayjs(`${specifyYear}-12-31`);
})

const data: Ref<Array<any>> = ref([]);

function loadKeyCountData() {
	QueryUserKeyCountInYear({
		RivalId: props.rivalId,
	} as any)
		.then((result) => {
			if (result.Code != 200) {
				return Promise.reject(result.Msg);
			}
			console.log(result);
			data.value = [
				...result.Rows.map((row: dto.KeyCountDto) => {
					return {
						count: row.KeyCount,
						date: row.RecordDate,
					};
				}),
			];
			console.log(data.value);
		})
		.catch(err => window.$notifyError(err));
}

function tooltipFormatter(value: {
	date: Date | string;
	count: number;
}): string {
	const { date, count } = value;
	const mmdd = dayjs(date).format("YYYY-MM-DD");
	return `${count} notes <br> ${mmdd}`;
}

watch(
	() => props.rivalId,
	(_, newValue) => {
		loadKeyCountData();
	},
	{ once: true },
);
</script>
