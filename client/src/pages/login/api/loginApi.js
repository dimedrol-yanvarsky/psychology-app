import axios from "axios";
import { API_BASE_URL } from "../../../shared/config/api";

const api = axios.create({
    baseURL: `${API_BASE_URL}/api/login`,
});

export const loginWithPassword = (payload) => api.post("/password", payload);

export const loginWithProvider = (provider) => api.post(`/${provider}`);
