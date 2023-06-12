import { createRouter, createWebHistory } from "vue-router";
import Home from "@/views/Home/HomeIndex.vue";

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
      name: "Home",
      component: Home,
      meta: {
        isAuth: false,
        title: "Blog",
      },
    },
    {
      path: "/login",
      name: "Login",
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("@/views/UserInfo/UserLogin.vue"),
      meta: {
        isAuth: false,
        title: "登录",
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
