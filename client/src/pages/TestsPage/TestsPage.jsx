import React from 'react';
import styles from './TestsPage.module.css';

const TestsPage = () => {
    return (
        <div className={styles.pageContainer}>
            <div className={styles.content}>
                <h1>Тестирования</h1>
                <p>Доступные тесты и опросы</p>
                <div className={styles.grid}>
                    <div className={styles.testCard}>
                        <h3>Тест знаний</h3>
                        <p>Проверьте свои знания</p>
                        <button className={styles.startButton}>Начать</button>
                    </div>
                    <div className={styles.testCard}>
                        <h3>Опрос предпочтений</h3>
                        <p>Помогите улучшить рекомендации</p>
                        <button className={styles.startButton}>Начать</button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default TestsPage;