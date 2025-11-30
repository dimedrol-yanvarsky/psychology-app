import React, { useState, useEffect } from "react";
import axios from "axios";
import styles from "./LoginPage.module.css";
import logoYandex from "../../pictures/yandex-logo.png";
import logoGoogle from "../../pictures/google-logo.png";
import { Link, useNavigate } from "react-router-dom";

const LoginPage = ({ showAlert, setIsAdmin, setIsAuth, setProfileData }) => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [recoveryStep, setRecoveryStep] = useState(1);
    const [recoveryLogin, setRecoveryLogin] = useState("");
    const [recoveryError, setRecoveryError] = useState("");

    const navigate = useNavigate();
    const loginApiBase = "http://localhost:8080/api/login";

    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const token = urlParams.get("token");
        const error = urlParams.get("error");

        if (token) {
            localStorage.setItem("token", token);
            window.location.href = "/account";
        }

        if (error) {
            showAlert("error", `Ошибка авторизации: ${error}`);
        }
    }, []);

    const handleSubmit = async (event) => {
        event.preventDefault();

        if (!email || !password) {
            showAlert("error", "Не оставляйте поля пустыми");
            return;
        }

        try {
            const response = await axios.post(`${loginApiBase}/password`, {
                email,
                password,
            });

            if (response.data.status === "Администратор") {
                setIsAdmin(true);
            }

            setProfileData(response.data);

            console.log(response)
            setIsAuth(true);
            showAlert("success", "Авторизация успешна");
            navigate("/account");
        } catch (error) {
            console.log(error);
            if (error.status === 401) {
                return showAlert("error", "Пароль неверный");
            }
            if (error.status === 500) {
                return showAlert("error", "База данных не отвечает");
            }
            if (error.status === 404) {
                return showAlert("error", "Пользователь не найден");
            }
            if (error.status === 400) {
                return showAlert("error", "Введите корректные данные");
            }
            if (error.status === 403) {
                return showAlert("error", "Доступ заблокирован");
            }
        }
    };

    const handleOAuth = async (provider) => {
        try {
            const { data } = await axios.post(`${loginApiBase}/${provider}`);

            if (data.error) {
                showAlert("error", data.error);
                return;
            }

            showAlert("success", "Авторизация выполнена");
            navigate("/account");
        } catch (error) {
            showAlert(
                "error",
                `Авторизация через ${provider} временно невозможна`
            );
        }
    };

    const openRecoveryModal = () => {
        setIsModalOpen(true);
        setRecoveryStep(1);
        setRecoveryLogin("");
        setRecoveryError("");
    };

    const closeRecoveryModal = () => {
        setIsModalOpen(false);
        setRecoveryStep(1);
        setRecoveryLogin("");
        setRecoveryError("");
    };

    const handleRecoverySubmit = (event) => {
        event.preventDefault();

        if (!recoveryLogin.trim()) {
            setRecoveryError("Введите логин или email, чтобы продолжить");
            return;
        }

        setRecoveryError("");
        setRecoveryStep(2);
        if (showAlert) {
            showAlert("success", "Инструкция по восстановлению отправлена");
        }
    };

    return (
        <div className={styles.page}>
            <div className={styles.layout}>
                <section className={styles.introCard}>
                    <p className={styles.overline}>Доступ к рекомендациям</p>
                    <h1 className={styles.title}>
                        Вернитесь к спокойной работе над собой
                    </h1>
                    <p className={styles.subtitle}>
                        Сохраняем ваш прогресс, подборки и результаты тестов.
                        Войдите, чтобы продолжить работу с рекомендациями и
                        личным кабинетом.
                    </p>
                    <div className={styles.pills}>
                        <span className={styles.pill}>Безопасный вход</span>
                        <span className={styles.pill}>Поддержка 24/7</span>
                        <span className={styles.pill}>
                            Синхронизация данных
                        </span>
                    </div>
                </section>

                <div className={styles.loginCard}>
                    <div className={styles.loginHeader}>
                        <span className={styles.badge}>Авторизация</span>
                        <div>
                            <h3 className={styles.cardTitle}>
                                Войдите любым удобным способом
                            </h3>
                            <p className={styles.cardSubtitle}>
                                Используйте корпоративный Google или Яндекс,
                                либо продолжите по email.
                            </p>
                        </div>
                    </div>

                    <div className={styles.oauthRow}>
                        <button
                            type="button"
                            className={styles.oauthButton}
                            onClick={() => handleOAuth("google")}
                        >
                            <img src={logoGoogle} alt="Google" />
                            <span>Войти через Google</span>
                        </button>
                        <button
                            type="button"
                            className={styles.oauthButton}
                            onClick={() => handleOAuth("yandex")}
                        >
                            <img src={logoYandex} alt="Яндекс" />
                            <span>Войти через Яндекс</span>
                        </button>
                    </div>

                    <div className={styles.divider}>
                        <span>или авторизация по почте</span>
                    </div>

                    <form
                        onSubmit={handleSubmit}
                        className={styles.form}
                        noValidate
                    >
                        <label className={styles.label} htmlFor="email">
                            Email
                        </label>
                        <input
                            className={styles.input}
                            type="email"
                            id="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder="Введите почтовый адрес"
                        />

                        <div className={styles.passwordRow}>
                            <label className={styles.label} htmlFor="password">
                                Пароль
                            </label>
                            <Link
                                to="#"
                                className={styles.linkButton}
                                onClick={(event) => {
                                    event.preventDefault();
                                    openRecoveryModal();
                                }}
                            >
                                Забыли пароль?
                            </Link>
                        </div>
                        <input
                            className={styles.input}
                            type="password"
                            id="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="Введите пароль"
                        />

                        <button type="submit" className={styles.primaryButton}>
                            Войти
                        </button>
                    </form>

                    <p className={styles.footerText}>
                        Нет аккаунта?{" "}
                        <Link to="/register" className={styles.linkButton}>
                            Зарегистрируйтесь
                        </Link>
                    </p>
                </div>
            </div>

            {isModalOpen && (
                <div
                    className={styles.modalOverlay}
                    onClick={closeRecoveryModal}
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
                            onClick={closeRecoveryModal}
                            aria-label="Закрыть модальное окно"
                        >
                            ×
                        </button>
                        {/* <p className={styles.modalOverline}>Модальное окно</p> */}
                        <h4 className={styles.modalTitle}>
                            Восстановление пароля
                        </h4>
                        <div className={styles.stepper}>
                            <div className={styles.stepLine}>
                                <div
                                    className={styles.stepLineFill}
                                    style={{
                                        width:
                                            recoveryStep === 2 ? "100%" : "50%",
                                    }}
                                />
                            </div>
                            <div className={styles.stepDots}>
                                {[1, 2].map((step) => (
                                    <div
                                        key={step}
                                        className={`${styles.stepDot} ${
                                            recoveryStep >= step
                                                ? styles.stepDotActive
                                                : ""
                                        }`}
                                    >
                                        <span>{step}</span>
                                        <p className={styles.stepLabel}>
                                            {step === 1
                                                ? "Логин"
                                                : "Подтверждение"}
                                        </p>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <form
                            className={styles.modalForm}
                            onSubmit={handleRecoverySubmit}
                        >
                            <label
                                className={styles.label}
                                htmlFor="recovery-login"
                            >
                                Введите логин или email
                            </label>
                            <input
                                id="recovery-login"
                                className={styles.input}
                                type="text"
                                value={recoveryLogin}
                                onChange={(e) =>
                                    setRecoveryLogin(e.target.value)
                                }
                                placeholder="Ваш логин"
                            />
                            {recoveryError && (
                                <p className={styles.modalError}>
                                    {recoveryError}
                                </p>
                            )}
                            <p className={styles.modalHint}>
                                {recoveryStep === 1
                                    ? "Отправим ссылку на восстановление пароля на указанный адрес."
                                    : "Проверьте почту и следуйте инструкции для обновления пароля."}
                            </p>
                            <button
                                type="submit"
                                className={styles.primaryButton}
                            >
                                Восстановить
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default LoginPage;
