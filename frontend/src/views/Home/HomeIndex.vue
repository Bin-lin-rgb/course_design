<template>
  <CommHeader />
  <div class="common-layout">
    <el-container class="container">
      <el-main class="main"> hhhakhfiashd </el-main>
    </el-container>
  </div>
</template>

<script lang="ts" setup>
import CommHeader from "@/components/common/CommHeader.vue";
// import { GetPostListWithTime } from "@/utils/api/article";
// import { ref, onBeforeMount } from "vue";
// import { ElMessage } from "element-plus";
// import { useRouter } from "vue-router";
// const router = useRouter();
</script>

<style scoped lang="less">
@normal_text_color: #86909c;

.common-layout {
  margin-top: 120px !important;
  .container {
    width: 1100px;
    .main {
      padding: 10px;
      .list-wrap {
        background: rgb(255, 255, 255);
        height: 600px;
        .item {
          height: 125px;
          border-bottom: solid 1px #dcdfe6;
          padding: 10px;
          box-sizing: border-box;
          cursor: pointer;
          .info-bar {
            height: 30px;
            display: inline-block;
            color: @normal_text_color;
            .divider {
              width: 1px;
              height: 14px;
              background: #e5e6eb;
              margin: 0 8px;
              display: inline-block;
            }
          }
          .info-main {
            .main-title {
              font-weight: bold;
              font-size: 20px;
              margin-bottom: 12px;
              /* 文本过长自动截断 */
              overflow: hidden;
              white-space: nowrap;
              text-overflow: ellipsis;
            }
            .main-context {
              color: @normal_text_color;
              overflow: hidden;
              white-space: nowrap;
              text-overflow: ellipsis;
            }
          }
        }

        .item:hover {
          background-color: #fafafa;
        }
      }
    }
    .aside {
      padding: 10px;

      .aside-wrap {
        padding: 10px;
        background: rgb(255, 255, 255);
      }
    }
  }
}
</style>
