import dayjs from 'dayjs';
import { defineStore } from 'pinia';

export const useSpecifyYearStore = defineStore("specifyYear", {
  state: () => ({
    specifyYear: dayjs().get('year').toString(),
  }),

  actions: {
    setter(newValue: string) {
      this.specifyYear = newValue;
    }
  }
})
