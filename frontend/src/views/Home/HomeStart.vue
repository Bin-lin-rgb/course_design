<template>
  <CommHeader />

  <div v-if="!isLogin">
    请先登录！
    <a-divider />
    <div class="bottom-save">
      <a-button type="primary" size="large" @click="gotoLogin"
        >去登录</a-button
      >
    </div>
  </div>

  <div class="main-wrap" v-if="isLogin">
    <div class="font-2">个人资料</div>
    <a-divider />
    <div class="input-item">
      <div class="item-text-area">用户名</div>
      <el-input
        v-model="data.username"
        maxlength="20"
        placeholder="请填写你的用户名"
        show-word-limit
        type="text"
        readonly="true"
      />
    </div>
    <a-divider />
    <div class="input-item">
      <div class="item-text-area">学号</div>
      <el-input
        v-model="data.account"
        max="710"
        min="0"
        placeholder="请填写你的职位"
        show-word-limit
        type="text"
        readonly="true"
      />
    </div>
    <a-divider />
    <div>
      <!-- 四级成绩 -->    
      <div class="input-item">
      <div class="item-text-area">四级成绩</div>

      <el-select 
        v-if="data.fourGrade === ''" v-model="fourValue" class="m-2" placeholder="Select" size="large">

    <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    />
  </el-select>

      <el-input
        v-if="fourValue === '1000' || data.fourGrade !== ''"
        v-model="data.fourGrade"
        placeholder="请填写你的四级成绩"
        show-word-limit
        type="number"
      />
    </div>
    </div>

    <a-divider />

    <div>
      <!-- 六级成绩 -->    
      <div class="input-item">
      <div class="item-text-area">六级成绩</div>

      <el-select v-if="data.sixGrade === ''" v-model="sixValue" class="m-2" placeholder="Select" size="large">
    <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    />
  </el-select>

      <el-input
        v-if="sixValue === '1000' || data.sixGrade !== ''"
        v-model="data.sixGrade"
        placeholder="请填写你的六级成绩"
        show-word-limit
        type="number"
      />
    </div>
    </div>
    <a-divider />
    <div class="input-item">
      <div class="item-text-area">最近测试结果</div>
      <el-input
        v-model="data.basicVocabulary"
        maxlength="100"
        placeholder="请填写你的个人介绍"
        show-word-limit
        type="textarea"
        readonly="true"
      />
    </div>

    <div class="bottom-save">
      <a-button type="primary" size="large" @click="HandleCommit"
        >保存修改，然后开测</a-button
      >
    </div>

    <div class="bottom-save">
      <a-button type="primary" size="large" @click="gotoList1"
        >不修改，直接开始测试</a-button
      >
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ElMessage } from "element-plus";
import CommHeader from "@/components/common/CommHeader.vue";
import { ref, onBeforeMount, reactive } from 'vue';
import { getUserInfo, updateScore } from "@/utils/api/login";
import router from "@/router";

//用户是否是登录状态
const isLogin = ref(false);

const data = reactive({
  username: "",
  account: "",
  fourGrade: "",
  sixGrade: "",
  basicVocabulary: "",
});

const data2 = reactive({
  fourGrade: "",
  sixGrade: "",
});

const fourValue = ref('')
const sixValue = ref('')

const options = [
  {
    value: '还没考呢',
    label: '还没考呢',
  },
  {
    value: '考了没过',
    label: '考了没过',
  },
  {
    value: '1000',
    label: '考了过了',
  }
]

const HandleCommit = async () => {
  if (fourValue.value !== '1000') {
    data2.fourGrade = fourValue.value
  }
  if (sixValue.value !== '1000') {
    data2.sixGrade = sixValue.value
  }
  if (data.fourGrade != ''){
    data2.fourGrade = data.fourGrade
  }
  if (data.sixGrade != ''){
    data2.sixGrade = data.sixGrade
  }

  console.log(data2);

  const res: any = await updateScore(data2);

  console.log(res);
  if (res.code === 0) {
    ElMessage({
      message: "修改成功",
      type: "success",
    });
    gotoList1();
  }
};

const gotoLogin = () => {
  router.push("/login");
};

const gotoList1 = () => {
  router.push("/index");
};

const loadUserInfo = async () => {
  const res: any = await getUserInfo();
  if (res.code === 0) {
    isLogin.value = true;
    data.username = res.result.username;
    data.account = res.result.account;
    data.fourGrade = res.result.fourGrade;
    data.sixGrade = res.result.sixGrade;
    data.basicVocabulary = res.result.basicVocabulary;
  } else {
    return;
  }
};

onBeforeMount(() => {
  loadUserInfo();
});
</script>


<style scoped lang="less">
.main-wrap {
  width: 700px;
  .font-2 {
    font-size: large;
  }
  .input-item {
    display: flex;
    .item-text-area {
      width: 100px;
      height: 35px;
      text-align: center;
      line-height: 35px;
    }
  }
  .bottom-save {
    margin: 50px;
  }
}
</style>
