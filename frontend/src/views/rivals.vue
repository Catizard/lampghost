<template>
	<perfect-scrollbar>
		<n-flex justify="space-between">
			<n-h1 prefix="bar" style="text-align: start;">
				<n-text type="primary">Rivals</n-text>
			</n-h1>
			<n-flex justify="end">
				<n-button>Reload All</n-button>
			</n-flex>

			<n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false" :row-key="row => row.ID" />
		</n-flex>
	</perfect-scrollbar>
</template>

<script lang="ts" setup>
import { DataTableColumns } from 'naive-ui';
import { RivalInfo } from '../../wailsjs/go/models';

const pagination = false as const;

function createColumns(): DataTableColumns<dto.RivalInfo> {

}
</script>