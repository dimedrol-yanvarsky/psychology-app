import { createSlice } from "@reduxjs/toolkit";
import { createDefaultProfileData } from "../../entities/user";

const initialState = {
    isAuth: false,
    isAdmin: false,
    profileData: createDefaultProfileData(),
};

// Слайс авторизации: управляет состоянием сессии и профилем пользователя.
const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        setIsAuth(state, action) {
            state.isAuth = action.payload;
        },
        setIsAdmin(state, action) {
            state.isAdmin = action.payload;
        },
        setProfileData(state, action) {
            state.profileData = action.payload;
        },
        updateProfileData(state, action) {
            Object.assign(state.profileData, action.payload);
        },
        resetAuth(state) {
            state.isAuth = false;
            state.isAdmin = false;
            state.profileData = createDefaultProfileData();
        },
    },
});

export const {
    setIsAuth,
    setIsAdmin,
    setProfileData,
    updateProfileData,
    resetAuth,
} = authSlice.actions;

export default authSlice.reducer;
