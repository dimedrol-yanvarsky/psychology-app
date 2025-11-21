import React from 'react';
import styles from './RecommendationsPage.module.css';

const RecommendationsPage = () => {
    return (
        <div className={styles.pageContainer}>
            <div className={styles.content}>
                <h1>Рекомендации</h1>
                <p>Здесь будут ваши персональные рекомендации</p>
                <div className={styles.card}>
                    <h3>Рекомендуемые товары</h3>
                    <ul>
                        <li>Товар 1</li>
                        <li>Товар 2</li>
                        <li>Товар 3</li>
                    </ul>
                </div>
            </div>
        </div>
    );
};

export default RecommendationsPage;