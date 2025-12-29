import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { loginWithPassword, loginWithProvider } from "../api/loginApi";

export const useLoginPage = ({
    showAlert,
    setIsAdmin,
    setIsAuth,
    setProfileData,
}) => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [recoveryStep, setRecoveryStep] = useState(1);
    const [recoveryLogin, setRecoveryLogin] = useState("");
    const [recoveryError, setRecoveryError] = useState("");

    const navigate = useNavigate();

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
            const response = await loginWithPassword({ email, password });

            if (response.data.status === "Администратор") {
                setIsAdmin(true);
            }

            setProfileData(response.data);

            setIsAuth(true);
            showAlert("success", "Авторизация успешна");
            navigate("/account");
        } catch (error) {
            const status = error?.response?.status || error?.status;
            if (status === 401) {
                return showAlert("error", "Пароль неверный");
            }
            if (status === 500) {
                return showAlert("error", "База данных не отвечает");
            }
            if (status === 404) {
                return showAlert("error", "Пользователь не найден");
            }
            if (status === 400) {
                return showAlert("error", "Введите корректные данные");
            }
            if (status === 403) {
                return showAlert("error", "Доступ заблокирован");
            }
            showAlert("error", "Не удалось выполнить авторизацию");
        }
    };

    const handleOAuth = async (provider) => {
        try {
            const { data } = await loginWithProvider(provider);

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

    return {
        email,
        handleOAuth,
        handleRecoverySubmit,
        handleSubmit,
        isModalOpen,
        openRecoveryModal,
        password,
        recoveryError,
        recoveryLogin,
        recoveryStep,
        setEmail,
        setPassword,
        setRecoveryLogin,
        closeRecoveryModal,
    };
};
