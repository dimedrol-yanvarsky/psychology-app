import React, { useState, useEffect } from "react";
import axios from "axios";
import styles from "./LoginPage.module.css";
import logoYandex from "../../pictures/yandex-logo.png";
import logoGoogle from "../../pictures/google-logo.png";
import { Link, useNavigate } from "react-router-dom";

const LoginPage = (props) => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const navigate = useNavigate();
    const loginApiBase = "http://localhost:8080/api/login";

    useEffect(() => {
        console.log("tut", props.showAlert)
        // Проверяем, есть ли токен в URL (после OAuth редиректа)
        const urlParams = new URLSearchParams(window.location.search);
        const token = urlParams.get("token");
        const error = urlParams.get("error");

        if (token) {
            localStorage.setItem("token", token);
            window.location.href = "/dashboard";
        }

        if (error) {
            props.showAlert("error", `Ошибка авторизации: ${error}`);
        }
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!email || !password) {
            props.showAlert("error", "Не оставляйте поля пустыми");
            console.log("hi")
            return;
        }

        try {
            const { data } = await axios.post(`${loginApiBase}/password`, {
                email,
                password,
            });

            if (data.error) {
                props.showAlert("error", data.error);
                return;
            }

            props.showAlert("success", data.message || "Авторизация успешна");
            navigate(data.redirect || "/personal");
        } catch (error) {
            props.showAlert("error", "Ошибка авторизации");
        }
    };

    const handleOAuth = async (provider) => {
        try {
            const { data } = await axios.post(`${loginApiBase}/${provider}`);

            if (data.error) {
                props.showAlert("error", data.error);
                return;
            }

            props.showAlert("success", data.message || "Авторизация выполнена");
            navigate(data.redirect || "/account");
        } catch (error) {
            props.showAlert("error", "Авторизация невозможна");
        }
    };

    return (
        <div className={styles.login_form}>
            <form onSubmit={handleSubmit} noValidate>
                <h3>Авторизация</h3>
                <div className={styles.login_option}>
                    <div className={styles.option}>
                        <Link onClick={() => handleOAuth("google")}>
                            <img src={logoGoogle} />
                            <span>Google</span>
                        </Link>
                    </div>
                    <div className={styles.option}>
                        <Link onClick={() => handleOAuth("yandex")}>
                            <img src={logoYandex} />
                            <span>Яндекс</span>
                        </Link>
                    </div>
                </div>
                <p className={styles.separator}>
                    <span>или</span>
                </p>
                <div className={styles.input_box}>
                    <label htmlFor="email">Email</label>
                    <input
                        type="email"
                        id="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="Введите почтовый адрес"
                    />
                </div>
                <div className={styles.input_box}>
                    <div className={styles.password_title}>
                        <label htmlFor="password">Пароль</label>
                        <Link to="#">Забыли пароль?</Link>
                    </div>
                    <input
                        type="password"
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="Введите пароль"
                    />
                </div>
                <button type="submit">Войти</button>
                <p className={styles.sign_up}>
                    Нет аккаунта? <Link to="/register">Зарегистрируйтесь</Link>
                </p>
            </form>
        </div>
    );
};

export default LoginPage;
