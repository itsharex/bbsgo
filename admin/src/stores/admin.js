import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/api";

export const useAdminStore = defineStore("admin", () => {
  const token = ref(localStorage.getItem("admin_token") || "");
  const user = ref(JSON.parse(localStorage.getItem("admin_user") || "null"));

  const isLoggedIn = computed(() => !!token.value && !!user.value);
  const isAdmin = computed(() => user.value?.role >= 2);

  async function login(loginData) {
    const res = await api.post("/login", loginData);
    if (res.user.role < 2) {
      throw new Error("需要管理员权限");
    }
    token.value = res.token;
    user.value = res.user;
    localStorage.setItem("admin_token", res.token);
    localStorage.setItem("admin_user", JSON.stringify(res.user));
  }

  function logout() {
    token.value = "";
    user.value = null;
    localStorage.removeItem("admin_token");
    localStorage.removeItem("admin_user");
  }

  return {
    token,
    user,
    isLoggedIn,
    isAdmin,
    login,
    logout,
  };
});
