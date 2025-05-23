export enum ClearType {
  NO_PLAY = 0,
  Failed,
  AssistEasy,
  LightAssistEasy,
  Easy,
  Normal,
  Hard,
  ExHard,
  FullCombo,
  Perfect,
  Max,
}

export class ClearTypeDef {
  value: number;
  text: string;
  color: string;
  textColor: string;
}

export type ClearTypeColorStyle = Record<number, ClearTypeDef>;

export const DefaultClearTypeColorStyle: Record<number, ClearTypeDef> = {
  0: {
    value: 0,
    text: "NO_PLAY",
    color: "#8C8C8C",
    textColor: "#FFFFFF",
  },
  1: {
    value: 1,
    text: "FAILED",
    color: "#CC5C76",
    textColor: "#ffffff",
  },
  2: {
    value: 2,
    text: "A-Easy",
    color: "#FF9FF9",
    textColor: "#ffffff",
  },
  3: {
    value: 3,
    text: "LA-Easy",
    color: "#FF9FF9",
    textColor: "#ffffff",
  },
  4: {
    value: 4,
    text: "EASY",
    color: "#49E670",
    textColor: "#ffffff",
  },
  5: {
    value: 5,
    text: "NORMAL",
    color: "#4FBCF7",
    textColor: "#ffffff",
  },
  6: {
    value: 6,
    text: "HARD",
    color: "#FF6B74",
    textColor: "#ffffff",
  },
  7: {
    value: 7,
    text: "EX-HARD",
    color: "#FFAD70",
    textColor: "#ffffff",
  },
  8: {
    value: 8,
    text: "FULL-COMBO",
    color: "#FFD251",
    textColor: "#ffffff",
  },
  9: {
    value: 9,
    text: "PERFECT",
    color: "#FFD251",
    textColor: "#ffffff",
  },
  10: {
    value: 10,
    text: "MAX",
    color: "#FFD251",
    textColor: "#ffffff",
  },
};

export function queryClearTypeColorStyle(clearType: number): ClearTypeDef {
  return DefaultClearTypeColorStyle[clearType];
}
