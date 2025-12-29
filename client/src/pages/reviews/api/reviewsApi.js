import axios from "axios";
import { API_BASE_URL } from "../../../shared/config/api";

const api = axios.create({
    baseURL: `${API_BASE_URL}/api`,
});

export const getReviews = () => api.get("/reviews/getReviews");

export const createReview = (payload) =>
    api.post("/reviews/createReview", payload);

export const deleteReview = (payload) =>
    api.post("/reviews/deleteReview", payload);

export const decideReview = (payload) =>
    api.post("/reviews/approveOrDeny", payload);

export const updateReview = (payload) =>
    api.post("/reviews/updateReview", payload);
