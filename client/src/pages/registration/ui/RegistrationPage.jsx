import React from "react";
import styles from "./RegistrationPage.module.css";
import RegistrationCard from "./components/RegistrationCard";
import RegistrationIntro from "./components/RegistrationIntro";
import { useRegistration } from "../../../features/auth";

const RegistrationPage = ({ showAlert }) => {
    const {
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
    } = useRegistration({ showAlert });

    return (
        <div className={styles.page}>
            {/* Основной блок регистрации */}
            <div className={styles.layout}>
                <RegistrationIntro />
                <RegistrationCard
                    name={name}
                    login={login}
                    password={password}
                    passwordRepeated={passwordRepeated}
                    onNameChange={setName}
                    onLoginChange={setLogin}
                    onPasswordChange={setPassword}
                    onPasswordRepeatedChange={setPasswordRepeated}
                    onGeneratePassword={handleGeneratePassword}
                    onSubmit={handleSubmit}
                />
            </div>
        </div>
    );
};

export default RegistrationPage;
