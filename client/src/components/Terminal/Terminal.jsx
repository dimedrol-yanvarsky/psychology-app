import React, { useState } from "react";
import styles from "./Terminal.module.css";

const themeOptions = [
    { id: "dark", name: "Темная" },
    { id: "light", name: "Светлая" },
];

const Terminal = () => {
    const [activeTheme, setActiveTheme] = useState("dark");

    const terminalLines = [
        "Добро пожаловать в терминал веб-приложения!",
        "Он признан упростить работу администратора по управлению веб-приложением!",
        <>
            Для получения перечня допустимых команд используйте команду{" "}
            <span className={styles.command}>help</span>
        </>,
    ];

    return (
        <div className={styles.app}>
            <div
                className={`${styles.window} ${
                    activeTheme === "light" ? styles.windowLight : ""
                }`}
            >
                {/* Верхняя полоска с табами и контролами окна */}
                <div className={styles.topBar}>
                    <div className={styles.topBarLeft}>
                        {/* <button
                            className={`${styles.pill} ${styles.pillMuted}`}
                        >
                            <span className={styles.pillIcon} />
                            <span className={styles.pillLabel}>Vaults</span>
                        </button>

                        <button
                            className={`${styles.pill} ${styles.pillMuted}`}
                        >
                            <span className={styles.pillIcon} />
                            <span className={styles.pillLabel}>SFTP</span>
                        </button> */}

                        <button
                            className={`${styles.pill} ${styles.pillActive}`}
                        >
                            <span className={styles.pillLabel}>Терминал </span>
                        </button>
                    </div>

                    <div className={styles.windowControls}>
                        <span className={styles.controlIcon} />
                        <span className={styles.controlIcon} />
                        <span
                            className={`${styles.controlIcon} ${styles.controlClose}`}
                        />
                    </div>
                </div>

                {/* Основное содержимое: терминал + правая панель */}
                <div className={styles.content}>
                    {/* Левая часть — терминал */}
                    <div className={styles.terminalColumn}>
                        <div className={styles.terminalToolbar}>
                            <div className={styles.connectionBadge}>
                                <span
                                    className={styles.connectionStatusDot}
                                ></span>
                                <span className={styles.connectionText}>
                                    Дмитрий Голубев
                                </span>
                            </div>

                            <span className={styles.toolbarChip}>HTTP</span>
                        </div>

                        <div className={styles.terminalBody}>
                            <div className={styles.terminalLines}>
                                {terminalLines.map((line, idx) => (
                                    <div
                                        key={idx}
                                        className={styles.terminalLine}
                                    >
                                        {line}
                                    </div>
                                ))}
                            </div>

                            <div className={styles.promptLine}>
                                <span>root#</span>
                                <span className={styles.cursor} />
                            </div>
                        </div>
                    </div>

                    {/* Правая часть — настройки темы */}
                    <aside className={styles.sidebarColumn}>
                        <div className={styles.sidebarTopIcons}>
                            <button
                                className={styles.roundIcon}
                                aria-label="Очистить терминал"
                                title="Очистить терминал"
                            >
                                <svg
                                    className={styles.roundIconIcon}
                                    viewBox="0 0 24 24"
                                    aria-hidden="true"
                                >
                                    <path
                                        d="M6.187 8h11.625l-.695 11.125A2 2 0 0 1 15.121 21H8.879a2 2 0 0 1-1.996-1.875zM19 5v2H5V5h3V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v1zm-9 0h4V4h-4z"
                                        fill="currentColor"
                                    />
                                </svg>
                            </button>
                            <button
                                className={styles.roundIcon}
                                aria-label="Выбрать тему"
                                title="Выбрать тему"
                            >
                                <svg
                                    className={styles.roundIconIcon}
                                    viewBox="0 0 24 24"
                                    aria-hidden="true"
                                >
                                    <path
                                        d="M12 2C6.49 2 2 6.49 2 12s4.49 10 10 10a2.5 2.5 0 0 0 2.5-2.5c0-.61-.23-1.2-.64-1.67a.53.53 0 0 1-.13-.33c0-.28.22-.5.5-.5H16c3.31 0 6-2.69 6-6c0-4.96-4.49-9-10-9m5.5 11c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5s1.5.67 1.5 1.5s-.67 1.5-1.5 1.5m-3-4c-.83 0-1.5-.67-1.5-1.5S13.67 6 14.5 6s1.5.67 1.5 1.5S15.33 9 14.5 9M5 11.5c0-.83.67-1.5 1.5-1.5s1.5.67 1.5 1.5S7.33 13 6.5 13S5 12.33 5 11.5m6-4c0 .83-.67 1.5-1.5 1.5S8 8.33 8 7.5S8.67 6 9.5 6s1.5.67 1.5 1.5"
                                        fill="currentColor"
                                    />
                                </svg>
                            </button>
                        </div>

                        <div className={styles.sidebarPanel}>
                            <div className={styles.sectionHeading}>
                                Цветовая тема
                            </div>

                            <div className={styles.themesList}>
                                {themeOptions.map((theme) => {
                                    const isActive = theme.id === activeTheme;
                                    return (
                                        <button
                                            key={theme.id}
                                            className={`${styles.themeCard} ${
                                                isActive
                                                    ? styles.themeCardActive
                                                    : styles.themeCardInactive
                                            }`}
                                            onClick={() =>
                                                setActiveTheme(theme.id)
                                            }
                                        >
                                            <div
                                                className={styles.themePreview}
                                            >
                                                <div
                                                    className={`${
                                                        styles.previewBar
                                                    } ${
                                                        theme.id === "light"
                                                            ? styles.previewBarLightPrimary
                                                            : ""
                                                    }`}
                                                />
                                                <div
                                                    className={`${
                                                        styles.previewBarShort
                                                    } ${
                                                        theme.id === "light"
                                                            ? styles.previewBarLightSecondary
                                                            : ""
                                                    }`}
                                                />
                                                <div
                                                    className={`${
                                                        styles.previewBarBottom
                                                    } ${
                                                        theme.id === "light"
                                                            ? styles.previewBarLightAccent
                                                            : ""
                                                    }`}
                                                />
                                            </div>

                                            <div className={styles.themeInfo}>
                                                <div
                                                    className={styles.themeName}
                                                >
                                                    {theme.name}
                                                </div>
                                            </div>
                                        </button>
                                    );
                                })}
                            </div>
                        </div>
                    </aside>
                </div>
            </div>
        </div>
    );
};

export default Terminal;
