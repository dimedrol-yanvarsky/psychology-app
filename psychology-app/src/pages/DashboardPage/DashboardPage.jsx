import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import styles from './DashboardPage.module.css';

const DashboardPage = () => {
    const [user, setUser] = useState(null);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            window.location.href = '/login';
            return;
        }

        axios.get('http://localhost:8080/auth/me', {
            headers: { Authorization: `Bearer ${token}` }
        })
        .then(response => {
            setUser(response.data.user);
        })
        .catch(() => {
            localStorage.removeItem('token');
            window.location.href = '/login';
        });
    }, []);

    if (!user) return <div className={styles.pageContainer}>Загрузка...</div>;

    return (
        <div className={styles.pageContainer}>
            <div className={styles.content}>
                <h1>Личный кабинет</h1>
                <div className={styles.card}>
                    <h2>Добро пожаловать, {user.login}!</h2>
                    <p><strong>Email:</strong> {user.email}</p>
                    <p><strong>Провайдер:</strong> {user.provider}</p>
                    <p><strong>Дата регистрации:</strong> {new Date(user.createdAt).toLocaleDateString()}</p>
                </div>
                
                <div className={styles.grid}>
                    <Link to="/recommendations" className={styles.dashboardCard}>
                        <h3>🎯 Рекомендации</h3>
                        <p>Персональные рекомендации на основе ваших предпочтений</p>
                    </Link>
                    
                    <Link to="/tests" className={styles.dashboardCard}>
                        <h3>📊 Тестирования</h3>
                        <p>Пройдите тесты для улучшения рекомендаций</p>
                    </Link>
                </div>
            </div>
        </div>
    );
};

export default DashboardPage;