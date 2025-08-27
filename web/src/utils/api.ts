import { showNotification } from "@mantine/notifications";
import axios, { AxiosError, AxiosResponse } from "axios";

// Create axios instance
const api = axios.create({
  baseURL: `${window.location.origin}/api`,
  timeout: 10000,
});

// Request interceptor
api.interceptors.request.use(
  (config: any) => {
    return config;
  },
  (error: AxiosError) => {
    showNotification({
      title: "Request Error",
      message: error.message,
      color: "red",
    });
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  (error: AxiosError) => {
    showNotification({
      title: "Error",
      message: (error.response?.data as any)?.message || error.message,
      color: "red",
    });
    return Promise.reject(error);
  }
);

export default api;
