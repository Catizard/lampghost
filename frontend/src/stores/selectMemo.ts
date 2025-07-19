import { defineStore } from "pinia";

export type Memo = {
  difficultTableId?: number
  customTableId?: number
};

export const useSelectMemo = defineStore("selectMemo", {
  state: () => ({
    difficultTableId: null,
    customTableId: null
  } as Memo),

  actions: {
    setDifficultTableId(id: number | null) {
      this.difficultTableId = id
    },
    setCustomTableId(id: number | null) {
      this.customTableId = id;
    }
  }
});
