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
    isLR2User: (state) => state.user.Type == "LR2"
  },
  actions: {
    setter(user: dto.RivalInfoDto) {
      this.user = user;
    }
  }
})
