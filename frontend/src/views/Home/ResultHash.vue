<template>
    <router-view></router-view>
    <CommHeader />
    <div class="common-layout">
        <el-card class="box-card">
            <el-row>
                <el-col :span="8">
                    <el-form-item>
                        <h1 class="Tips">Awesome! We calculated your vocabulary </h1>
                        <el-divider />
                    </el-form-item>

                </el-col>
                <el-col :span="8">
                </el-col>
                <el-col :span="6">
                    <el-card class="number">
                        <h1 class="result"> {{ parsedData }}</h1>
                    </el-card>
                </el-col>
            </el-row>
        </el-card>
    </div>
</template>
  
<script lang="ts" setup>
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import CommHeader from "@/components/common/CommHeader.vue";

const route = useRoute();
const parsedData = ref(parseInt(route.query.data as string) || 0);

</script>
<style scoped lang="less">
.common-layout {
    background-color: rgb(228, 237, 237);
    background-size: cover;
    background-position: center;
    width: 8000px;
    height: 600px;
    margin-top: 120px;
    border-radius: 50px;
}

.Tips {
    margin-top: 100px;
    font-size: 40px;
    margin-left: 50px;
    margin-bottom: 80px;
    line-height: 1.4;
}

.number {
    margin-top: 100px;
    margin-left: 50px;
    margin-bottom: 80px;
}

.result {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 40px;
}
</style>
  