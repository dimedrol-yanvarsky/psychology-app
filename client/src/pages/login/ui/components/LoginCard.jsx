import React from "react";
import { Button } from "../../../../shared/ui/button";
import logoYandex from "../../../../shared/assets/images/yandex-logo.png";
import logoGoogle from "../../../../shared/assets/images/google-logo.png";
import styles from "../LoginPage.module.css";

const LoginCard = ({
    email,
    password,
    onEmailChange,
    onPasswordChange,
    onSubmit,
    onOAuth,
    onOpenRecoveryModal,
}) => {
    return (
        <div className={styles.loginCard}>
            <div className={styles.loginHeader}>
                <span className={styles.badge}>Авторизация</span>
                <div>
                    <h3 className={styles.cardTitle}>
                        Войдите любым удобным способом
                    </h3>
                    <p className={styles.cardSubtitle}>
                        Используйте корпоративный Google или Яндекс, либо
                        продолжите по email.
                    </p>
                </div>
            </div>

            <div className={styles.oauthRow}>
                <Button
                    type="button"
                    className={styles.oauthButton}
                    onClick={() => onOAuth("google")}
                >
                    <img src={logoGoogle} alt="Google" />
                    <span>Войти через Google</span>
                </Button>
                <Button
                    type="button"
                    className={styles.oauthButton}
                    onClick={() => onOAuth("yandex")}
                >
                    <img src={logoYandex} alt="Яндекс" />
                    <span>Войти через Яндекс</span>
                </Button>
            </div>

            <div className={styles.divider}>
                <span>или авторизация по почте</span>
            </div>

            <form onSubmit={onSubmit} className={styles.form} noValidate>
                <label className={styles.label} htmlFor="email">
                    Email
                </label>
                <input
                    className={styles.input}
                    type="email"
                    id="email"
                    value={email}
                    onChange={(event) => onEmailChange(event.target.value)}
                    placeholder="Введите почтовый адрес"
                />

                <div className={styles.passwordRow}>
                    <label className={styles.label} htmlFor="password">
                        Пароль
                    </label>
                    <Button
                        as="link"
                        to="#"
                        className={styles.linkButton}
                        onClick={(event) => {
                            event.preventDefault();
                            onOpenRecoveryModal();
                        }}
                    >
                        Забыли пароль?
                    </Button>
                </div>
                <input
                    className={styles.input}
                    type="password"
                    id="password"
                    value={password}
                    onChange={(event) => onPasswordChange(event.target.value)}
                    placeholder="Введите пароль"
                />

                <Button type="submit" className={styles.primaryButton}>
                    Войти
                </Button>
            </form>

            <p className={styles.footerText}>
                Нет аккаунта?{" "}
                <Button as="link" to="/register" className={styles.linkButton}>
                    Зарегистрируйтесь
                </Button>
            </p>
        </div>
    );
};

export default LoginCard;
