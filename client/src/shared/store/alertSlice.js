import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    statusAlert: "",
    messageAlert: "",
};

// Слайс уведомлений: управляет показом алертов.
const alertSlice = createSlice({
    name: "alert",
    initialState,
    reducers: {
        setAlert(state, action) {
            state.statusAlert = action.payload.status;
            state.messageAlert = action.payload.message;
        },
        clearAlert(state) {
            state.statusAlert = "";
            state.messageAlert = "";
        },
    },
});

export const { setAlert, clearAlert } = alertSlice.actions;

// Thunk для показа уведомления с автоочисткой через 3 секунды.
let alertTimer = null;
export const showAlert = (status, message) => (dispatch) => {
    clearTimeout(alertTimer);
    dispatch(setAlert({ status, message }));
    alertTimer = setTimeout(() => {
        dispatch(clearAlert());
    }, 3000);
    return true;
};

export default alertSlice.reducer;
