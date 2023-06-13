<template>
  <header>
    <div class="header-content">
      <h1 class="content-logo">
        <img src="@/assets/DrinkUpTheSea.png" @click="gotoHome" />
      </h1>

      <div class="bell-user">
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
              <el-dropdown-item @click="gotoUserinfo">我的主页</el-dropdown-item>
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
import { BellFilled } from "@element-plus/icons-vue";
// pinia
import { useUserStore } from "@/stores/user";
// api
import { getUserInfo } from "@/utils/api/login";
//element
import { ElMessage, ElMessageBox } from "element-plus";
// Router
import { useRouter } from "vue-router";
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
    if (res.code == 1000) {
      isLogin.value = true;
      userInfo.value = res.data;
    }
  });
});

function gotoHome() {
  router.push("/");
}
function gotoUserinfo() {
  router.push("/userinfo");
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
  width: 110px;
  height: 30px;
  margin: 20px 900px 20px 0px;
  cursor: pointer;
  font-size: 24px;
}

.content-logo img {
  height: 100%;
}

.bell-user {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 120px;
}

/* .title {
  font-size: 24px;
  box-sizing: border-box;
} */
</style>
