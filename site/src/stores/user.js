import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api, { authApi } from "@/api";

export const useUserStore = defineStore("user", () => {
  const user = ref(JSON.parse(localStorage.getItem("user") || "null"));

  const isLoggedIn = computed(() => !!user.value);

  async function login(loginData) {
    const res = await api.post("/login", loginData);
    user.value = res.user;
    localStorage.setItem("user", JSON.stringify(res.user));
  }

  async function register(registerData) {
    const res = await api.post("/register", registerData);
    // 注册成功后刷新用户信息，确保状态同步
    await refreshProfile();
  }

  async function logout() {
    try {
      await authApi.logout();
    } catch (e) {
      // 即使 API 调用失败，也清除本地状态
    }
    user.value = null;
    localStorage.removeItem("user");
  }

  async function refreshProfile() {
    const res = await api.get("/user/profile");
    user.value = res;
    localStorage.setItem("user", JSON.stringify(res));
  }

  return {
    user,
    isLoggedIn,
    login,
    register,
    logout,
    refreshProfile,
  };
});
