import axios from "axios";
import { ElMessage } from 'element-plus';

const api = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

api.interceptors.response.use(
  (response) => {
    const res = response.data;
    if (res.code === 0) {
      return res.data;
    } else {
      ElMessage.error(res.message || "请求失败")
      return Promise.reject(new Error(res.message));
    }
  },
  (error) => {
    if (error.response) {
      // 服务器返回了错误响应
      const data = error.response.data;
      const message = data?.message || data?.msg || "请求失败";
      ElMessage.error(message)
    } else if (error.message) {
      // 网络错误等
      ElMessage.error(error.message)
    }
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);

export default api;
