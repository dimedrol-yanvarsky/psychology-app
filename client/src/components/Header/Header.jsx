import React from "react";
import { Link, useNavigate, useLocation } from "react-router-dom";
import styles from "./Header.module.css";
import logo from "../../pictures/logo.png"; // Предполагается, что логотип лежит в папке assets

const Header = () => {
    const navigate = useNavigate();
    const location = useLocation();
    // const isLoggedIn = localStorage.getItem('token');
    const isLoggedIn = true;

    const handleLogout = () => {
        localStorage.removeItem("token");
        navigate("/login");
    };

    const isActive = (path) => {
        return location.pathname === path ? styles.active : "";
    };

    return (
        <header className={styles.header}>
            <div className={styles.logoContainer}>
                <Link to="/">
                    <img src={logo} alt="Логотип" className={styles.logo} />
                </Link>
            </div>

            <nav className={styles.nav}>
                {isLoggedIn ? (
                    <>
                        <Link
                            to="/recommendations"
                            className={`${styles.navLink} ${isActive(
                                "/recommendations"
                            )}`}
                        >
                            Рекомендации
                        </Link>
                        <Link
                            to="/tests"
                            className={`${styles.navLink} ${isActive(
                                "/tests"
                            )}`}
                        >
                            Тестирования
                        </Link>
                        <Link
                            to="/reviews"
                            className={`${styles.navLink} ${isActive(
                                "/reviews"
                            )}`}
                        >
                            Отзывы
                        </Link>
                        <Link
                            to="/dashboard"
                            className={`${styles.navLink} ${
                                styles.homeLink
                            } ${isActive("/dashboard")}`}
                            title="Личный кабинет"
                        >
                            🏠
                        </Link>
                        {/* <button 
                            onClick={handleLogout} 
                            className={styles.logoutButton}
                        >
                            Выйти
                        </button> */}
                    </>
                ) : (
                    <Link to="/login" className={styles.navLink}>
                        Войти
                    </Link>
                )}
            </nav>
        </header>
    );
};

export default Header;
