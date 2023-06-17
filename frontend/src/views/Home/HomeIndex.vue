<template>
  <CommHeader />
  <div class="common-layout">
    <el-card class="box-card">
      <el-form-item>
        <h1 class="Tips">Step 1 of 2: test your broad vocab level</h1>
        <el-divider />
      </el-form-item>
      <el-container class="container">
        <el-header>Check the box if you know at least one definition for a word. If you’re not sure about the exact
          meaning, leave it blank.</el-header>
      </el-container>
      <div>
        <div class="checkbox-row" v-for="(row, index) in checkboxRows" :key="index">
          <el-checkbox class="checkbox-item" v-for="item in row" :key="item.id" :model-value="item.checked"
            @update:model-value="updateValue(item, item.checked)">
            <span class="checkbox-label">{{ item.word }}</span>
          </el-checkbox>
        </div>
      </div>
      <el-button type="primary" @click="sendData">Continue</el-button>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import CommHeader from "@/components/common/CommHeader.vue";
import { ref, computed, onMounted } from 'vue';
import axios from "axios";
import { GetWordList1 } from '@/utils/api/wordBook';

import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';

const router = useRouter();
interface CheckboxOption {
  id: number;
  word: string;
  checked: boolean; // 添加选中状态的属性
}
interface ResultItem {
  id: number;
  word: string;
}//接收返回的数组
type ResultArray = ResultItem[];

const checkboxOptions = ref<CheckboxOption[]>([]);

onMounted(async () => {
  await getList1();
});

const getList1 = async () => {
  const res: any = await GetWordList1();
  checkboxOptions.value = res.result.map((item: { id: number; word: string }) => ({
    id: item.id,
    word: item.word,
    checked: false, // 初始化选中状态为false
  }));
};

const checkboxRows = computed(() => {
  const rows: CheckboxOption[][] = [];
  const options: CheckboxOption[] = [...checkboxOptions.value];
  while (options.length) {
    rows.push(options.splice(0, 4));
  }
  return rows;
});

const updateValue = (item: CheckboxOption, checked: boolean) => {
  item.checked = !item.checked; // 更新选中状态
};

const sendData = async () => {
  const selectedOptions = checkboxOptions.value.map((item) => ({
    id: item.id,
    word: item.word,
    known: item.checked ? 1 : 0, // 将布尔值转换为数值
  }));

  //console.log(selectedOptions); // 往后端传的数组
  const requestData = { WordList: selectedOptions };
  try {
    const userStore = useUserStore();
    const token = userStore.token;
    const response = await axios.post('/api/basic-api/word/postList1', requestData, {
      headers: { Authorization: `${token}` }

    });

    const ResResult: ResultArray = response.data.result;
    console.log("result",ResResult)
    router.push({
      name: 'Home2',
      params: { StringWordList: JSON.stringify(ResResult) }, // 将数组转换为字符串进行传递
    });
  } catch (error) {
    console.error(error);
  }
};

</script>

<style scoped lang="less">
@normal_text_color: #86909c;

.Tips {
  margin-top: 10px;
}

.el-card {
  margin-top: 50px;
}

.checkbox-row {
  display: flex;
  left: 5%;
}

.el-button {
  margin-left: 500px;
}

.checkbox-item {
  flex: 1;
  margin-right: 10px;
  zoom: 150%;
  align-items: center;
}

.checkbox-label {
  font-size: 12px;
}

.common-layout {
  background-color: rgb(228, 237, 237);
  background-size: cover;
  background-position: center;
  width: 8000px;
  //height: 600px;
  //margin-top: 60px;
  border-radius: 50px;
}

.container {
  left: 3%;
}
</style>
