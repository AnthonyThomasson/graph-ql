import { ApolloClients } from "@vue/apollo-composable";
import "vite/modulepreload-polyfill";
import { createApp, h, provide } from "vue";
import App from "./App.vue";
import { apolloClient } from "./apollo/client";
import "./style.css";

createApp({
  setup() {
    provide(ApolloClients, {
      default: apolloClient,
    });
  },

  render: () => h(App),
}).mount("#app");
