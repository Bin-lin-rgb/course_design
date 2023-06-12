<template>
  <a-result status="403" subtitle="抱歉，暂不支持该功能.">
    <template #extra>
      <a-space>
        <a-button type="primary" @click="Back">返回</a-button>
      </a-space>
    </template>
  </a-result>
</template>

<script setup lang="ts">
// router
import { useRouter } from "vue-router";

const router = useRouter();
function Back() {
  router.go(-1);
}
</script>
