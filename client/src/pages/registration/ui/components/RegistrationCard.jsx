import React from "react";
import { Button } from "../../../../shared/ui/button";
import styles from "../RegistrationPage.module.css";

const RegistrationCard = ({
    name,
    login,
    password,
    passwordRepeated,
    onNameChange,
    onLoginChange,
    onPasswordChange,
    onPasswordRepeatedChange,
    onGeneratePassword,
    onSubmit,
}) => {
    return (
        <div className={styles.registerCard}>
            <div className={styles.registerHeader}>
                <span className={styles.badge}>Регистрация</span>
                <div>
                    <h3 className={styles.cardTitle}>
                        Создайте аккаунт за пару шагов
                    </h3>
                    <p className={styles.cardSubtitle}>
                        Заполните форму ниже. Можно сгенерировать надежный
                        пароль прямо здесь.
                    </p>
                </div>
            </div>

            <form onSubmit={onSubmit} className={styles.form} noValidate>
                <label className={styles.label} htmlFor="name">
                    Как Вас зовут?
                </label>
                <input
                    className={styles.input}
                    type="text"
                    id="name"
                    value={name}
                    onChange={(event) => onNameChange(event.target.value)}
                    placeholder="Введите имя"
                    required
                />

                <label className={styles.label} htmlFor="email">
                    Email
                </label>
                <input
                    className={styles.input}
                    type="email"
                    id="email"
                    value={login}
                    onChange={(event) => onLoginChange(event.target.value)}
                    placeholder="Введите почтовый адрес"
                    required
                />

                <div className={styles.passwordRow}>
                    <label className={styles.label} htmlFor="password">
                        Придумайте пароль
                    </label>
                    <Button
                        type="button"
                        className={styles.linkButton}
                        onClick={onGeneratePassword}
                        title="Генерирует сложный пароль"
                    >
                        Сгенерировать
                    </Button>
                </div>
                <input
                    className={styles.input}
                    type="password"
                    id="password"
                    value={password}
                    onChange={(event) => onPasswordChange(event.target.value)}
                    placeholder="Введите пароль"
                    required
                />
                <p className={styles.helperText}>
                    Используйте буквы, цифры и символы для надежности.
                </p>

                <label className={styles.label} htmlFor="password-repeat">
                    Повторите пароль
                </label>
                <input
                    className={styles.input}
                    type="password"
                    id="password-repeat"
                    value={passwordRepeated}
                    onChange={(event) =>
                        onPasswordRepeatedChange(event.target.value)
                    }
                    placeholder="Повторите пароль"
                    required
                />

                <Button type="submit" className={styles.primaryButton}>
                    Зарегистрироваться
                </Button>
            </form>

            <p className={styles.footerText}>
                Есть аккаунт?{" "}
                <Button as="link" to="/login" className={styles.linkButton}>
                    Авторизуйтесь
                </Button>
            </p>
        </div>
    );
};

export default RegistrationCard;
