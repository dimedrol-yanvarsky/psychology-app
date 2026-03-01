import { useEffect, useMemo } from "react";
import { useSelector, useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { useLockBodyScroll } from "../../../shared/lib/hooks/useLockBodyScroll";
import { createDefaultProfileData } from "../../../entities/user";
import {
    blockUser,
    changeUserData,
    deleteAccount,
    deleteUser,
    fetchCompletedTests,
    fetchUserAnswers,
    fetchUsers,
} from "../../../entities/session";
import { useAuthContext } from "../../../shared/context/AuthContext";
import { useAlertContext } from "../../../shared/context/AlertContext";
import {
    profileSaveStart,
    profileSaveSuccess,
    profileSaveError,
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
} from "./dashboardSlice";

const DEFAULT_EMOTION_DATA = [
    { id: "calm", label: "Спокойствие", value: 72 },
    { id: "energy", label: "Энергия", value: 54 },
    { id: "focus", label: "Фокус", value: 61 },
    { id: "stress", label: "Стресс", value: 28 },
];

export const useDashboard = () => {
    const { isAdmin, profileData, setProfileData, setIsAuth, setIsAdmin } =
        useAuthContext();
    const { showAlert } = useAlertContext();
    const navigate = useNavigate();
    const reduxDispatch = useDispatch();

    const state = useSelector((s) => s.dashboard);

    const selectedAnswers = useMemo(() => {
        const map = new Map();
        if (!Array.isArray(state.answersModal.answers)) {
            return map;
        }

        state.answersModal.answers.forEach((item) => {
            if (Array.isArray(item) && item.length > 1) {
                const [qNumber, ...rest] = item;
                const normalized = Number(qNumber);
                const answersSet = new Set(
                    rest.map((num) => Number(num)).filter((num) => !isNaN(num))
                );
                if (!isNaN(normalized)) {
                    map.set(normalized, answersSet);
                }
            }
        });

        return map;
    }, [state.answersModal.answers]);

    const hasGoogle = Boolean(profileData?.isGoogleAdded);
    const hasYandex = Boolean(profileData?.isYandexAdded);
    const showLinkButtons = !hasGoogle || !hasYandex;

    const handleFieldChange = (field, value) => {
        setProfileData((prev) => ({ ...prev, [field]: value }));
    };

    const handleProfileSave = async (event) => {
        event.preventDefault();

        if (state.isSavingProfile) {
            return;
        }

        const userId = profileData?.id;
        const firstName = (profileData?.firstName || "").trim();

        if (!userId) {
            showAlert("error", "Не указан идентификатор пользователя для обновления");
            return;
        }

        if (!firstName) {
            showAlert("error", "Введите имя, чтобы сохранить изменения");
            return;
        }

        reduxDispatch(profileSaveStart());

        try {
            const { data } = await changeUserData({ userId, firstName });

            if (data?.status === "success") {
                const updatedUser = data?.user || {};

                setProfileData((prev) => ({
                    ...prev,
                    id: updatedUser.id || userId,
                    firstName: updatedUser.firstName || firstName,
                    email: updatedUser.email || prev.email,
                    status: updatedUser.status || prev.status,
                }));

                reduxDispatch(
                    profileSaveSuccess({
                        user: {
                            id: updatedUser.id || userId,
                            firstName: updatedUser.firstName || firstName,
                            email: updatedUser.email || profileData.email,
                            status: updatedUser.status || profileData.status,
                        },
                    })
                );

                showAlert("success", data?.message || "Данные профиля обновлены");
            } else {
                reduxDispatch(profileSaveError());
                showAlert("error", data?.message || "Не удалось сохранить изменения");
            }
        } catch (error) {
            reduxDispatch(profileSaveError());
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось сохранить изменения";
            showAlert("error", message);
        }
    };

    const handleLogout = () => {
        setIsAuth(false);
        setIsAdmin(false);
        setProfileData(createDefaultProfileData());
        showAlert("success", "Вы вышли из аккаунта");
        navigate("/login");
    };

    const handleDeleteAccount = async () => {
        if (state.isDeletingAccount) {
            return;
        }

        if (!profileData?.id) {
            showAlert("error", "Не указан идентификатор пользователя для удаления");
            return;
        }

        reduxDispatch(deleteAccountStart());

        try {
            const { data } = await deleteAccount({ userId: profileData.id });

            if (data?.status === "success") {
                reduxDispatch(deleteAccountSuccess());
                setIsAdmin(false);
                setIsAuth(false);
                setProfileData(createDefaultProfileData());

                showAlert("success", data?.message || "Аккаунт удален");
                navigate(data?.redirect || "/login", { replace: true });
            } else {
                reduxDispatch(deleteAccountError());
                showAlert("error", data?.message || "Не удалось удалить аккаунт");
            }
        } catch (error) {
            reduxDispatch(deleteAccountError());
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось удалить аккаунт";
            showAlert("error", message);
        }
    };

    const handleChangePassword = () => {
        showAlert("success", "Мы отправим ссылку для смены пароля на вашу почту");
    };

    const handleLinkProvider = async (provider) => {
        if (!profileData?.id) {
            showAlert("error", "Пользователь не авторизован");
            return;
        }

        try {
            const response = await fetch(
                `http://localhost:8080/api/auth/link/${provider}`,
                {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ userId: profileData.id }),
                }
            );

            const data = await response.json();

            if (data.redirectUrl) {
                window.location.href = data.redirectUrl;
            } else if (data.error) {
                showAlert("error", data.error);
            } else {
                showAlert("error", "Не получен URL для авторизации");
            }
        } catch (error) {
            showAlert("error", `Ошибка при привязке аккаунта ${provider}`);
        }
    };

    const handleStartTesting = () => {
        reduxDispatch(closeTestModal());
        return showAlert("error", "Тестирование пока недоступно");
    };

    const handleAdminAction = (action, account) => {
        return showAlert("error", `${action} для ${account.email} недоступно`);
    };

    const handleBlockUser = async (account) => {
        if (!profileData?.id) {
            showAlert("error", "Не указан идентификатор администратора для запроса");
            return;
        }

        const targetId = account?.id;
        const status = (account?.status || "").trim();

        if (!targetId) {
            showAlert("error", "Не удалось определить пользователя для блокировки");
            return;
        }

        if (status === "Заблокирован" || status === "Удален") {
            return;
        }

        reduxDispatch(blockUserStart(targetId));

        try {
            const { data } = await blockUser({
                adminId: profileData.id,
                targetUserId: targetId,
            });

            if (data?.status === "success") {
                reduxDispatch(
                    updateAccountStatus({
                        id: targetId,
                        status: data?.user?.status || "Заблокирован",
                    })
                );
                showAlert("success", data?.message || "Пользователь заблокирован");
            } else {
                showAlert(
                    "error",
                    data?.message || "Не удалось заблокировать пользователя"
                );
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось заблокировать пользователя";
            showAlert("error", message);
        } finally {
            reduxDispatch(blockUserEnd(targetId));
        }
    };

    const handleDeleteUser = async (account) => {
        if (!profileData?.id) {
            showAlert("error", "Не указан идентификатор администратора для запроса");
            return;
        }

        const targetId = account?.id;
        const status = (account?.status || "").trim();

        if (!targetId) {
            showAlert("error", "Не удалось определить пользователя для удаления");
            return;
        }

        if (status === "Удален") {
            return;
        }

        reduxDispatch(deleteUserStart(targetId));

        try {
            const { data } = await deleteUser({
                adminId: profileData.id,
                targetUserId: targetId,
            });

            if (data?.status === "success") {
                reduxDispatch(
                    updateAccountStatus({
                        id: targetId,
                        status: data?.user?.status || "Удален",
                    })
                );
                showAlert("success", data?.message || "Пользователь удален");
            } else {
                showAlert(
                    "error",
                    data?.message || "Не удалось удалить пользователя"
                );
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось удалить пользователя";
            showAlert("error", message);
        } finally {
            reduxDispatch(deleteUserEnd(targetId));
        }
    };

    // Загрузка списка пользователей для админ-панели.
    useEffect(() => {
        if (profileData?.status !== "Администратор") {
            reduxDispatch(adminListReset());
            return;
        }

        if (!profileData?.id) {
            reduxDispatch(
                adminListError("Не указан идентификатор пользователя")
            );
            return;
        }

        const fetchAdminUsers = async () => {
            reduxDispatch(adminListLoading());

            try {
                const { data } = await fetchUsers({
                    userId: profileData.id,
                    status: profileData.status,
                });

                if (data?.status === "success") {
                    reduxDispatch(
                        adminListSuccess(
                            Array.isArray(data.users) ? data.users : []
                        )
                    );
                } else {
                    reduxDispatch(
                        adminListError(
                            data?.message ||
                                "Не удалось загрузить пользователей"
                        )
                    );
                }
            } catch (error) {
                const message =
                    error?.response?.data?.message ||
                    error?.message ||
                    "Ошибка загрузки пользователей";
                reduxDispatch(adminListError(message));
            }
        };

        fetchAdminUsers();
    }, [profileData?.id, profileData?.status, reduxDispatch]);

    const handleOpenAnswersModal = async (test) => {
        reduxDispatch(
            openAnswersModal(test.testName || "Пройденный тест")
        );

        if (!profileData?.id || !test?.id || !test?.testId) {
            reduxDispatch(
                answersModalError(
                    "Недостаточно данных для загрузки результата"
                )
            );
            return;
        }

        try {
            const { data } = await fetchUserAnswers({
                userId: profileData.id,
                completedTestId: test.id,
                testId: test.testId,
            });

            if (data?.status === "success") {
                reduxDispatch(
                    answersModalSuccess({
                        answers: Array.isArray(data.answers)
                            ? data.answers
                            : [],
                        questions: Array.isArray(data.questions)
                            ? data.questions
                            : [],
                    })
                );
            } else {
                reduxDispatch(
                    answersModalError(
                        data?.message ||
                            "Не удалось загрузить ответы пользователя"
                    )
                );
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Ошибка загрузки ответов пользователя";
            reduxDispatch(answersModalError(message));
        }
    };

    const handleCloseAnswersModal = () => {
        reduxDispatch(closeAnswersModal());
    };

    useLockBodyScroll(state.answersModal.open);

    // Загрузка пройденных тестов.
    useEffect(() => {
        if (!profileData?.id) {
            reduxDispatch(completedTestsReset());
            return;
        }

        const loadCompletedTests = async () => {
            reduxDispatch(completedTestsLoading());

            try {
                const { data } = await fetchCompletedTests({
                    userId: profileData.id,
                });

                if (data?.status === "success") {
                    reduxDispatch(
                        completedTestsSuccess(
                            Array.isArray(data.tests) ? data.tests : []
                        )
                    );
                } else {
                    reduxDispatch(
                        completedTestsError(
                            data?.message ||
                                "Не удалось загрузить список пройденных тестов"
                        )
                    );
                }
            } catch (error) {
                const message =
                    error?.response?.data?.message ||
                    error?.message ||
                    "Ошибка загрузки пройденных тестов";
                reduxDispatch(completedTestsError(message));
            }
        };

        loadCompletedTests();
    }, [profileData?.id, reduxDispatch]);

    // Обработка OAuth linking callback.
    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const linkGoogle = urlParams.get("linkGoogle");
        const linkYandex = urlParams.get("linkYandex");
        const message = urlParams.get("message");

        if (linkGoogle !== null) {
            const success = linkGoogle === "true";

            if (success && message) {
                showAlert("success", decodeURIComponent(message));
                setProfileData((prev) => ({ ...prev, isGoogleAdded: true }));
            } else if (message) {
                showAlert("error", decodeURIComponent(message));
            }

            window.history.replaceState({}, document.title, "/account");
        }

        if (linkYandex !== null) {
            const success = linkYandex === "true";

            if (success && message) {
                showAlert("success", decodeURIComponent(message));
                setProfileData((prev) => ({ ...prev, isYandexAdded: true }));
            } else if (message) {
                showAlert("error", decodeURIComponent(message));
            }

            window.history.replaceState({}, document.title, "/account");
        }
    }, [showAlert, setProfileData]);

    const handleToggleTerminal = () => {
        reduxDispatch(toggleTerminal());
    };

    return {
        adminAccounts: state.adminAccounts,
        adminListError: state.adminListError,
        answersModal: state.answersModal,
        blockingUsers: state.blockingUsers,
        closeAnswersModal: handleCloseAnswersModal,
        closeTerminal: () => reduxDispatch(closeTerminal()),
        closeTestModal: () => reduxDispatch(closeTestModal()),
        completedTests: state.completedTests,
        completedTestsError: state.completedTestsError,
        deletingUsers: state.deletingUsers,
        emotionData: DEFAULT_EMOTION_DATA,
        handleAdminAction,
        handleBlockUser,
        handleChangePassword,
        handleDeleteAccount,
        handleDeleteUser,
        handleFieldChange,
        handleLinkProvider,
        handleLogout,
        handleProfileSave,
        handleStartTesting,
        handleToggleTerminal,
        hasGoogle,
        hasYandex,
        isAdmin,
        isAdminListLoading: state.isAdminListLoading,
        isCompletedTestsLoading: state.isCompletedTestsLoading,
        isDeletingAccount: state.isDeletingAccount,
        isSavingProfile: state.isSavingProfile,
        isTerminalOpen: state.isTerminalOpen,
        isTestModalOpen: state.isTestModalOpen,
        openAnswersModal: handleOpenAnswersModal,
        openTestModal: () => reduxDispatch(openTestModal()),
        profileData,
        selectedAnswers,
        setTerminalOpen: (value) => reduxDispatch(setTerminalOpen(value)),
        showLinkButtons,
    };
};
