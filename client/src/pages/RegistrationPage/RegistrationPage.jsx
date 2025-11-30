import React, { useState, useEffect } from "react";
import axios from "axios";
import styles from "./RegistrationPage.module.css";
import { Link, useNavigate } from "react-router-dom";
import { generateStrongPassword } from "./passwordGenerator";

const RegistrationPage = ({ showAlert }) => {
    const [login, setLogin] = useState("");
    const [name, setName] = useState("");
    const [password, setPassword] = useState("");
    const [passwordRepeated, setPasswordRepeated] = useState("");
    const [message, setMessage] = useState("");
    const [showPassword, setShowPassword] = useState(false);

    const navigate = useNavigate();

    const handleGeneratePassword = () => {
        const newPassword = generateStrongPassword();
        setPassword(newPassword);
        setPasswordRepeated(newPassword);
    };

    useEffect(() => {
        // Проверяем, есть ли токен в URL (после OAuth редиректа)
        const urlParams = new URLSearchParams(window.location.search);
        const token = urlParams.get("token");
        const error = urlParams.get("error");

        if (token) {
            localStorage.setItem("token", token);
            window.location.href = "/dashboard";
        }

        if (error) {
            setMessage(`Ошибка авторизации: ${error}`);
        }
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!(name && login && password && passwordRepeated)) {
            return showAlert("error", "Не оставляйте поля пустыми");
        }

        if (password !== passwordRepeated) {
            return showAlert("error", "Пароли отличаются");
        }

        try {
            const {response} = await axios.post(
                "http://localhost:8080/api/createAccount",
                {
                    name,
                    login,
                    password,
                    passwordRepeated,
                }
            );
            console.log(1, response);
            if (response.success) {
                showAlert("success", "Аккаунт зарегистрирован");
                navigate("/login");
            }
        } catch (error) {
            if (error.status === 409) {
                return showAlert("error", "Аккаунт уже существует");
            }
            if (error.status === 403) {
                return showAlert("error", "Доступ запрещен");
            }
            if (error.status === 401) {
                return showAlert("error", "Проверьте введенные данные");
            }
             if (error.status === 500) {
                return showAlert("error", "База данных недоступна");
            }
        }
    };

    const handleOAuth = (provider) => {
        window.location.href = `http://localhost:8080/auth/${provider}`;
    };

    return (
        <div className={styles.page}>
            <div className={styles.layout}>
                <section className={styles.introCard}>
                    <p className={styles.overline}>
                        Старт для новых пользователей
                    </p>
                    <h1 className={styles.title}>
                        Создайте аккаунт и продолжайте работу над собой
                    </h1>
                    <p className={styles.subtitle}>
                        Сохраняем подборки рекомендаций, результаты тестов и
                        личные цели. Зарегистрируйтесь, чтобы ничего не потерять
                        и возвращаться к прогрессу в один клик.
                    </p>
                    <div className={styles.pills}>
                        <span className={styles.pill}>Личный кабинет</span>
                        <span className={styles.pill}>
                            Сохранение прогресса
                        </span>
                        <span className={styles.pill}>Поддержка 24/7</span>
                    </div>
                </section>

                <div className={styles.registerCard}>
                    <div className={styles.registerHeader}>
                        <span className={styles.badge}>Регистрация</span>
                        <div>
                            <h3 className={styles.cardTitle}>
                                Создайте аккаунт за пару шагов
                            </h3>
                            <p className={styles.cardSubtitle}>
                                Заполните форму ниже. Можно сгенерировать
                                надежный пароль прямо здесь.
                            </p>
                        </div>
                    </div>

                    <form
                        onSubmit={handleSubmit}
                        className={styles.form}
                        noValidate
                    >
                        <label className={styles.label} htmlFor="name">
                            Как Вас зовут?
                        </label>
                        <input
                            className={styles.input}
                            type="text"
                            id="name"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Введите имя"
                            required
                        />

                        <label className={styles.label} htmlFor="email">
                            Email
                        </label>
                        <input
                            className={styles.input}
                            type="email"
                            id="email"
                            value={login}
                            onChange={(e) => setLogin(e.target.value)}
                            placeholder="Введите почтовый адрес"
                            required
                        />

                        <div className={styles.passwordRow}>
                            <label className={styles.label} htmlFor="password">
                                Придумайте пароль
                            </label>
                            <button
                                type="button"
                                className={styles.linkButton}
                                onClick={handleGeneratePassword}
                                title="Генерирует сложный пароль"
                            >
                                Сгенерировать
                            </button>
                        </div>
                        <input
                            className={styles.input}
                            type={showPassword ? "text" : "password"}
                            id="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="Введите пароль"
                            required
                        />
                        <p className={styles.helperText}>
                            Используйте буквы, цифры и символы для надежности.
                        </p>

                        <label
                            className={styles.label}
                            htmlFor="password-repeat"
                        >
                            Повторите пароль
                        </label>
                        <input
                            className={styles.input}
                            type={showPassword ? "text" : "password"}
                            id="password-repeat"
                            value={passwordRepeated}
                            onChange={(e) =>
                                setPasswordRepeated(e.target.value)
                            }
                            placeholder="Повторите пароль"
                            required
                        />

                        <button type="submit" className={styles.primaryButton}>
                            Зарегистрироваться
                        </button>
                    </form>

                    <p className={styles.footerText}>
                        Есть аккаунт?{" "}
                        <Link to="/login" className={styles.linkButton}>
                            Авторизуйтесь
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default RegistrationPage;
