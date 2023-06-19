<template>
  <router-view></router-view>
  <CommHeader />
  <div class="common-layout">
    <el-card class="box-card">
      <el-form-item>
        <h1 class="Tips">Step 2 of 2: test your narrow vocab level</h1>
        <el-divider />
      </el-form-item>
      <el-container class="container">
        <el-header>
          Use this larger list to calculate your vocab size with greater
          precision. Check the box if you know at least one definition for a
          word.
        </el-header>
      </el-container>
      <div>
        <div
          class="checkbox-row"
          v-for="(row, index) in checkboxRows"
          :key="index"
        >
          <el-checkbox
            class="checkbox-item"
            v-for="item in row"
            :key="item.id"
            v-model="item.checked"
          >
            <span class="checkbox-label">{{ item.word }}</span>
          </el-checkbox>
        </div>
      </div>
      <el-button type="primary" @click="sendData">Commit</el-button>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import CommHeader from "@/components/common/CommHeader.vue";
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import axios from "axios";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";
interface CheckboxOption {
  id: number;
  word: string;
  checked: boolean;
}

const { StringWordList, StringWordList1 } = defineProps({
  StringWordList: {
    type: String,
    required: true,
  },
  StringWordList1: {
    type: String,
    required: true,
  },
});

const parsedResult = computed(() => {
  const result = JSON.parse(StringWordList || "[]") as {
    id: number;
    word: string;
  }[];
  return Array.isArray(result) ? result : [];
});

const parsedResult1 = computed(() => {
  const result = JSON.parse(StringWordList1 || "[]") as {
    id: number;
    word: string;
    known: number;
  }[];
  return Array.isArray(result) ? result : [];
});

const route = useRoute();
const router = useRouter();
const checkboxOptions = ref<CheckboxOption[]>([]);

const initializeOptions = () => {
  checkboxOptions.value = parsedResult.value.map((item) => ({
    id: item.id,
    word: item.word,
    checked: false,
  }));
};

onMounted(() => {
  initializeOptions();
});

const checkboxRows = computed(() => {
  const options = checkboxOptions.value;
  const rows: CheckboxOption[][] = [];
  for (let i = 0; i < options.length; i += 4) {
    rows.push(options.slice(i, i + 4));
  }
  return rows;
});
const sendData = async () => {
  const selectedOptions = checkboxOptions.value.map((item) => ({
    id: item.id,
    word: item.word,
    known: item.checked ? 1 : 0,
  }));
  const userStore = useUserStore();
  const token = userStore.token;
  // result = selectedOptions + StringWordList1;
  parsedResult1.value.forEach((item) => {
    selectedOptions.push(item);
  });
  const requestData = { WordList: selectedOptions };
  console.log(requestData);
  try {
    const response = await axios.post(
      "/api/basic-api/word/getVocabulary",
      requestData,
      {
        headers: { Authorization: `${token}` },
      }
    );
    router.push({
      name: "resulthash",
      query: {
        data: response.data.result.toString(),
      },
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
  // background-image: url('@/assets/img/动态图.jpg');
  background-size: cover;
  background-position: center;
  width: 8000px;
  // height: 600px;
  // margin-top: 60px;
  border-radius: 50px;
}

// .continue-test {
//   left: 5%;
//   top: 7%;
//   background: linear-gradient(to bottom, #8ab7ee, #2582ee);
//   border: none;
//   color: white;
//   padding: 15px 32px;
//   text-align: center;
//   text-decoration: none;
//   display: inline-block;
//   font-size: 16px;
//   border-radius: 25px;
//   box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.25);
//   transition: all 0.3s ease-in-out;
// }

// .list-1 {
//   left: 20%;
//   top: 7%;
//   background: linear-gradient(to bottom, #8ab7ee, #2582ee);
//   border: none;
//   color: white;
//   padding: 15px 32px;
//   text-align: center;
//   text-decoration: none;
//   display: inline-block;
//   font-size: 16px;
//   border-radius: 25px;
//   box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.25);
//   transition: all 0.3s ease-in-out;
// }

.container {
  left: 3%;
}
</style>
