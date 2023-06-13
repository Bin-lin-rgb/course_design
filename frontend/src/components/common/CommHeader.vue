<template>
  <header>
    <div class="header-content">
      <h1 class="content-logo">
        <img src="@/assets/DrinkUpTheSea.png" @click="gotoHome" />
      </h1>
      <div class="content-nav">
        <ul>
          <router-link class="link" to="/">关于我们</router-link>
          <li>课程</li>
          <li>直播</li>
          <li>活动</li>
        </ul>
      </div>
      <div class="search-login">
        <div class="content-search">
          <input type="" placeholder="请输入要搜索的内容" />
          <el-icon color="#808080" :size="22">
            <search style="width: 22px; height: 22px" />
          </el-icon>
        </div>
        <div class="content-create">
          <button class="cta" @click="gotoCreator">
            <span>创作者中心</span>
            <svg viewBox="0 0 13 10" height="10px" width="15px">
              <path d="M1,5 L11,5"></path>
              <polyline points="8 1 12 5 8 9"></polyline>
            </svg>
          </button>
        </div>
        <div class="content-BellFilled">
          <el-icon :size="25" color="#8A919F"><BellFilled /></el-icon>
        </div>
        <div class="content-login" v-if="!isLogin">
          <router-link to="/login">登录 / 注册</router-link>
        </div>
        <!-- 划过头像显示  -->
        <el-dropdown v-else>
          <span class="el-dropdown-link">
            <el-avatar
              @mouseenter="isShow = true"
              src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
          /></span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="gotoWrite">写文章</el-dropdown-item>
              <el-dropdown-item @click="gotoUserinfo"
                >我的主页</el-dropdown-item
              >
              <el-dropdown-item @click="GotoEdit">设置</el-dropdown-item>
              <el-dropdown-item @click="Logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { Search, BellFilled } from "@element-plus/icons-vue";
// pinia
import { useUserStore } from "@/stores/user";
// Router
import { useRouter } from "vue-router";
// api
import { getUserInfo } from "@/utils/api/login";
//element
import { ElMessage, ElMessageBox } from "element-plus";

const router = useRouter();
//用户是否是登录状态
const isLogin = ref(false);
// 划过头像展示
const isShow = ref(false);
//显示用户更多数据
const userInfo = ref();

function Logout() {
  ElMessageBox.confirm("确定退出登录吗？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      const userStore = useUserStore();
      userStore.clearToken();
      router.go(0);
      ElMessage({
        type: "success",
        message: "退出成功！",
      });
    })
    .catch(() => {
      ElMessage({
        type: "info",
        message: "已取消",
      });
    });
}

onBeforeMount(() => {
  getUserInfo().then((res: any) => {
    // console.log(res);
    if (res.code == 0) {
      isLogin.value = true;
      userInfo.value = res.data;
    }
  });
});

function gotoHome() {
  router.push("/");
}
function gotoCreator() {
  router.push("/creator");
}
function gotoUserinfo() {
  router.push("/userinfo");
}
function gotoWrite() {
  const routeData = router.resolve({
    path: "/write",
  });
  window.open(routeData.href, "_blank");
}
function GotoEdit() {
  router.push("/userinfo/setting/profile");
}
</script>

<style scoped>
header {
  position: fixed;
  top: 0;
  left: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 70px;
  background: white;
  box-shadow: 0px 5px 6px rgba(0, 0, 0, 0.16);
  opacity: 1;
  z-index: 100;
}

.header-content {
  display: flex;
  justify-content: space-around;
  width: 1700px;
}

.content-logo {
  width: 100px;
  height: 30px;
  margin: 20px 0;
  cursor: pointer;
}

.content-logo img {
  height: 100%;
}

.content-nav {
  width: 300px;
  height: 75px;
}

.content-nav ul {
  display: flex;
  justify-content: space-around;
  align-items: center;
  width: 100%;
  height: 70px;
}

.content-nav ul li {
  font-size: 16px;
  color: #323232;
}

.content-nav ul .link {
  font-size: 16px;
  color: #323232;
}

.search-login {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 650px;
}

.content-search {
  display: flex;
  align-items: center;
  padding: 5px 10px;
  width: 350px;
  height: 35px;
  border-radius: 8px;
  background: #f0f2f4;
}

.content-search input {
  padding: 0 10px;
  width: 430px;
  height: 40px;
  border: 0;
  border-radius: 8px;
  color: #808080;
  background: #f0f2f4;
  font-size: 16px;
  outline: none;
}

.content-BellFilled {
  /* font-size: 23px; */
  margin-right: 23px;
}

/* From uiverse.io by @alexmaracinaru */
.cta {
  position: relative;
  margin: auto;
  padding: 12px 18px;
  transition: all 0.2s ease;
  border: none;
  background: none;
}

.cta:before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  display: block;
  border-radius: 50px;
  background: #b1dae7;
  width: 45px;
  height: 45px;
  transition: all 0.3s ease;
}

.cta span {
  position: relative;
  font-family: "Ubuntu", sans-serif;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.05em;
  color: #234567;
}

.cta svg {
  position: relative;
  top: 0;
  margin-left: 10px;
  fill: none;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke: #234567;
  stroke-width: 2;
  transform: translateX(-5px);
  transition: all 0.3s ease;
}

.cta:hover:before {
  width: 100%;
  background: #b1dae7;
}

.cta:hover svg {
  transform: translateX(0);
}

.cta:active {
  transform: scale(0.95);
}

.content-login {
  font-size: 18px;
  color: #808080;
  text-align: center;
  cursor: pointer;
}

#avatar {
  position: fixed;
  right: 20px;
}
</style>
