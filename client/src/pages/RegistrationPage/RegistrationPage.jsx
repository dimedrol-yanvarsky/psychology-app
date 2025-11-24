import React, { useState, useEffect } from "react";
import axios from "axios";
import styles from "./RegistrationPage.module.css";
import { Link, useNavigate } from "react-router-dom";
import { generateStrongPassword } from "./passwordGenerator";

const RegistrationPage = () => {
    const [login, setLogin] = useState("");
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
        try {
            const response = await axios.post(
                "http://localhost:8080/api/login",
                {
                    login,
                    password,
                }
            );

            if (response.data.success) {
                localStorage.setItem("token", response.data.token);
                window.location.href = "/dashboard";
            }
            setMessage(response.data.message);
        } catch (error) {
            setMessage("Ошибка авторизации");
        }
    };

    const handleOAuth = (provider) => {
        window.location.href = `http://localhost:8080/auth/${provider}`;
    };

    return (
        <div className={styles.login_form}>
            <form onSubmit={handleSubmit}>
                <h3>Регистрация</h3>
                <div className={styles.input_box}>
                    <label for="email">Как Вас зовут?</label>
                    <input
                        type="email"
                        id="email"
                        value={login}
                        onChange={(e) => setLogin(e.target.value)}
                        placeholder="Введите имя"
                        required
                    />
                </div>
                <div className={styles.input_box}>
                    <label for="email">Email</label>
                    <input
                        type="email"
                        id="email"
                        value={login}
                        onChange={(e) => setLogin(e.target.value)}
                        placeholder="Введите почтовый адрес"
                        required
                    />
                </div>
                <div className={styles.input_box_passw}>
                    <div className={styles.password_title}>
                        <label for="password">Придумайте пароль</label>
                        <a
                            onClick={handleGeneratePassword}
                            title="Генерирует сложный пароль"
                        >
                            Сгенерировать
                        </a>
                    </div>
                    <input
                        type={showPassword ? "text" : "password"}
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="Введите пароль"
                        required
                    />
                </div>
                <div className={styles.input_box_passw}>
                    <input
                        type={showPassword ? "text" : "password"}
                        id="password"
                        value={passwordRepeated}
                        onChange={(e) => setPasswordRepeated(e.target.value)}
                        placeholder="Повторите пароль"
                        required
                    />
                    {/* <input
                        type="checkbox"
                        checked={showPassword}
                        onChange={() => setShowPassword(!showPassword)}
                        className="checkbox-input"
                    /> */}
                </div>
                <button type="submit">Зарегистрироваться</button>
                <p className={styles.sign_up}>
                    Есть аккаунт?{" "}
                    <Link to="/login">Авторизуйтесь</Link>
                </p>
            </form>
        </div>
    );
};

export default RegistrationPage;
