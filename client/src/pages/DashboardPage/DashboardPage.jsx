import React, { useEffect, useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import axios from "axios";
import Terminal from "../../components/Terminal/Terminal";
import logoYandex from "../../pictures/yandex-logo.png";
import logoGoogle from "../../pictures/google-logo.png";
import styles from "./DashboardPage.module.css";

const DashboardPage = ({
    showAlert,
    setIsAuth,
    setIsAdmin,
    isAdmin,
    profileData,
    setProfileData,
}) => {
    const navigate = useNavigate();
    const dashboardApiBase = "http://localhost:8080/api/dashboard";
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
    const [emotionData] = useState([
        { id: "calm", label: "Спокойствие", value: 72 },
        { id: "energy", label: "Энергия", value: 54 },
        { id: "focus", label: "Фокус", value: 61 },
        { id: "stress", label: "Стресс", value: 28 },
    ]);

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
            const { data } = await axios.post(
                `${dashboardApiBase}/change-user-data`,
                {
                    userId,
                    firstName,
                }
            );

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
        // localStorage.removeItem("token");
        setIsAuth(false);
        setIsAdmin(false);
        setProfileData({
            firstName: "",
            lastName: "",
            email: "",
            psychotype: "",
        });
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
            const { data } = await axios.post(
                `${dashboardApiBase}/delete-account`,
                {
                    userId: profileData.id,
                }
            );

            if (data?.status === "success") {
                setIsAdmin(false);
                setIsAuth(false);
                setUser(null);
                setProfileData({
                    id: null,
                    firstName: "",
                    lastName: "",
                    email: "",
                    status: "",
                    psychoType: "",
                    date: "",
                    isGoogleAdded: false,
                    isYandexAdded: false,
                });

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
        // navigate("/tests");
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
            const { data } = await axios.post(
                `${dashboardApiBase}/block-user`,
                {
                    adminId: profileData.id,
                    targetUserId: targetId,
                }
            );

            if (data?.status === "success") {
                setAdminAccounts((prev) =>
                    prev.map((user) =>
                        user.id === targetId
                            ? {
                                  ...user,
                                  status:
                                      data?.user?.status || "Заблокирован",
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
            const { data } = await axios.post(
                `${dashboardApiBase}/delete-user`,
                {
                    adminId: profileData.id,
                    targetUserId: targetId,
                }
            );

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

        const fetchUsers = async () => {
            setIsAdminListLoading(true);
            setAdminListError("");

            try {
                const { data } = await axios.post(`${dashboardApiBase}/users`, {
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

        fetchUsers();
    }, [dashboardApiBase, profileData?.id, profileData?.status]);

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
            console.log("Тест id");
            console.log(test.id);
            const { data } = await axios.post(
                `${dashboardApiBase}/user-answers`,
                {
                    userId: profileData.id,
                    completedTestId: test.id,
                    testId: test.testId,
                }
            );

            console.log("Данные теста:");
            console.log(data);
            console.log("Ответы");
            console.log(data.answers);
            console.log("Вопросы");
            console.log(data.questions);

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

    useEffect(() => {
        if (answersModal.open) {
            const originalOverflow = document.body.style.overflow;
            document.body.style.overflow = "hidden";

            return () => {
                document.body.style.overflow = originalOverflow;
            };
        }

        return undefined;
    }, [answersModal.open]);

    useEffect(() => {
        if (!profileData?.id) {
            setIsCompletedTestsLoading(false);
            setCompletedTests([]);
            return;
        }

        const fetchCompletedTests = async () => {
            setIsCompletedTestsLoading(true);
            setCompletedTestsError("");

            try {
                const { data } = await axios.post(
                    `${dashboardApiBase}/completed-tests`,
                    { userId: profileData.id }
                );

                if (data?.status === "success") {
                    console.log(data.tests);
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

        fetchCompletedTests();
    }, [dashboardApiBase, profileData?.id]);

    return (
        <div className={styles.page}>
            <div className={styles.layout}>
                <section className={styles.heroCard}>
                    <p className={styles.overline}>Личный кабинет</p>
                    <h1 className={styles.title}>
                        Управляйте профилем, тестированиями и безопасностью в
                        одном месте
                    </h1>
                    <p className={styles.subtitle}>
                        Редактируйте основные данные, следите за статусом и
                        запускайте тесты на психотип. Все действия -
                        централизовано и в едином стиле.
                    </p>
                    <div className={styles.pills}>
                        <span className={styles.pill}>
                            Шифрование данных с https
                        </span>
                        <span className={styles.pill}>
                            Хеширование паролей с bcrypt
                        </span>
                        <span className={styles.pill}>Поддержка CI/CD</span>
                    </div>
                </section>

                <section className={styles.profileCard}>
                    <div className={styles.sectionHead}>
                        <div>
                            <span className={styles.badge}>Профиль</span>
                            <h3 className={styles.cardTitle}>
                                Персональные данные
                            </h3>
                            <p className={styles.cardSubtitle}>
                                Обновляйте имя, а почта закреплена за
                                аккаунтом. Мы бережно храним изменения и
                                используем их в рекомендациях.
                            </p>
                        </div>
                    </div>

                    <div className={styles.profileGrid}>
                        <form
                            className={styles.form}
                            onSubmit={handleProfileSave}
                        >
                            <label
                                className={styles.label}
                                htmlFor="first-name"
                            >
                                Имя
                            </label>
                            <input
                                id="first-name"
                                className={styles.input}
                                type="text"
                                value={profileData.firstName}
                                onChange={(event) =>
                                    handleFieldChange(
                                        "firstName",
                                        event.target.value
                                    )
                                }
                                placeholder="Имя"
                            />

                            <label className={styles.label} htmlFor="email">
                                Почтовый адрес
                            </label>
                            <input
                                id="email"
                                className={styles.input}
                                type="email"
                                value={profileData.email}
                                placeholder="example@domain.com"
                                readOnly
                                aria-readonly="true"
                            />

                            <button
                                type="submit"
                                className={styles.primaryButton}
                                disabled={isSavingProfile}
                            >
                                {isSavingProfile
                                    ? "Сохраняем..."
                                    : "Сохранить изменения"}
                            </button>
                        </form>

                        <div className={styles.infoPanel}>
                            <div className={styles.infoRow}>
                                <span className={styles.infoLabel}>Статус</span>
                                <span className={styles.infoValue}>
                                    {profileData.status}
                                </span>
                            </div>
                            <div className={styles.actionsRow}>
                                <button
                                    type="button"
                                    className={`${styles.cardActionButton} ${styles.cardEditButton}`}
                                    onClick={handleChangePassword}
                                >
                                    Сменить пароль
                                </button>
                                <button
                                    type="button"
                                    className={styles.secondaryButton}
                                    onClick={handleLogout}
                                >
                                    Выйти из аккаунта
                                </button>
                                <button
                                    type="button"
                                    className={`${styles.cardActionButton} ${styles.cardDeleteButton}`}
                                    onClick={handleDeleteAccount}
                                    disabled={isDeletingAccount}
                                >
                                    {isDeletingAccount
                                        ? "Удаляем..."
                                        : "Удалить аккаунт"}
                                </button>
                            </div>
                        </div>

                        {showLinkButtons && (
                            <div className={styles.infoPanel}>
                                <div className={styles.infoRow}>
                                    <span className={styles.infoLabel}>
                                        Привязать аккаунт
                                    </span>
                                </div>
                                <div className={styles.connectRow}>
                                    {!hasGoogle && (
                                        <button
                                            type="button"
                                            className={styles.oauthButton}
                                            onClick={() =>
                                                handleLinkProvider("Google")
                                            }
                                        >
                                            <img
                                                src={logoGoogle}
                                                alt="Google"
                                            />
                                            <span>
                                                Привязать аккаунт Google
                                            </span>
                                        </button>
                                    )}
                                    {!hasYandex && (
                                        <button
                                            type="button"
                                            className={styles.oauthButton}
                                            onClick={() =>
                                                handleLinkProvider("Яндекс")
                                            }
                                        >
                                            <img
                                                src={logoYandex}
                                                alt="Yandex"
                                            />
                                            <span>
                                                Привязать аккаунт Яндекс
                                            </span>
                                        </button>
                                    )}
                                </div>
                            </div>
                        )}
                    </div>
                </section>

                <section className={styles.psychotypeCard}>
                    <div className={styles.sectionHead}>
                        <div>
                            <span className={styles.badge}>Психотип</span>
                            <h3 className={styles.cardTitle}>
                                Результаты тестирования на психотип
                            </h3>
                            <p className={styles.cardSubtitle}>
                                Мы покажем Ваш психотип или предложим пройти
                                тестирование, если Вы еще этого не сделали.
                            </p>
                        </div>
                    </div>

                    {profileData.psychoType ? (
                        <div className={styles.psychotypeContent}>
                            <div className={styles.psychotypeTag}>
                                {profileData.psychoType}
                            </div>
                            <p className={styles.subtitle}>
                                Вы всегда можете обновить результат, пройдя тест
                                повторно.
                            </p>
                            <div className={styles.retakeRow}>
                                <button
                                    type="button"
                                    className={styles.secondaryButton}
                                    onClick={() => setIsTestModalOpen(true)}
                                >
                                    Пройти заново
                                </button>
                                <Link to="/tests" className={styles.linkButton}>
                                    Открыть все тесты
                                </Link>
                            </div>
                        </div>
                    ) : (
                        <div className={styles.psychotypeEmpty}>
                            <div>
                                <h4 className={styles.emptyTitle}>
                                    Нет сохраненного психотипа
                                </h4>
                                <p className={styles.cardSubtitle}>
                                    Пройдите короткое тестирование, чтобы
                                    получить персональные рекомендации и дерево
                                    эмоций.
                                </p>
                            </div>
                            <button
                                type="button"
                                className={styles.primaryButton}
                                onClick={() => setIsTestModalOpen(true)}
                            >
                                Пройти тестирование
                            </button>
                        </div>
                    )}
                </section>

                <section className={styles.psychotypeCard}>
                    <div className={styles.sectionHead}>
                        <div>
                            <span className={styles.badge}>
                                Мои тестирования
                            </span>
                            <h3 className={styles.cardTitle}>
                                Пройденные тесты
                            </h3>
                            <p className={styles.cardSubtitle}>
                                Отслеживайте завершенные тестирования и
                                возвращайтесь к результатам при необходимости.
                            </p>
                        </div>
                    </div>

                    {isCompletedTestsLoading ? (
                        <div className={styles.testsLoader}>
                            <div
                                className={styles.loaderSpinner}
                                aria-hidden="true"
                            />
                            <span className={styles.loaderText}>
                                Загрузка...
                            </span>
                        </div>
                    ) : completedTestsError ? (
                        <div className={styles.testsError}>
                            {completedTestsError}
                        </div>
                    ) : completedTests.length === 0 ? (
                        <div className={styles.psychotypeEmpty}>
                            <div>
                                <h4 className={styles.emptyTitle}>
                                    Пока нет пройденных тестов
                                </h4>
                                <p className={styles.cardSubtitle}>
                                    Пройдите тестирование, чтобы мы могли
                                    показать историю результатов.
                                </p>
                            </div>
                            <Link
                                to="/tests"
                                className={styles.primaryButtonLink}
                            >
                                Перейти к тестам
                            </Link>
                        </div>
                    ) : (
                        <div className={styles.testsList}>
                            {completedTests.map((test) => (
                                <div
                                    key={test.testId || test.id}
                                    className={styles.testRow}
                                >
                                    <div className={styles.testInfo}>
                                        <div className={styles.testTitle}>
                                            {test.testName}
                                        </div>
                                        <div className={styles.testMeta}>
                                            {test.date || "Дата не указана"} ·{" "}
                                            {test.result || "Без результата"}
                                        </div>
                                    </div>
                                    <button
                                        type="button"
                                        className={styles.linkButton}
                                        onClick={() => openAnswersModal(test)}
                                    >
                                        Открыть
                                    </button>
                                </div>
                            ))}
                        </div>
                    )}
                </section>

                <section className={styles.psychotypeCard}>
                    <div className={styles.sectionHead}>
                        <div>
                            <span className={styles.badge}>Мои эмоции</span>
                            <h3 className={styles.cardTitle}>
                                Граф эмоционального состояния
                            </h3>
                            <p className={styles.cardSubtitle}>
                                Сводка по тестированиям помогает понять динамику
                                вашего состояния.
                            </p>
                        </div>
                    </div>

                    {!profileData.psychotype ? (
                        <div className={styles.psychotypeEmpty}>
                            <div>
                                <h4 className={styles.emptyTitle}>
                                    Нет данных для построения графика
                                </h4>
                                <p className={styles.cardSubtitle}>
                                    Пройдите тест на психотип, чтобы увидеть
                                    динамику эмоций.
                                </p>
                            </div>
                            <button
                                type="button"
                                className={styles.primaryButton}
                                onClick={() => setIsTestModalOpen(true)}
                            >
                                Пройти тестирование
                            </button>
                        </div>
                    ) : (
                        <div className={styles.emotionGraph}>
                            {emotionData.map((item) => (
                                <div
                                    key={item.id}
                                    className={styles.emotionRow}
                                >
                                    <span className={styles.emotionLabel}>
                                        {item.label}
                                    </span>
                                    <div className={styles.emotionBarTrack}>
                                        <div
                                            className={styles.emotionBar}
                                            style={{ width: `${item.value}%` }}
                                        />
                                    </div>
                                    <span className={styles.emotionValue}>
                                        {item.value}%
                                    </span>
                                </div>
                            ))}
                        </div>
                    )}
                </section>

                {isAdmin && (
                    <section className={styles.adminPanel}>
                        <div className={styles.sectionHead}>
                            <div>
                                <span className={styles.badge}>
                                    Панель администратора
                                </span>
                                <h3 className={styles.cardTitle}>
                                    Управление пользователями
                                </h3>
                                <p className={styles.cardSubtitle}>
                                    Просматривайте зарегистрированные аккаунты,
                                    блокируйте или удаляйте доступ, открывайте
                                    тестирования и дерево эмоций.
                                </p>
                            </div>
                            <button
                                type="button"
                                className={`${styles.cardActionButton} ${styles.cardEditButton}`}
                                onClick={() =>
                                    setIsTerminalOpen((prev) => !prev)
                                }
                            >
                                {isTerminalOpen
                                    ? "Скрыть терминал"
                                    : "Открыть терминал"}
                            </button>
                        </div>

                        <div className={styles.adminList}>
                            {isAdminListLoading ? (
                                <div className={styles.adminLoader}>
                                    <div
                                        className={styles.loaderSpinner}
                                        aria-hidden="true"
                                    />
                                    <span className={styles.loaderText}>
                                        Загрузка данных...
                                    </span>
                                </div>
                            ) : adminListError ? (
                                <div className={styles.adminError}>
                                    {adminListError}
                                </div>
                            ) : adminAccounts.length === 0 ? (
                                <div className={styles.adminEmpty}>
                                    Пользователи не найдены
                                </div>
                            ) : (
                                adminAccounts.map((account) => {
                                    const status = (
                                        account.status || ""
                                    ).trim();
                                    const isBlocked =
                                        status === "Заблокирован";
                                    const isDeleted =
                                        status === "Удален";
                                    const isBlocking = Boolean(
                                        blockingUsers[account.id]
                                    );
                                    const isDeleting = Boolean(
                                        deletingUsers[account.id]
                                    );
                                    const blockButtonClass = `${styles.secondaryButton} ${
                                        isBlocked ||
                                        isBlocking ||
                                        isDeleted ||
                                        isDeleting
                                            ? styles.blockedButton
                                            : ""
                                    }`;
                                    const deleteButtonClass = `${styles.cardActionButton} ${
                                        isDeleted
                                            ? styles.deletedButton
                                            : styles.cardDeleteButton
                                    }`;

                                    return (
                                        <div
                                            key={account.id}
                                            className={styles.adminUserCard}
                                        >
                                            <div className={styles.adminUserInfo}>
                                                <div
                                                    className={
                                                        styles.adminUserName
                                                    }
                                                >
                                                    {account.firstName ||
                                                        "Без имени"}
                                                    {account.lastName
                                                        ? ` ${account.lastName}`
                                                        : ""}
                                                </div>
                                                <div
                                                    className={
                                                        styles.adminUserEmail
                                                    }
                                                >
                                                    {account.email || "—"}
                                                </div>
                                            </div>
                                            <div className={styles.adminActions}>
                                                <button
                                                    type="button"
                                                    className={`${styles.cardActionButton} ${styles.cardPrimaryButton}`}
                                                    onClick={() =>
                                                        handleAdminAction(
                                                            "Просмотр тестирований",
                                                            account
                                                        )
                                                    }
                                                >
                                                    Просмотреть тестирования
                                                </button>
                                                <button
                                                    type="button"
                                                    className={`${styles.cardActionButton} ${styles.cardPrimaryButton}`}
                                                    onClick={() =>
                                                        handleAdminAction(
                                                            "Дерево эмоций",
                                                            account
                                                        )
                                                    }
                                                >
                                                    Дерево эмоций
                                                </button>
                                                <button
                                                    type="button"
                                                    className={blockButtonClass}
                                                    onClick={() =>
                                                        handleBlockUser(account)
                                                    }
                                                    disabled={
                                                        isBlocked ||
                                                        isBlocking ||
                                                        isDeleted ||
                                                        isDeleting
                                                    }
                                                >
                                                    {isDeleted
                                                        ? "Удален"
                                                        : isBlocked
                                                        ? "Заблокирован"
                                                        : isBlocking
                                                        ? "Блокируем..."
                                                        : "Заблокировать аккаунт"}
                                                </button>
                                                <button
                                                    type="button"
                                                    className={deleteButtonClass}
                                                    onClick={() =>
                                                        handleDeleteUser(account)
                                                    }
                                                    disabled={
                                                        isDeleted || isDeleting
                                                    }
                                                >
                                                    {isDeleted
                                                        ? "Удален"
                                                        : isDeleting
                                                        ? "Удаляем..."
                                                        : "Удалить аккаунт"}
                                                </button>
                                            </div>
                                        </div>
                                    );
                                })
                            )}
                        </div>
                    </section>
                )}
            </div>

            {isTerminalOpen && (
                <div
                    className={styles.terminalOverlay}
                    onClick={() => setIsTerminalOpen(false)}
                >
                    <div
                        className={styles.terminalModal}
                        onClick={(event) => event.stopPropagation()}
                        role="dialog"
                        aria-modal="true"
                    >
                        <Terminal
                            profileData={profileData}
                            setIsTerminalOpen={setIsTerminalOpen}
                        />
                    </div>
                </div>
            )}

            {isTestModalOpen && (
                <div
                    className={styles.modalOverlay}
                    onClick={() => setIsTestModalOpen(false)}
                >
                    <div
                        className={styles.modal}
                        onClick={(event) => event.stopPropagation()}
                        role="dialog"
                        aria-modal="true"
                    >
                        <button
                            type="button"
                            className={styles.modalClose}
                            onClick={() => setIsTestModalOpen(false)}
                            aria-label="Закрыть модальное окно"
                        >
                            ×
                        </button>
                        <p className={styles.modalOverline}>Психотип</p>
                        <h4 className={styles.modalTitle}>
                            Тестирование на психотип
                        </h4>
                        <p className={styles.modalText}>
                            Краткий опрос помогает понять эмоциональные реакции
                            и предпочтения. Результат влияет на рекомендации,
                            дерево эмоций и подборку тестов.
                        </p>
                        <div className={styles.modalChips}>
                            <span className={styles.pill}>15-20 минут</span>
                            <span className={styles.pill}>30 вопросов</span>
                            <span className={styles.pill}>Результат сразу</span>
                        </div>
                        <button
                            type="button"
                            className={styles.primaryButton}
                            onClick={handleStartTesting}
                        >
                            Начать тестирование
                        </button>
                    </div>
                </div>
            )}

            {answersModal.open && (
                <div
                    className={styles.modalOverlay}
                    onClick={closeAnswersModal}
                >
                    <div
                        className={styles.answersModal}
                        onClick={(event) => event.stopPropagation()}
                        role="dialog"
                        aria-modal="true"
                    >
                        <button
                            type="button"
                            className={styles.modalClose}
                            onClick={closeAnswersModal}
                            aria-label="Закрыть модальное окно"
                        >
                            ×
                        </button>
                        <p className={styles.modalOverline}>Пройденный тест</p>
                        <h4 className={styles.modalTitle}>
                            {answersModal.title || "Результаты теста"}
                        </h4>
                        {answersModal.loading ? (
                            <div className={styles.testsLoader}>
                                <div
                                    className={styles.loaderSpinner}
                                    aria-hidden="true"
                                />
                                <span className={styles.loaderText}>
                                    Загрузка...
                                </span>
                            </div>
                        ) : answersModal.error ? (
                            <div className={styles.testsError}>
                                {answersModal.error}
                            </div>
                        ) : (
                            <div className={styles.questionsList}>
                                    {answersModal.questions.map(
                                        (question, questionIndex) => {
                                            console.log('Вывод вопроса')
                                            console.log(question)
                                            const selectType =
                                                question.selectype ||
                                                question.selectType ||
                                                "";
                                            const isSingle = selectType === "one";
                                            const questionNumber =
                                                question.id ||
                                                question.number ||
                                                questionIndex + 1;
                                            const selected =
                                                selectedAnswers.get(
                                                    Number(questionNumber)
                                                ) || new Set();

                                        const rawQuestionText =
                                            question.questionBody ||
                                            question.question ||
                                            question.body ||
                                            "";
                                        const questionText =
                                            typeof rawQuestionText === "string"
                                                ? rawQuestionText
                                                : rawQuestionText?.body ||
                                                  rawQuestionText?.text ||
                                                  rawQuestionText?.title ||
                                                  "Вопрос";

                                        const options = Array.isArray(
                                            question.answerOptions
                                        )
                                            ? question.answerOptions
                                            : Array.isArray(question.answers)
                                            ? question.answers
                                            : [];
                                        return (
                                            <div
                                                key={
                                                    question.id || questionIndex
                                                }
                                                className={styles.questionCard}
                                            >
                                                <div
                                                className={
                                                    styles.questionTitle
                                                }
                                            >
                                                {questionNumber}.{" "}
                                                    {questionText}
                                                </div>
                                                <div
                                                    className={
                                                        styles.optionsList
                                                    }
                                                >
                                                    {options.map(
                                                        (
                                                            option,
                                                            optionIndex
                                                        ) => {
                                                            const optionNumber =
                                                                optionIndex + 1;
                                                            const optionLabel =
                                                                typeof option ===
                                                                "string"
                                                                    ? option
                                                                    : option?.body ||
                                                                      option
                                                                          ?.text ||
                                                                      option
                                                                          ?.title ||
                                                                      option?.label ||
                                                                      String(
                                                                          optionNumber
                                                                      );
                                                            const isChecked =
                                                                selected instanceof
                                                                    Set &&
                                                                selected.has(
                                                                    optionNumber
                                                                );
                                                            const inputType =
                                                                isSingle
                                                                    ? "radio"
                                                                    : "checkbox";

                                                            return (
                                                                <label
                                                                    key={`${question.id || questionIndex}-${option.id || optionIndex}`}
                                                                    className={`${
                                                                        styles.optionItem
                                                                    } ${
                                                                        isChecked
                                                                            ? styles.optionSelected
                                                                            : ""
                                                                    }`}
                                                                >
                                                                    <input
                                                                        type={
                                                                            inputType
                                                                        }
                                                                        name={`question-${questionNumber}`}
                                                                        checked={
                                                                            isChecked
                                                                        }
                                                                        disabled
                                                                        readOnly
                                                                    />
                                                                    <span>
                                                                        {optionLabel}
                                                                    </span>
                                                                </label>
                                                            );
                                                        }
                                                    )}
                                                </div>
                                            </div>
                                        );
                                    }
                                )}
                            </div>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
};

export default DashboardPage;
