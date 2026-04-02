import axios from "axios";

const api = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("admin_token");
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
      alert(res.message || "请求失败");
      return Promise.reject(new Error(res.message));
    }
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("admin_token");
      localStorage.removeItem("admin_user");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);

export default api;
