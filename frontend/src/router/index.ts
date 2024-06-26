import { createRouter, createWebHistory } from "vue-router";
import Home from "@/views/Home/HomeIndex.vue";
import HomeStart from "@/views/Home/HomeStart.vue";
import Vocabulary2 from "@/views/Home/Vocabulary2.vue";
import ResultHash from "@/views/Home/ResultHash.vue";
// class meta {
//   isAuth: boolean;
//   title: string;
//   // 构造方法
//   constructor(isAuth: boolean, title: string) {
//     this.isAuth = isAuth;
//     this.title = title;
//   }
// }

// const creator = new meta(true, "创造者中心");

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "HomeStart",
      component: HomeStart,
      meta: {
        isAuth: false,
        title: "词汇量测试首页",
      },
    },
    {
      path: "/index",
      name: "Home",
      component: Home,
      meta: {
        isAuth: true,
        title: "词汇量测试列表1",
      },
    },
    {
      path: "/index2",
      name: "Home2",
      component: Vocabulary2,
      props:true,
      meta: {
        isAuth: true,
        title: "词汇量测试列表2",
      },
    },
    {
      path: "/resulthash",
      name: "resulthash",
      component: ResultHash,
      props:true,
      meta: {
        isAuth: true,
        title: "词汇量测试结果",
      },
    },
    {
      path: "/login",
      name: "Login",
      component: () => import("@/views/UserInfo/UserLogin.vue"),
      meta: {
        isAuth: false,
        title: "词汇量测试->登录",
      },
    },
    {
      path: "/NoFunction",
      name: "NoFunction",
      component: () => import("@/components/common/ErrPage.vue"),
      meta: {
        isAuth: false,
        title: "额哦,403啦",
      },
    },
  ],
});

// pinia
import { useUserStore } from "@/stores/user";
router.beforeEach((to, from, next) => {
  const userStore = useUserStore();
  if (to.meta.isAuth === "ture" && userStore.token !== "")
    next({ name: "Login" });
  else next();
});

//全局后置路由守卫————初始化的时候被调用、每次路由切换之后被调用
router.afterEach((to, from) => {
  let str = "Blog";
  if (String(to.meta.title) !== "") {
    str = String(to.meta.title);
  }
  document.title = str;
});

export default router;
