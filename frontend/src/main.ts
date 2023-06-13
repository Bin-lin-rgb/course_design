import { createApp } from "vue";

import App from "./App.vue";
import router from "./router";
// !!!记得要引用store啊
import store from "./stores";

import "./assets/main.css";

// arco-design
import ArcoVue from "@arco-design/web-vue";
// 额外引入图标库
import ArcoVueIcon from "@arco-design/web-vue/es/icon";
import "@arco-design/web-vue/dist/arco.css";

const app = createApp(App);

app.use(router);
app.use(store);
app.use(ArcoVue);
app.use(ArcoVueIcon);

app.mount("#app");
