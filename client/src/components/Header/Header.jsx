import React from "react";
import { Link, useLocation } from "react-router-dom";
import styles from "./Header.module.css";
import logo from "../../pictures/logo.png"; // Предполагается, что логотип лежит в папке assets

const Header = ({ isAuth }) => {
    const location = useLocation();

    const isActive = (path) => {
        return location.pathname === path ? styles.active : "";
    };

    // const homeTarget = isAuth ? "/dashboard" : "/login";
    const navItems = [
        { to: "/tree", label: "Мои эмоции"},
        { to: "/recommendations", label: "Рекомендации" },
        { to: "/tests", label: "Тестирования" },
        { to: "/reviews", label: "Отзывы" },
        {
            to: "/account",
            label: "Личный кабинет",
            className: styles.homeLink,
            icon: (
                <svg
                    className={styles.homeIcon}
                    viewBox="0 0 24 24"
                    aria-hidden="true"
                >
                    <path d="M12 3.5 3 11h2.5v8h5v-5h3v5h5v-8H21z" />
                </svg>
            ),
        },
    ];

    return (
        <header className={styles.header}>
            <div className={styles.inner}>
                <div className={styles.logoContainer}>
                    <img src={logo} alt="Логотип" className={styles.logo} />
                </div>

                <nav className={styles.nav}>
                    {navItems.map((item) => (
                        <Link
                            key={item.to}
                            to={item.to}
                            className={`${styles.navLink} ${item.className || ""} ${isActive(
                                item.to
                            )}`}
                            title={item.icon ? item.label : undefined}
                        >
                            {item.icon || item.label}
                        </Link>
                    ))}
                </nav>
            </div>
        </header>
    );
};

export default Header;
