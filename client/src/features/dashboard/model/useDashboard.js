import { useEffect, useMemo, useState } from "react";
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

const DEFAULT_EMOTION_DATA = [
    { id: "calm", label: "Спокойствие", value: 72 },
    { id: "energy", label: "Энергия", value: 54 },
    { id: "focus", label: "Фокус", value: 61 },
    { id: "stress", label: "Стресс", value: 28 },
];

export const useDashboard = ({
    showAlert,
    setIsAuth,
    setIsAdmin,
    isAdmin,
    profileData,
    setProfileData,
}) => {
    const navigate = useNavigate();
    const [user, setUser] = useState(null);
    const [isTestModalOpen, setIsTestModalOpen] = useState(false);
    const [isTerminalOpen, setIsTerminalOpen] = useState(false);
    const [adminAccounts, setAdminAccounts] = useState([]);
    const [isAdminListLoading, setIsAdminListLoading] = useState(
        profileData?.status === "Администратор"
    );
    const [adminListError, setAdminListError] = useState("");
    const [blockingUsers, setBlockingUsers] = useState({});
    const [deletingUsers, setDeletingUsers] = useState({});
    const [isDeletingAccount, setIsDeletingAccount] = useState(false);
    const [isSavingProfile, setIsSavingProfile] = useState(false);
    const [answersModal, setAnswersModal] = useState({
        open: false,
        loading: false,
        error: "",
        answers: [],
        questions: [],
        title: "",
    });
    const [completedTests, setCompletedTests] = useState([]);
    const [isCompletedTestsLoading, setIsCompletedTestsLoading] =
        useState(true);
    const [completedTestsError, setCompletedTestsError] = useState("");

    const selectedAnswers = useMemo(() => {
        const map = new Map();
        if (!Array.isArray(answersModal.answers)) {
            return map;
        }

        answersModal.answers.forEach((item) => {
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
    }, [answersModal.answers]);

    const providers = Array.isArray(user?.providers) ? user.providers : [];
    const providerValue = user?.provider || profileData?.provider || "";
    const hasGoogle =
        providers.includes("google") ||
        providerValue === "google" ||
        Boolean(user?.googleLinked);
    const hasYandex =
        providers.includes("yandex") ||
        providerValue === "yandex" ||
        Boolean(user?.yandexLinked);
    const showLinkButtons = !hasGoogle || !hasYandex;

    const handleFieldChange = (field, value) => {
        setProfileData((prev) => ({ ...prev, [field]: value }));
    };

    const handleProfileSave = async (event) => {
        event.preventDefault();

        if (isSavingProfile) {
            return;
        }

        const userId = profileData?.id;
        const firstName = (profileData?.firstName || "").trim();

        if (!userId) {
            showAlert &&
                showAlert(
                    "error",
                    "Не указан идентификатор пользователя для обновления"
                );
            return;
        }

        if (!firstName) {
            showAlert &&
                showAlert("error", "Введите имя, чтобы сохранить изменения");
            return;
        }

        setIsSavingProfile(true);

        try {
            const { data } = await changeUserData({
                userId,
                firstName,
            });

            if (data?.status === "success") {
                const updatedUser = data?.user || {};

                setProfileData((prev) => ({
                    ...prev,
                    id: updatedUser.id || userId,
                    firstName: updatedUser.firstName || firstName,
                    email: updatedUser.email || prev.email,
                    status: updatedUser.status || prev.status,
                }));

                setUser((prev) => ({
                    ...(prev || {}),
                    id: updatedUser.id || userId,
                    firstName: updatedUser.firstName || firstName,
                    email:
                        updatedUser.email ||
                        (prev && prev.email) ||
                        profileData.email,
                    status:
                        updatedUser.status ||
                        (prev && prev.status) ||
                        profileData.status,
                }));

                showAlert &&
                    showAlert(
                        "success",
                        data?.message || "Данные профиля обновлены"
                    );
            } else {
                showAlert &&
                    showAlert(
                        "error",
                        data?.message || "Не удалось сохранить изменения"
                    );
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось сохранить изменения";

            showAlert && showAlert("error", message);
        } finally {
            setIsSavingProfile(false);
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
        if (isDeletingAccount) {
            return;
        }

        if (!profileData?.id) {
            showAlert &&
                showAlert(
                    "error",
                    "Не указан идентификатор пользователя для удаления"
                );
            return;
        }

        setIsDeletingAccount(true);

        try {
            const { data } = await deleteAccount({
                userId: profileData.id,
            });

            if (data?.status === "success") {
                setIsAdmin(false);
                setIsAuth(false);
                setUser(null);
                setProfileData(createDefaultProfileData());

                showAlert &&
                    showAlert("success", data?.message || "Аккаунт удален");

                navigate(data?.redirect || "/login", { replace: true });
            } else {
                showAlert &&
                    showAlert(
                        "error",
                        data?.message || "Не удалось удалить аккаунт"
                    );
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось удалить аккаунт";

            showAlert && showAlert("error", message);
        } finally {
            setIsDeletingAccount(false);
        }
    };

    const handleChangePassword = () => {
        if (showAlert) {
            showAlert(
                "success",
                "Мы отправим ссылку для смены пароля на вашу почту"
            );
        }
    };

    const handleLinkProvider = (provider) => {
        return showAlert(
            "error",
            `Привязка аккаунта ${provider} будет доступна позже`
        );
    };

    const handleStartTesting = () => {
        setIsTestModalOpen(false);
        return showAlert("error", "Тестирование пока недоступно");
    };

    const handleAdminAction = (action, account) => {
        return (
            showAlert &&
            showAlert("error", `${action} для ${account.email} недоступно`)
        );
    };

    const handleBlockUser = async (account) => {
        if (!profileData?.id) {
            showAlert &&
                showAlert(
                    "error",
                    "Не указан идентификатор администратора для запроса"
                );
            return;
        }

        const targetId = account?.id;
        const status = (account?.status || "").trim();

        if (!targetId) {
            showAlert &&
                showAlert(
                    "error",
                    "Не удалось определить пользователя для блокировки"
                );
            return;
        }

        if (status === "Заблокирован" || status === "Удален") {
            return;
        }

        setBlockingUsers((prev) => ({ ...prev, [targetId]: true }));

        try {
            const { data } = await blockUser({
                adminId: profileData.id,
                targetUserId: targetId,
            });

            if (data?.status === "success") {
                setAdminAccounts((prev) =>
                    prev.map((user) =>
                        user.id === targetId
                            ? {
                                  ...user,
                                  status: data?.user?.status || "Заблокирован",
                              }
                            : user
                    )
                );

                showAlert &&
                    showAlert(
                        "success",
                        data?.message || "Пользователь заблокирован"
                    );
            } else {
                showAlert &&
                    showAlert(
                        "error",
                        data?.message ||
                            "Не удалось заблокировать пользователя"
                    );
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось заблокировать пользователя";

            showAlert && showAlert("error", message);
        } finally {
            setBlockingUsers((prev) => {
                const updated = { ...prev };
                delete updated[targetId];
                return updated;
            });
        }
    };

    const handleDeleteUser = async (account) => {
        if (!profileData?.id) {
            showAlert &&
                showAlert(
                    "error",
                    "Не указан идентификатор администратора для запроса"
                );
            return;
        }

        const targetId = account?.id;
        const status = (account?.status || "").trim();

        if (!targetId) {
            showAlert &&
                showAlert(
                    "error",
                    "Не удалось определить пользователя для удаления"
                );
            return;
        }

        if (status === "Удален") {
            return;
        }

        setDeletingUsers((prev) => ({ ...prev, [targetId]: true }));

        try {
            const { data } = await deleteUser({
                adminId: profileData.id,
                targetUserId: targetId,
            });

            if (data?.status === "success") {
                setAdminAccounts((prev) =>
                    prev.map((user) =>
                        user.id === targetId
                            ? {
                                  ...user,
                                  status: data?.user?.status || "Удален",
                              }
                            : user
                    )
                );

                showAlert &&
                    showAlert("success", data?.message || "Пользователь удален");
            } else {
                showAlert &&
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

            showAlert && showAlert("error", message);
        } finally {
            setDeletingUsers((prev) => {
                const updated = { ...prev };
                delete updated[targetId];
                return updated;
            });
        }
    };

    useEffect(() => {
        if (profileData?.status !== "Администратор") {
            setAdminAccounts([]);
            setIsAdminListLoading(false);
            setAdminListError("");
            return;
        }

        if (!profileData?.id) {
            setAdminAccounts([]);
            setIsAdminListLoading(false);
            setAdminListError("Не указан идентификатор пользователя");
            return;
        }

        const fetchAdminUsers = async () => {
            setIsAdminListLoading(true);
            setAdminListError("");

            try {
                const { data } = await fetchUsers({
                    userId: profileData.id,
                    status: profileData.status,
                });

                if (data?.status === "success") {
                    setAdminAccounts(
                        Array.isArray(data.users) ? data.users : []
                    );
                } else {
                    setAdminAccounts([]);
                    setAdminListError(
                        data?.message || "Не удалось загрузить пользователей"
                    );
                }
            } catch (error) {
                const message =
                    error?.response?.data?.message ||
                    error?.message ||
                    "Ошибка загрузки пользователей";
                setAdminAccounts([]);
                setAdminListError(message);
            } finally {
                setIsAdminListLoading(false);
            }
        };

        fetchAdminUsers();
    }, [profileData?.id, profileData?.status]);

    const openAnswersModal = async (test) => {
        setAnswersModal({
            open: true,
            loading: true,
            error: "",
            answers: [],
            questions: [],
            title: test.testName || "Пройденный тест",
        });

        if (!profileData?.id || !test?.id || !test?.testId) {
            setAnswersModal((prev) => ({
                ...prev,
                loading: false,
                error: "Недостаточно данных для загрузки результата",
            }));
            return;
        }

        try {
            const { data } = await fetchUserAnswers({
                userId: profileData.id,
                completedTestId: test.id,
                testId: test.testId,
            });

            if (data?.status === "success") {
                setAnswersModal((prev) => ({
                    ...prev,
                    loading: false,
                    answers: Array.isArray(data.answers) ? data.answers : [],
                    questions: Array.isArray(data.questions)
                        ? data.questions
                        : [],
                }));
            } else {
                setAnswersModal((prev) => ({
                    ...prev,
                    loading: false,
                    error:
                        data?.message ||
                        "Не удалось загрузить ответы пользователя",
                }));
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Ошибка загрузки ответов пользователя";
            setAnswersModal((prev) => ({
                ...prev,
                loading: false,
                error: message,
            }));
        }
    };

    const closeAnswersModal = () => {
        setAnswersModal({
            open: false,
            loading: false,
            error: "",
            answers: [],
            questions: [],
            title: "",
        });
    };

    useLockBodyScroll(answersModal.open);

    useEffect(() => {
        if (!profileData?.id) {
            setIsCompletedTestsLoading(false);
            setCompletedTests([]);
            return;
        }

        const loadCompletedTests = async () => {
            setIsCompletedTestsLoading(true);
            setCompletedTestsError("");

            try {
                const { data } = await fetchCompletedTests({
                    userId: profileData.id,
                });

                if (data?.status === "success") {
                    setCompletedTests(
                        Array.isArray(data.tests) ? data.tests : []
                    );
                } else {
                    setCompletedTests([]);
                    setCompletedTestsError(
                        data?.message ||
                            "Не удалось загрузить список пройденных тестов"
                    );
                }
            } catch (error) {
                const message =
                    error?.response?.data?.message ||
                    error?.message ||
                    "Ошибка загрузки пройденных тестов";
                setCompletedTests([]);
                setCompletedTestsError(message);
            } finally {
                setIsCompletedTestsLoading(false);
            }
        };

        loadCompletedTests();
    }, [profileData?.id]);

    const handleToggleTerminal = () => {
        setIsTerminalOpen((prev) => !prev);
    };

    const openTestModal = () => setIsTestModalOpen(true);
    const closeTestModal = () => setIsTestModalOpen(false);
    const closeTerminal = () => setIsTerminalOpen(false);

    return {
        adminAccounts,
        adminListError,
        answersModal,
        blockingUsers,
        closeAnswersModal,
        closeTerminal,
        closeTestModal,
        completedTests,
        completedTestsError,
        deletingUsers,
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
        isAdminListLoading,
        isCompletedTestsLoading,
        isDeletingAccount,
        isSavingProfile,
        isTerminalOpen,
        isTestModalOpen,
        openAnswersModal,
        openTestModal,
        profileData,
        selectedAnswers,
        setTerminalOpen: setIsTerminalOpen,
        showLinkButtons,
    };
};
