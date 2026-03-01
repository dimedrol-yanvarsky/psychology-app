// Начальное состояние дашборда, объединяющее профиль, админ-панель, модалки и тесты.
export const getInitialDashboardState = (isAdminStatus) => ({
    // Профиль
    user: null,
    isSavingProfile: false,
    isDeletingAccount: false,

    // Админ-панель
    adminAccounts: [],
    isAdminListLoading: isAdminStatus,
    adminListError: "",
    blockingUsers: {},
    deletingUsers: {},

    // Модальные окна
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

    // Пройденные тесты
    completedTests: [],
    isCompletedTestsLoading: true,
    completedTestsError: "",
});

// Редьюсер дашборда, обрабатывающий все переходы состояний.
export const dashboardReducer = (state, action) => {
    switch (action.type) {
        // --- Профиль ---
        case "PROFILE_SAVE_START":
            return { ...state, isSavingProfile: true };
        case "PROFILE_SAVE_SUCCESS":
            return {
                ...state,
                isSavingProfile: false,
                user: action.payload.user,
            };
        case "PROFILE_SAVE_ERROR":
            return { ...state, isSavingProfile: false };
        case "SET_USER":
            return { ...state, user: action.payload };

        // --- Удаление аккаунта ---
        case "DELETE_ACCOUNT_START":
            return { ...state, isDeletingAccount: true };
        case "DELETE_ACCOUNT_SUCCESS":
            return { ...state, isDeletingAccount: false, user: null };
        case "DELETE_ACCOUNT_ERROR":
            return { ...state, isDeletingAccount: false };

        // --- Админ-панель ---
        case "ADMIN_LIST_LOADING":
            return { ...state, isAdminListLoading: true, adminListError: "" };
        case "ADMIN_LIST_SUCCESS":
            return {
                ...state,
                isAdminListLoading: false,
                adminAccounts: action.payload,
            };
        case "ADMIN_LIST_ERROR":
            return {
                ...state,
                isAdminListLoading: false,
                adminAccounts: [],
                adminListError: action.payload,
            };
        case "ADMIN_LIST_RESET":
            return {
                ...state,
                isAdminListLoading: false,
                adminAccounts: [],
                adminListError: "",
            };
        case "BLOCK_USER_START":
            return {
                ...state,
                blockingUsers: { ...state.blockingUsers, [action.payload]: true },
            };
        case "BLOCK_USER_END": {
            const updated = { ...state.blockingUsers };
            delete updated[action.payload];
            return { ...state, blockingUsers: updated };
        }
        case "UPDATE_ACCOUNT_STATUS":
            return {
                ...state,
                adminAccounts: state.adminAccounts.map((user) =>
                    user.id === action.payload.id
                        ? { ...user, status: action.payload.status }
                        : user
                ),
            };
        case "DELETE_USER_START":
            return {
                ...state,
                deletingUsers: { ...state.deletingUsers, [action.payload]: true },
            };
        case "DELETE_USER_END": {
            const updated = { ...state.deletingUsers };
            delete updated[action.payload];
            return { ...state, deletingUsers: updated };
        }

        // --- Модальные окна ---
        case "OPEN_TEST_MODAL":
            return { ...state, isTestModalOpen: true };
        case "CLOSE_TEST_MODAL":
            return { ...state, isTestModalOpen: false };
        case "TOGGLE_TERMINAL":
            return { ...state, isTerminalOpen: !state.isTerminalOpen };
        case "SET_TERMINAL_OPEN":
            return { ...state, isTerminalOpen: action.payload };
        case "CLOSE_TERMINAL":
            return { ...state, isTerminalOpen: false };

        // --- Модалка ответов ---
        case "OPEN_ANSWERS_MODAL":
            return {
                ...state,
                answersModal: {
                    open: true,
                    loading: true,
                    error: "",
                    answers: [],
                    questions: [],
                    title: action.payload,
                },
            };
        case "ANSWERS_MODAL_SUCCESS":
            return {
                ...state,
                answersModal: {
                    ...state.answersModal,
                    loading: false,
                    answers: action.payload.answers,
                    questions: action.payload.questions,
                },
            };
        case "ANSWERS_MODAL_ERROR":
            return {
                ...state,
                answersModal: {
                    ...state.answersModal,
                    loading: false,
                    error: action.payload,
                },
            };
        case "CLOSE_ANSWERS_MODAL":
            return {
                ...state,
                answersModal: {
                    open: false,
                    loading: false,
                    error: "",
                    answers: [],
                    questions: [],
                    title: "",
                },
            };

        // --- Пройденные тесты ---
        case "COMPLETED_TESTS_LOADING":
            return {
                ...state,
                isCompletedTestsLoading: true,
                completedTestsError: "",
            };
        case "COMPLETED_TESTS_SUCCESS":
            return {
                ...state,
                isCompletedTestsLoading: false,
                completedTests: action.payload,
            };
        case "COMPLETED_TESTS_ERROR":
            return {
                ...state,
                isCompletedTestsLoading: false,
                completedTests: [],
                completedTestsError: action.payload,
            };
        case "COMPLETED_TESTS_RESET":
            return {
                ...state,
                isCompletedTestsLoading: false,
                completedTests: [],
            };

        default:
            return state;
    }
};
