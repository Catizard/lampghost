export enum ScoreRank {
  NO_PLAY = 0,
  F,
  E,
  D,
  C,
  B,
  A,
  AA,
  AAA,
}

export class ScoreRankDef {
  value: number;
  text: string;
  color: string;
  textColor: string;
}

export type ScoreRankColorStyle = Record<number, ScoreRankDef>;

export const DefaultScoreRankColorStyle: Record<number, ScoreRankDef> = {
  0: {
    value: 0,
    text: "NO_PLAY",
    color: "#8C8C8C",
    textColor: "#FFFFFF"
  },
  1: {
    value: 1,
    text: "F",
    color: "#CC5C76",
    textColor: "#ffffff",
  },
  2: {
    value: 2,
    text: "E",
    color: "#FF9FF9",
    textColor: "#ffffff",
  },
  3: {
    value: 3,
    text: "D",
    color: "#FF9FF9",
    textColor: "#ffffff",
  },
  4: {
    value: 4,
    text: "C",
    color: "#49E670",
    textColor: "#ffffff",
  },
  5: {
    value: 5,
    text: "B",
    color: "#4FBCF7",
    textColor: "#ffffff",
  },
  6: {
    value: 6,
    text: "A",
    color: "#FF6B74",
    textColor: "#ffffff",
  },
  7: {
    value: 7,
    text: "AA",
    color: "#FFAD70",
    textColor: "#ffffff",
  },
  8: {
    value: 8,
    text: "AAA",
    color: "#FFD251",
    textColor: "#ffffff",
  },
}
