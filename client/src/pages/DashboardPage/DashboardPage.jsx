import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Terminal from "../../components/Terminal/Terminal";
import logoYandex from "../../pictures/yandex-logo.png";
import logoGoogle from "../../pictures/google-logo.png";
import styles from "./DashboardPage.module.css";

const defaultAdminAccounts = [
    {
        id: 1,
        firstName: "Анна",
        lastName: "Иванова",
        email: "anna.ivanova@example.com",
    },
    {
        id: 2,
        firstName: "Дмитрий",
        lastName: "Кузнецов",
        email: "d.kuznetsov@example.com",
    },
    {
        id: 3,
        firstName: "София",
        lastName: "Лебедева",
        email: "s.lebedeva@example.com",
    },
];

const DashboardPage = ({
    showAlert,
    setIsAuth,
    setIsAdmin,
    isAdmin,
    profileData,
    setProfileData,
}) => {
    const navigate = useNavigate();
    const [user, setUser] = useState(null);
    const [isLoading] = useState(false);
    const [isTestModalOpen, setIsTestModalOpen] = useState(false);
    const [isTerminalOpen, setIsTerminalOpen] = useState(false);
    const [adminAccounts] = useState(defaultAdminAccounts);
    const [completedTests] = useState([
        {
            id: "t1",
            title: "Шкала депрессии Бека",
            score: "Средний уровень",
            date: "12.04.2024",
        },
        {
            id: "t2",
            title: "Опросник Спилбергера-Ханина",
            score: "Низкая тревожность",
            date: "03.03.2024",
        },
    ]);
    const [emotionData] = useState([
        { id: "calm", label: "Спокойствие", value: 72 },
        { id: "energy", label: "Энергия", value: 54 },
        { id: "focus", label: "Фокус", value: 61 },
        { id: "stress", label: "Стресс", value: 28 },
    ]);

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

    const handleProfileSave = (event) => {
        event.preventDefault();
        setUser((prev) => ({ ...(prev || {}), ...profileData }));

        if (showAlert) {
            showAlert("success", "Данные профиля обновлены");
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

    const handleDeleteAccount = () => {
        showAlert("error", "Удаление аккаунта будет доступно позже");
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
        if (showAlert) {
            showAlert(
                "error",
                `Привязка аккаунта ${provider} будет доступна позже`
            );
        }
    };

    const handleStartTesting = () => {
        setIsTestModalOpen(false);
        navigate("/tests");
    };

    const handleAdminAction = (action, account) => {
        if (showAlert) {
            showAlert(
                "success",
                `${action} для ${account.firstName} запрошено`
            );
        }
    };

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
                                Обновляйте имя, фамилию и почту. Мы бережно
                                храним изменения и используем их в
                                рекомендациях.
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

                            <label className={styles.label} htmlFor="last-name">
                                Фамилия
                            </label>
                            <input
                                id="last-name"
                                className={styles.input}
                                type="text"
                                value={profileData.lastName}
                                onChange={(event) =>
                                    handleFieldChange(
                                        "lastName",
                                        event.target.value
                                    )
                                }
                                placeholder="Фамилия"
                            />

                            <label className={styles.label} htmlFor="email">
                                Почтовый адрес
                            </label>
                            <input
                                id="email"
                                className={styles.input}
                                type="email"
                                value={profileData.email}
                                onChange={(event) =>
                                    handleFieldChange(
                                        "email",
                                        event.target.value
                                    )
                                }
                                placeholder="example@domain.com"
                            />

                            <button
                                type="submit"
                                className={styles.primaryButton}
                            >
                                Сохранить изменения
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
                                >
                                    Удалить аккаунт
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

                    {profileData.psychotype ? (
                        <div className={styles.psychotypeContent}>
                            <div className={styles.psychotypeTag}>
                                {profileData.psychotype}
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

                    {completedTests.length === 0 ? (
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
                                <div key={test.id} className={styles.testRow}>
                                    <div className={styles.testInfo}>
                                        <div className={styles.testTitle}>
                                            {test.title}
                                        </div>
                                        <div className={styles.testMeta}>
                                            {test.date} · {test.score}
                                        </div>
                                    </div>
                                    <Link
                                        to="/tests"
                                        className={styles.linkButton}
                                    >
                                        Открыть
                                    </Link>
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

                <div className={styles.quickGrid}>
                    <Link to="/recommendations" className={styles.quickCard}>
                        <div className={styles.quickIcon}>🎯</div>
                        <div>
                            <h4 className={styles.quickTitle}>Рекомендации</h4>
                            <p className={styles.quickText}>
                                Персональные подборки с учетом текущего статуса
                                и результатов тестов.
                            </p>
                        </div>
                    </Link>

                    <Link to="/tests" className={styles.quickCard}>
                        <div className={styles.quickIcon}>📊</div>
                        <div>
                            <h4 className={styles.quickTitle}>Тестирования</h4>
                            <p className={styles.quickText}>
                                Запустите новое тестирование или вернитесь к уже
                                пройденным.
                            </p>
                        </div>
                    </Link>
                </div>

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
                            {adminAccounts.map((account) => (
                                <div
                                    key={account.id}
                                    className={styles.adminUserCard}
                                >
                                    <div className={styles.adminUserInfo}>
                                        <div className={styles.adminUserName}>
                                            {account.firstName}{" "}
                                            {account.lastName}
                                        </div>
                                        <div className={styles.adminUserEmail}>
                                            {account.email}
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
                                            className={styles.secondaryButton}
                                            onClick={() =>
                                                handleAdminAction(
                                                    "Блокировка аккаунта",
                                                    account
                                                )
                                            }
                                        >
                                            Заблокировать аккаунт
                                        </button>
                                        <button
                                            type="button"
                                            className={`${styles.cardActionButton} ${styles.cardDeleteButton}`}
                                            onClick={() =>
                                                handleAdminAction(
                                                    "Удаление аккаунта",
                                                    account
                                                )
                                            }
                                        >
                                            Удалить аккаунт
                                        </button>
                                    </div>
                                </div>
                            ))}
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
                        <Terminal profileData={profileData} setIsTerminalOpen={setIsTerminalOpen}/>
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
        </div>
    );
};

export default DashboardPage;
