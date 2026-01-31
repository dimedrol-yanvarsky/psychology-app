import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { generateStrongPassword } from "../lib/passwordGenerator";
import { createAccount } from "../api/authApi";

export const useRegistration = ({ showAlert }) => {
    const [login, setLogin] = useState("");
    const [name, setName] = useState("");
    const [password, setPassword] = useState("");
    const [passwordRepeated, setPasswordRepeated] = useState("");

    const navigate = useNavigate();

    const handleGeneratePassword = () => {
        const newPassword = generateStrongPassword();
        setPassword(newPassword);
        setPasswordRepeated(newPassword);
    };

    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const token = urlParams.get("token");
        if (token) {
            localStorage.setItem("token", token);
            window.location.href = "/dashboard";
        }
    }, []);

    const handleSubmit = async (event) => {
        event.preventDefault();

        if (!(name && login && password && passwordRepeated)) {
            return showAlert("error", "Не оставляйте поля пустыми");
        }

        if (password !== passwordRepeated) {
            return showAlert("error", "Пароли отличаются");
        }

        try {
            const { data } = await createAccount({
                name,
                login,
                password,
                passwordRepeated,
            });
            if (data?.success) {
                showAlert("success", "Аккаунт зарегистрирован");
                navigate("/login");
                return;
            }
            showAlert("error", data?.error || "Не удалось создать аккаунт");
        } catch (error) {
            const status = error?.response?.status || error?.status;
            if (status === 409) {
                return showAlert("error", "Аккаунт уже существует");
            }
            if (status === 403) {
                return showAlert("error", "Доступ запрещен");
            }
            if (status === 401) {
                return showAlert("error", "Проверьте введенные данные");
            }
            if (status === 500) {
                return showAlert("error", "База данных недоступна");
            }
            showAlert("error", "Не удалось создать аккаунт");
        }
    };

    return {
        handleGeneratePassword,
        handleSubmit,
        login,
        name,
        password,
        passwordRepeated,
        setLogin,
        setName,
        setPassword,
        setPasswordRepeated,
    };
};
