import axios from "axios";

const api = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("admin_token");
    console.log(`请求拦截器: token=${token}, url=${config.url}`);
    if (token && token !== "null" && token !== "undefined") {
      config.headers.Authorization = `Bearer ${token}`;
    }
    console.log(`请求: ${config.method?.toUpperCase()} ${config.baseURL}${config.url}`, config.params || {}, config.data || {});
    return config;
  },
  (error) => {
    console.error("请求拦截器错误:", error);
    return Promise.reject(error);
  },
);

api.interceptors.response.use(
  (response) => {
    const res = response.data;
    if (res.code === 0) {
      return res.data;
    } else {
      console.error(`响应错误: ${response.config.url}`, res.code, res.message);
      alert(res.message || "请求失败");
      return Promise.reject(new Error(res.message));
    }
  },
  (error) => {
    if (error.response) {
      console.error(`响应错误: ${error.config?.url}`, error.response.status, error.response.data);
      if (error.response.status === 401) {
        localStorage.removeItem("admin_token");
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

export default api;
