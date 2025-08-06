import { dto } from "@wailsjs/go/models";
import { defineStore } from "pinia";

export const useUserStore = defineStore("userStore", {
  state: () => {
    return {
      user: undefined as dto.RivalInfoDto | undefined
    }
  },
  getters: {
    id: (state) => state.user.ID,
  },
  actions: {
    setter(user: dto.RivalInfoDto) {
      this.user = user;
    }
  }
})
