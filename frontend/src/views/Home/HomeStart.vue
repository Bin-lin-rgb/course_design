<template>
  <CommHeader />

  <div>
  ddd
  </div>
</template>

<script lang="ts" setup>
import CommHeader from "@/components/common/CommHeader.vue";
import { ref, onBeforeMount, onMounted } from 'vue';
import { getUserInfo, updateScore } from "@/utils/api/login";

//用户是否是登录状态
const isLogin = ref(false);
//显示用户更多数据
const userInfo = ref();

onBeforeMount(() => {
  getUserInfo().then((res: any) => {
    console.log(res);
    if (res.code == 0) {
      isLogin.value = true;
      userInfo.value = res.result;
    }
  });
});
</script>
