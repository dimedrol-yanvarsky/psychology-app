import React, { useState, useEffect } from "react";
import axios from "axios";
import styles from "./RegistrationPage.module.css";
import { useNavigate } from "react-router-dom";

const RegistrationPage = () => {
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");
    const [passwordRepeated, setPasswordRepeated] = useState("");
    const [message, setMessage] = useState("");
    const [showPassword, setShowPassword] = useState(false);

    const navigate = useNavigate();

    function generateStrongPassword() {
        // Определяем наборы символов
        const lowercase = "abcdefghijklmnopqrstuvwxyz";
        const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
        const numbers = "0123456789";
        const specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?";

        // Объединяем все символы
        const allChars = lowercase + uppercase + numbers + specialChars;

        // Генерируем случайную длину от 14 до 20 символов
        const length = Math.floor(Math.random() * 7) + 14;

        let passw = "";

        // Гарантируем наличие хотя бы одного символа из каждой категории
        passw += lowercase[Math.floor(Math.random() * lowercase.length)];
        passw += uppercase[Math.floor(Math.random() * uppercase.length)];
        passw += numbers[Math.floor(Math.random() * numbers.length)];
        passw += specialChars[Math.floor(Math.random() * specialChars.length)];

        // Заполняем оставшуюся длину случайными символами из всех категорий
        for (let i = 4; i < length; i++) {
            passw += allChars[Math.floor(Math.random() * allChars.length)];
        }

        const shuffledPassw = shuffleString(passw);
        setPassword(shuffledPassw);
        setPasswordRepeated(shuffledPassw);
        // Перемешиваем символы для большей случайности
        return true;
    }

    // Вспомогательная функция для перемешивания символов в строке
    function shuffleString(string) {
        const array = string.split("");
        for (let i = array.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [array[i], array[j]] = [array[j], array[i]];
        }
        return array.join("");
    }

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
                            onClick={() => generateStrongPassword()}
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
                    <input
                        type="checkbox"
                        checked={showPassword}
                        onChange={() => setShowPassword(!showPassword)}
                        className="checkbox-input"
                    />
                </div>
                <button type="submit">Войти</button>
                <p className={styles.sign_up}>
                    Есть аккаунт?{" "}
                    <a onClick={() => navigate(`/login`)}>Авторизуйтесь</a>
                </p>
            </form>
        </div>
    );
};

export default RegistrationPage;
