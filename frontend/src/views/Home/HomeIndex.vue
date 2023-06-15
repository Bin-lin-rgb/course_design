<template>
  <CommHeader />
  <div class="common-layout">
    <el-form-item>
      <h1 style="left: 5%;">Step 1 of 2: test your broad vocab level</h1>
      <el-divider />
    </el-form-item>
    <el-container class="container">
      <el-header>Check the box if you know at least one definition for a word. If you’re not sure about the exact meaning,
        leave it blank.</el-header>
    </el-container>
    <div>
      <div class="checkbox-row" v-for="(row, index) in checkboxRows" :key="index">
        <el-checkbox class="checkbox-item" v-for="item in row" :key="item.id" v-model="item.checked"
          @change="updateValue(item)">
          <span class="checkbox-label"> {{ item.label }}</span>
        </el-checkbox>
      </div>
    </div>
    <button class="continue-test" @click="sendData">继续测试</button>
    <button class="list-1" @click="aaa">List1</button>
  </div>
</template>

<script lang="ts" setup>
import CommHeader from "@/components/common/CommHeader.vue";
import { ref, computed, onMounted } from 'vue'
import axios from "axios";
import { GetWordList1 } from '@/utils/api/wordBook'
// import { GetPostListWithTime } from "@/utils/api/article";
// import { ref, onBeforeMount } from "vue";
// import { ElMessage } from "element-plus";
// import { useRouter } from "vue-router";
// const router = useRouter();
interface CheckboxOption {
  id: number;
  label: string;
  checked: boolean;
  value: number;
}

interface SelectedOption {
  id: number;
  label: string;
}

const checkboxOptions = ref<CheckboxOption[] >([
  { id: 1, label: '选项1', checked: false, value: 0 },
  { id: 2, label: '你好2', checked: false, value: 0 },
  { id: 3, label: '选项3', checked: false, value: 0 },
  { id: 4, label: '选项4', checked: false, value: 0 },
  { id: 5, label: '选项5', checked: false, value: 0 },
  { id: 6, label: '选项6', checked: false, value: 0 },
  { id: 7, label: '选项7', checked: false, value: 0 },
  { id: 8, label: '选项8', checked: false, value: 0 },
])

const checkboxRows = computed(() => {
  const rows : CheckboxOption[][] = [];
  const options: CheckboxOption[] = [...checkboxOptions.value]; // 复制一份选项数组，以免修改原始数据
  while (options.length) {
    rows.push(options.splice(0, 4));
  }
  return rows;
})
const selectedOptions = ref<SelectedOption[]>([]);

const updateValue =  (item: { id: number, label: string, checked: boolean, value: number }) => {
  item.value = item.checked ? 1 : 0
  //console.log(item.value)
}

async function sendData() {
    const selected: CheckboxOption[] = checkboxOptions.value.filter((item) => item.value === 1);
    selectedOptions.value = selected.map((item) => ({ id: item.id, label: item.label }));

    console.log(selectedOptions.value); // 打印数组内容

    try {
      const response = await axios.post('/api/sendData', selectedOptions.value);
      console.log(response.data);
    } catch (error) {
      console.error(error);
    }
  }


  const aaa = async () => {
    const res: any = await GetWordList1()
    console.log(res.code)
    console.log(res.result)
  }
    
  
</script>
<style scoped lang="less">
@normal_text_color: #86909c;

// .common-layout {
//   margin-top: 120px !important;
// }
.checkbox-row {
  display: flex;
  left: 5%;
}

.checkbox-item {
  flex: 1;
  margin-right: 10px;
  zoom: 150%;
  align-items: center
}

.checkbox-label {

  font-size: 12px;
}

.common-layout{
  background-color: rgb(228, 237, 237);
  // background-image: url('@/assets/img/动态图.jpg');
  background-size: cover;
  background-position: center;
  width: 8000px;
  height: 600px;
  margin-top: 60px; 
  border-radius: 50px;
}

.continue-test{
  left: 5%;
  top:7%;
  background: linear-gradient(to bottom, #8ab7ee, #2582ee);
  border: none;
  color: white;
  padding: 15px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  border-radius: 25px;
  box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.25);
  transition: all 0.3s ease-in-out;
}

.list-1{
  left: 20%;
  top:7%;
  background: linear-gradient(to bottom, #8ab7ee, #2582ee);
  border: none;
  color: white;
  padding: 15px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  border-radius: 25px;
  box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.25);
  transition: all 0.3s ease-in-out;
}
.container{
  left: 3%;
}

</style>
