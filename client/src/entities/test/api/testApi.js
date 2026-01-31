import axios from "axios";
import { API_BASE_URL } from "../../../shared/config/api";

const api = axios.create({
    baseURL: `${API_BASE_URL}/api/tests`,
});

export const fetchTests = (payload) => api.post("/getTests", payload);

export const deleteTest = (payload) => api.post("/deleteTest", payload);

export const fetchQuestions = (payload) => api.post("/getQuestions", payload);

export const submitAttempt = (payload) => api.post("/attemptTest", payload);

export const changeTest = (payload) => api.post("/changeTest", payload);

export const addTest = (payload) => api.post("/addTest", payload);
