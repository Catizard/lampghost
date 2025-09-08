import { config } from "@wailsjs/go/models";
import { defineStore } from "pinia";

export const useConfigStore = defineStore("config", {
  state: () => {
    return {
      config: undefined as config.ApplicationConfig | undefined
    }
  },
  actions: {
    setter(config: config.ApplicationConfig) {
      this.config = config;
    }
  }
})
