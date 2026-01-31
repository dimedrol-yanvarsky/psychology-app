import axios from "axios";
import { API_BASE_URL } from "../../../shared/config/api";

const api = axios.create({
    baseURL: `${API_BASE_URL}/api/dashboard`,
});

export const changeUserData = (payload) =>
    api.post("/change-user-data", payload);

export const deleteAccount = (payload) => api.post("/delete-account", payload);

export const fetchUsers = (payload) => api.post("/users", payload);

export const blockUser = (payload) => api.post("/block-user", payload);

export const deleteUser = (payload) => api.post("/delete-user", payload);

export const fetchCompletedTests = (payload) =>
    api.post("/completed-tests", payload);

export const fetchUserAnswers = (payload) =>
    api.post("/user-answers", payload);
