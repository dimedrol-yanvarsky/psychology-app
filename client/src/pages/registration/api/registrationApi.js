import axios from "axios";
import { API_BASE_URL } from "../../../shared/config/api";

export const createAccount = (payload) =>
    axios.post(`${API_BASE_URL}/api/createAccount`, payload);
