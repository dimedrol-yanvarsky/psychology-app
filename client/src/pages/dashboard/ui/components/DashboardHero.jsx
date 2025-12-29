import React from "react";
import styles from "../DashboardPage.module.css";

const DashboardHero = () => {
    return (
        <section className={styles.heroCard}>
            <p className={styles.overline}>Личный кабинет</p>
            <h1 className={styles.title}>
                Управляйте профилем, тестированиями и безопасностью в одном
                месте
            </h1>
            <p className={styles.subtitle}>
                Редактируйте основные данные, следите за статусом и запускайте
                тесты на психотип. Все действия - централизовано и в едином
                стиле.
            </p>
            <div className={styles.pills}>
                <span className={styles.pill}>Шифрование данных с https</span>
                <span className={styles.pill}>Хеширование паролей с bcrypt</span>
                <span className={styles.pill}>Поддержка CI/CD</span>
            </div>
        </section>
    );
};

export default DashboardHero;
