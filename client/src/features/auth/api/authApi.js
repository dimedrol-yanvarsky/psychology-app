import axios from "axios";
import { API_BASE_URL } from "../../../shared/config/api";

const api = axios.create({
    baseURL: `${API_BASE_URL}/api/login`,
});

const oauthApi = axios.create({
    baseURL: `${API_BASE_URL}/api/auth`,
});

export const loginWithPassword = (payload) => api.post("/password", payload);

export const loginWithProvider = (provider) => oauthApi.post(`/${provider}`);

export const createAccount = (payload) =>
    axios.post(`${API_BASE_URL}/api/createAccount`, payload);
