import { defineStore } from "pinia";

export type Memo = {
  difficultTableId?: number
};

export const useSelectMemo = defineStore("selectMemo", {
  state: () => ({
    difficultTableId: null
  } as Memo),

  actions: {
    setDifficultTableId(id: number | null) {
      this.difficultTableId = id
    }
  }
});
