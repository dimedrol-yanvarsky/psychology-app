import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    user: null,
    isSavingProfile: false,
    isDeletingAccount: false,

    adminAccounts: [],
    isAdminListLoading: false,
    adminListError: "",
    blockingUsers: {},
    deletingUsers: {},

    isTestModalOpen: false,
    isTerminalOpen: false,
    answersModal: {
        open: false,
        loading: false,
        error: "",
        answers: [],
        questions: [],
        title: "",
    },

    completedTests: [],
    isCompletedTestsLoading: true,
    completedTestsError: "",
};

// Слайс дашборда: профиль, админ-панель, модалки, пройденные тесты.
const dashboardSlice = createSlice({
    name: "dashboard",
    initialState,
    reducers: {
        // Профиль
        profileSaveStart(state) {
            state.isSavingProfile = true;
        },
        profileSaveSuccess(state, action) {
            state.isSavingProfile = false;
            state.user = action.payload.user;
        },
        profileSaveError(state) {
            state.isSavingProfile = false;
        },
        setUser(state, action) {
            state.user = action.payload;
        },

        // Удаление аккаунта
        deleteAccountStart(state) {
            state.isDeletingAccount = true;
        },
        deleteAccountSuccess(state) {
            state.isDeletingAccount = false;
            state.user = null;
        },
        deleteAccountError(state) {
            state.isDeletingAccount = false;
        },

        // Админ-панель
        adminListLoading(state) {
            state.isAdminListLoading = true;
            state.adminListError = "";
        },
        adminListSuccess(state, action) {
            state.isAdminListLoading = false;
            state.adminAccounts = action.payload;
        },
        adminListError(state, action) {
            state.isAdminListLoading = false;
            state.adminAccounts = [];
            state.adminListError = action.payload;
        },
        adminListReset(state) {
            state.isAdminListLoading = false;
            state.adminAccounts = [];
            state.adminListError = "";
        },
        blockUserStart(state, action) {
            state.blockingUsers[action.payload] = true;
        },
        blockUserEnd(state, action) {
            delete state.blockingUsers[action.payload];
        },
        updateAccountStatus(state, action) {
            const { id, status } = action.payload;
            const account = state.adminAccounts.find((u) => u.id === id);
            if (account) {
                account.status = status;
            }
        },
        deleteUserStart(state, action) {
            state.deletingUsers[action.payload] = true;
        },
        deleteUserEnd(state, action) {
            delete state.deletingUsers[action.payload];
        },

        // Модальные окна
        openTestModal(state) {
            state.isTestModalOpen = true;
        },
        closeTestModal(state) {
            state.isTestModalOpen = false;
        },
        toggleTerminal(state) {
            state.isTerminalOpen = !state.isTerminalOpen;
        },
        setTerminalOpen(state, action) {
            state.isTerminalOpen = action.payload;
        },
        closeTerminal(state) {
            state.isTerminalOpen = false;
        },

        // Модалка ответов
        openAnswersModal(state, action) {
            state.answersModal = {
                open: true,
                loading: true,
                error: "",
                answers: [],
                questions: [],
                title: action.payload,
            };
        },
        answersModalSuccess(state, action) {
            state.answersModal.loading = false;
            state.answersModal.answers = action.payload.answers;
            state.answersModal.questions = action.payload.questions;
        },
        answersModalError(state, action) {
            state.answersModal.loading = false;
            state.answersModal.error = action.payload;
        },
        closeAnswersModal(state) {
            state.answersModal = {
                open: false,
                loading: false,
                error: "",
                answers: [],
                questions: [],
                title: "",
            };
        },

        // Пройденные тесты
        completedTestsLoading(state) {
            state.isCompletedTestsLoading = true;
            state.completedTestsError = "";
        },
        completedTestsSuccess(state, action) {
            state.isCompletedTestsLoading = false;
            state.completedTests = action.payload;
        },
        completedTestsError(state, action) {
            state.isCompletedTestsLoading = false;
            state.completedTests = [];
            state.completedTestsError = action.payload;
        },
        completedTestsReset(state) {
            state.isCompletedTestsLoading = false;
            state.completedTests = [];
        },
    },
});

export const {
    profileSaveStart,
    profileSaveSuccess,
    profileSaveError,
    setUser,
    deleteAccountStart,
    deleteAccountSuccess,
    deleteAccountError,
    adminListLoading,
    adminListSuccess,
    adminListError,
    adminListReset,
    blockUserStart,
    blockUserEnd,
    updateAccountStatus,
    deleteUserStart,
    deleteUserEnd,
    openTestModal,
    closeTestModal,
    toggleTerminal,
    setTerminalOpen,
    closeTerminal,
    openAnswersModal,
    answersModalSuccess,
    answersModalError,
    closeAnswersModal,
    completedTestsLoading,
    completedTestsSuccess,
    completedTestsError,
    completedTestsReset,
} = dashboardSlice.actions;

export default dashboardSlice.reducer;
