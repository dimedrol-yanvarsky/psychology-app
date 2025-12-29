import React from "react";
import styles from "../LoginPage.module.css";

const LoginIntro = () => {
    return (
        <section className={styles.introCard}>
            <p className={styles.overline}>Доступ к рекомендациям</p>
            <h1 className={styles.title}>
                Вернитесь к спокойной работе над собой
            </h1>
            <p className={styles.subtitle}>
                Сохраняем ваш прогресс, подборки и результаты тестов. Войдите,
                чтобы продолжить работу с рекомендациями и личным кабинетом.
            </p>
            <div className={styles.pills}>
                <span className={styles.pill}>Безопасный вход</span>
                <span className={styles.pill}>Поддержка 24/7</span>
                <span className={styles.pill}>Синхронизация данных</span>
            </div>
        </section>
    );
};

export default LoginIntro;
