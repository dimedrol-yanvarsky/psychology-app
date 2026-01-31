import axios from "axios";
import { API_BASE_URL } from "../../../shared/config/api";

const api = axios.create({
    baseURL: `${API_BASE_URL}/api/recommendations`,
});

export const fetchRecommendations = () => api.get("/list");

export const updateBlock = (payload) => api.post("/updateBlock", payload);

export const addSection = () => api.post("/addSection");

export const addBlock = (payload) => api.post("/addBlock", payload);

export const deleteBlock = (payload) => api.post("/deleteBlock", payload);

export const deleteSection = (payload) => api.post("/deleteSection", payload);
