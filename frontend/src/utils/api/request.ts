import axios from "axios";
import { useUserStore } from "@/stores/user";

//1. 创建axios对象
const service = axios.create();

//2. 请求拦截器
service.interceptors.request.use(
  (config) => {
    const userStore = useUserStore();
    const token = userStore.token;
    if (token) {
      if (config && config.headers) {
        // 多一步判断
        config.headers.Authorization = `Bearer ${token}`; // 将accessToken放入到请求头里
      }
    }
    return config;
  },
  (error) => {
    Promise.reject(error);
  }
);

//3. 响应拦截器
service.interceptors.response.use(
  (response) => {
    if (response.status === 200) {
      return Promise.resolve(response.data);
    } else {
      return Promise.reject(response.data);
    }
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default service;
