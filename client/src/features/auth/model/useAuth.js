import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { loginWithPassword, loginWithProvider } from "../api/authApi";
import { useAuthContext } from "../../../shared/context/AuthContext";
import { useAlertContext } from "../../../shared/context/AlertContext";

export const useAuth = () => {
    const { setIsAdmin, setIsAuth, setProfileData } = useAuthContext();
    const { showAlert } = useAlertContext();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [recoveryStep, setRecoveryStep] = useState(1);
    const [recoveryLogin, setRecoveryLogin] = useState("");
    const [recoveryError, setRecoveryError] = useState("");

    const navigate = useNavigate();

    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const oauthStatus = urlParams.get("oauth");
        const error = urlParams.get("error");

        // Обработка OAuth callback
        if (oauthStatus === "success") {
            // Получаем данные пользователя из URL параметров
            const userData = {
                id: urlParams.get("id"),
                firstName: urlParams.get("firstName"),
                email: urlParams.get("email"),
                status: urlParams.get("status"),
                psychoType: urlParams.get("psychoType"),
                date: urlParams.get("date"),
                isGoogleAdded: urlParams.get("isGoogleAdded") === "true",
                isYandexAdded: urlParams.get("isYandexAdded") === "true",
            };

            // Устанавливаем данные пользователя
            setProfileData(userData);
            setIsAuth(true);
            setIsAdmin(userData.status === "Администратор");

            // Очищаем URL от параметров
            window.history.replaceState({}, document.title, "/account");

            showAlert("success", "Авторизация успешна");
            navigate("/account");
        }

        if (oauthStatus === "error" && error) {
            showAlert("error", `Ошибка OAuth: ${decodeURIComponent(error)}`);
            // Очищаем URL от параметров ошибки
            window.history.replaceState({}, document.title, "/login");
        }
    }, [navigate, setIsAdmin, setIsAuth, setProfileData, showAlert]);

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

            // Получаем URL для редиректа от сервера
            if (data.redirectUrl) {
                // Редиректим на OAuth провайдер (Google или Яндекс)
                window.location.href = data.redirectUrl;
            } else {
                showAlert("error", "Не получен URL для авторизации");
            }
        } catch (error) {
            const errorMessage = error?.response?.data?.error;

            if (errorMessage) {
                showAlert("error", errorMessage);
            } else if (error?.response?.status === 503) {
                showAlert(
                    "error",
                    `Авторизация через ${provider} не настроена на сервере`
                );
            } else {
                showAlert(
                    "error",
                    `Авторизация через ${provider} временно невозможна`
                );
            }
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
