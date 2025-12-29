import React from "react";
import styles from "./TreePage.module.css";
import underConstruction from "../../../shared/assets/images/under-construction.svg";

const TreePage = () => {
    return (
        <div className={styles.page}>
            {/* Заглушка для страницы в разработке */}
            <div className={styles.card}>
                <p className={styles.badge}>Скоро доступно</p>
                <h1 className={styles.title}>Страница находится в разработке</h1>
                <p className={styles.subtitle}>
                    Мы готовим визуализацию и трекер ваших эмоций. Загляните позже
                    - скоро все будет готово.
                </p>
                <div className={styles.imageWrapper}>
                    <img
                        src={underConstruction}
                        alt="Иконка предупреждения — страница в разработке"
                        className={styles.image}
                    />
                </div>
            </div>
        </div>
    );
};

export default TreePage;
