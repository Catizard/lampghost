export enum ScoreOption {
  NORMAL,
  MIRROR,
  RANDOM,
  R_RANDOM,
  S_RANDOM,
  SPIRAL,
  H_RANDOM,
  ALL_SCR,
  RANDOM_EX,
  S_RANDOM_EX
}

export class ScoreOptionDef {
  value: number;
  text: string;
  color: string;
}

export type ScoreOptionColorStyle = Record<number, ScoreOptionDef>;

export const DefaultScoreOptionColorStyle: Record<number, ScoreOptionDef> = {
  0: {
    value: 0,
    text: "NORMAL",
    color: "#FFFFFF",
  },
  1: {
    value: 1,
    text: "MIRROR",
    color: "#FFFFFF"
  },
  2: {
    value: 2,
    text: "RANDOM",
    color: "#FFFFFF",
  },
  3: {
    value: 3,
    text: "R-RANDOM",
    color: "#FFFFFF"
  },
  4: {
    value: 4,
    text: "S-RANDOM",
    color: "#FFFFFF",
  },
  5: {
    value: 5,
    text: "SPIRAL",
    color: "#FFFFFF"
  },
  6: {
    value: 6,
    text: "H-RANDOM",
    color: "#FFFFFF",
  },
  7: {
    value: 7,
    text: "ALL-SCR",
    color: "#FFFFFF"
  },
  8: {
    value: 8,
    text: "RANDOM-EX",
    color: "#FFFFFF",
  },
  9: {
    value: 9,
    text: "S-RANDOM-EX",
    color: "#FFFFFF"
  },
}

export function queryScoreOptionColorStyle(scoreOption: number): ScoreOptionDef {
  return DefaultScoreOptionColorStyle[scoreOption];
}
