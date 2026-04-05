import axios from "axios";

const api = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
  withCredentials: true, // 自动发送和接收 Cookie
});

api.interceptors.response.use(
  (response) => {
    const res = response.data;
    if (res.code === 0) {
      return res.data;
    } else {
      console.error(`响应错误: ${response.config.url}`, res.code, res.message);
      return Promise.reject(new Error(res.message));
    }
  },
  (error) => {
    if (error.response) {
      console.error(`响应错误: ${error.config?.url}`, error.response.status, error.response.data);
      if (error.response.status === 401) {
        localStorage.removeItem("admin_user");
        window.location.href = "/console/login";
      }
    } else if (error.request) {
      console.error("请求错误（无响应）:", error.config?.url, error.message);
    } else {
      console.error("请求错误:", error.message);
    }
    return Promise.reject(error);
  },
);

export const authApi = {
  logout: () => api.post('/logout'),
}

export default api;
