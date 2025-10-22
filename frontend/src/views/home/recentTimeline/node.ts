import { ClearType, ClearTypeDef, queryClearTypeColorStyle } from "@/constants/cleartype";
import { dto } from "@wailsjs/go/models";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
// HACK:See https://github.com/Catizard/lampghost/issues/28
dayjs.extend(utc);

// NOTE: No semantic meaning, used as the unique field of Node type
let nodeId = 0;

export type Node = {
  key: number, // unique id, required by vue's v-for
  type: string | null,
  title: string,
  color: string, // hex color
  logs: Array<dto.RivalScoreLogDto>,
};

export function parseOneDayLogs(rows: dto.RivalScoreLogDto[]): [dayjs.Dayjs, Array<Node>] {
  let queryDate = dayjs(rows[0].RecordTime).utc();
  let nodes = [];
  nodes.push(buildTimelineDateNode(queryDate.format("YYYY-MM-DD")));
  for (let i = 0; i < rows.length; ++i) {
    let j = i;
    while (j + 1 < rows.length && rows[j + 1].Clear == rows[i].Clear) j++;
    const clearTypeDef = queryClearTypeColorStyle(rows[i].Clear);
    let clearSummaryNode = buildTimelineSummaryNode(j - i + 1, clearTypeDef);
    for (let k = i; k <= j; ++k) {
      clearSummaryNode.logs.push(rows[k]);
    }
    i = j;
    nodes.push(clearSummaryNode);
  }
  return [queryDate, nodes];
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
