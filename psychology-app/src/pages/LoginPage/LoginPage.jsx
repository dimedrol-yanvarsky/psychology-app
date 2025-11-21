import React, { useState, useEffect } from "react";
import axios from "axios";
import styles from "./LoginPage.module.css";
import logoYandex from "../../pictures/yandex-logo.png";
import logoGoogle from "../../pictures/google-logo.png";
import { useNavigate } from "react-router-dom";

const LoginPage = () => {
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");
    const [message, setMessage] = useState("");

    const navigate = useNavigate();

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
                "http://localhost:8080/auth/login",
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
                <h3>Авторизация</h3>
                <div className={styles.login_option}>
                    <div className={styles.option}>
                        <a onClick={() => handleOAuth("google")}>
                            <img src={logoGoogle} />
                            <span>Google</span>
                        </a>
                    </div>
                    <div className={styles.option}>
                        <a onClick={() => handleOAuth("yandex")}>
                            <img src={logoYandex} />
                            <span>Яндекс</span>
                        </a>
                    </div>
                </div>
                <p className={styles.separator}>
                    <span>или</span>
                </p>
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
                <div className={styles.input_box}>
                    <div className={styles.password_title}>
                        <label for="password">Пароль</label>
                        <a href="#">Забыли пароль?</a>
                    </div>
                    <input
                        type="password"
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="Введите пароль"
                        required
                    />
                </div>
                <button type="submit">Войти</button>
                <p className={styles.sign_up}>
                    Нет аккаунта?{" "}
                    <a onClick={() => navigate(`/register`)}>
                        Зарегистрируйтесь
                    </a>
                </p>
            </form>
        </div>
    );
};

export default LoginPage;
