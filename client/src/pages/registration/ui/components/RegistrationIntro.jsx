import React from "react";
import styles from "../RegistrationPage.module.css";

const RegistrationIntro = () => {
    return (
        <section className={styles.introCard}>
            <p className={styles.overline}>Старт для новых пользователей</p>
            <h1 className={styles.title}>
                Создайте аккаунт и продолжайте работу над собой
            </h1>
            <p className={styles.subtitle}>
                Сохраняем подборки рекомендаций, результаты тестов и личные
                цели. Зарегистрируйтесь, чтобы ничего не потерять и
                возвращаться к прогрессу в один клик.
            </p>
            <div className={styles.pills}>
                <span className={styles.pill}>Личный кабинет</span>
                <span className={styles.pill}>Сохранение прогресса</span>
                <span className={styles.pill}>Поддержка 24/7</span>
            </div>
        </section>
    );
};

export default RegistrationIntro;
