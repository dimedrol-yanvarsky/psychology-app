import React from "react";
import styles from "./LoginPage.module.css";
import LoginCard from "./components/LoginCard";
import LoginIntro from "./components/LoginIntro";
import RecoveryModal from "./components/RecoveryModal";
import { useLoginPage } from "../model/useLoginPage";

const LoginPage = ({ showAlert, setIsAdmin, setIsAuth, setProfileData }) => {
    const {
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
    } = useLoginPage({ showAlert, setIsAdmin, setIsAuth, setProfileData });

    return (
        <div className={styles.page}>
            {/* Основной блок страницы входа */}
            <div className={styles.layout}>
                <LoginIntro />
                <LoginCard
                    email={email}
                    password={password}
                    onEmailChange={setEmail}
                    onPasswordChange={setPassword}
                    onSubmit={handleSubmit}
                    onOAuth={handleOAuth}
                    onOpenRecoveryModal={openRecoveryModal}
                />
            </div>

            {/* Модальное окно восстановления пароля */}
            <RecoveryModal
                isOpen={isModalOpen}
                onClose={closeRecoveryModal}
                recoveryStep={recoveryStep}
                recoveryLogin={recoveryLogin}
                recoveryError={recoveryError}
                onRecoveryLoginChange={setRecoveryLogin}
                onSubmit={handleRecoverySubmit}
            />
        </div>
    );
};

export default LoginPage;
