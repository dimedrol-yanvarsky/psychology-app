import React, { useEffect, useMemo, useRef, useState } from "react";
import axios from "axios";
import { API_BASE_URL } from "../../../shared/config/api";
import styles from "./Terminal.module.css";

const themeOptions = [
    { id: "light", name: "Светлая" },
    { id: "dark", name: "Темная" },
];

const dashboardApiBase = `${API_BASE_URL}/api/dashboard`;

const Terminal = ({ profileData, setIsTerminalOpen }) => {
    const [activeTheme, setActiveTheme] = useState("light");
    const [commandInput, setCommandInput] = useState("");
    const [isSending, setIsSending] = useState(false);

    const initialLines = useMemo(
        () => [
            {
                id: "welcome-1",
                type: "text",
                content: "Добро пожаловать в терминал веб-приложения!",
            },
            {
                id: "welcome-2",
                type: "text",
                content: "Он признан упростить работу администратора по управлению веб-приложением!",
            },
            {
                id: "welcome-3",
                type: "rich",
                content: (
                    <>
                        Для получения перечня допустимых команд используйте команду{" "}
                        <span className={styles.command}>help</span>
                    </>
                ),
            },
        ],
        []
    );

    const [terminalLines, setTerminalLines] = useState(initialLines);
    const commandInputRef = useRef(null);
    const scrollAreaRef = useRef(null);
    const lineIdRef = useRef(initialLines.length);

    const displayName = useMemo(() => {
        const parts = [profileData?.firstName, profileData?.lastName].filter(
            Boolean
        );
        return parts.join(" ") || "Администратор";
    }, [profileData]);

    const focusPrompt = () => {
        commandInputRef.current?.focus();
    };

    useEffect(() => {
        focusPrompt();
    }, []);

    useEffect(() => {
        const scrollArea = scrollAreaRef.current;
        if (scrollArea) {
            scrollArea.scrollTop = scrollArea.scrollHeight;
        }
    }, [terminalLines]);

    const normalizeCommand = (value = "") =>
        value
            .split(/\s+/)
            .filter(Boolean)
            .join(" ")
            .trim();

    const addLines = (lines = []) => {
        setTerminalLines((prev) => [
            ...prev,
            ...lines.map((line) => ({
                ...line,
                id: line.id || `line-${(lineIdRef.current += 1)}`,
            })),
        ]);
    };

    const handleClearTerminal = () => {
        setTerminalLines([]);
        focusPrompt();
    };

    const handleSubmitCommand = async (event) => {
        event?.preventDefault();
        if (isSending) return;

        const normalizedCommand = normalizeCommand(commandInput);
        if (!normalizedCommand) return;

        addLines([{ type: "command", command: normalizedCommand }]);
        setCommandInput("");
        setIsSending(true);

        try {
            const { data } = await axios.post(`${dashboardApiBase}/terminal`, {
                command: commandInput,
            });

            const serverCommand =
                normalizeCommand(data?.command) || normalizedCommand;

            if (data?.status === "success") {
                if (Array.isArray(data?.commands)) {
                    addLines([
                        {
                            type: "table",
                            title: "Доступные команды",
                            rows: data.commands,
                        },
                    ]);
                } else if (data?.message) {
                    addLines([{ type: "text", content: data.message }]);
                }

                if (
                    !Array.isArray(data?.commands) &&
                    !data?.message &&
                    serverCommand
                ) {
                    addLines([
                        {
                            type: "text",
                            content: `Команда "${serverCommand}" выполнена`,
                        },
                    ]);
                }
            } else {
                addLines([
                    {
                        type: "error",
                        content:
                            data?.message || "Не удалось распознать команду",
                    },
                ]);
            }
        } catch (error) {
            const serverMessage = error?.response?.data?.message;
            addLines([
                {
                    type: "error",
                    content:
                        serverMessage ||
                        "Не удалось выполнить команду. Проверьте подключение к серверу.",
                },
            ]);
        } finally {
            setIsSending(false);
            focusPrompt();
        }
    };

    const renderLineContent = (line) => {
        if (line.type === "command") {
            return (
                <div className={styles.commandLine}>
                    <span className={styles.promptLabel}>root#</span>
                    <span className={`${styles.command} ${styles.commandToken}`}>
                        {line.command}
                    </span>
                </div>
            );
        }

        if (line.type === "error") {
            return <span className={styles.errorText}>{line.content}</span>;
        }

        if (line.type === "table" && Array.isArray(line.rows)) {
            return (
                <div className={styles.terminalTableWrapper}>
                    <div className={styles.tableTitle}>
                        {line.title || "Команды"}
                    </div>
                    <table className={styles.terminalTable}>
                        <thead>
                            <tr>
                                <th>Команда</th>
                                <th>Описание</th>
                            </tr>
                        </thead>
                        <tbody>
                            {line.rows.map((row, index) => (
                                <tr key={row?.name || index}>
                                    <td>{row?.name}</td>
                                    <td>{row?.description}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            );
        }

        return line.content;
    };

    return (
        <div
            className={`${styles.window} ${
                activeTheme === "light" ? styles.windowLight : ""
            }`}
        >
            <div className={styles.topBar}>
                <div className={styles.topBarLeft}>
                    <button className={`${styles.pill} ${styles.pillActive}`}>
                        <span className={styles.pillLabel}>Терминал </span>
                    </button>
                </div>

                <div className={styles.windowControls}>
                    <span
                        className={`${styles.controlIcon} ${styles.controlClose}`}
                        onClick={() => setIsTerminalOpen?.(false)}
                    />
                </div>
            </div>

            {/* Основное содержимое: терминал + правая панель */}
            <div className={styles.content}>
                {/* Левая часть — терминал */}
                <div className={styles.terminalColumn}>
                    <div className={styles.terminalToolbar}>
                        <div className={styles.connectionBadge}>
                            <span className={styles.connectionStatusDot}></span>
                            <span className={styles.connectionText}>
                                {displayName}
                            </span>
                        </div>

                        <span className={styles.toolbarChip}>HTTPS</span>
                    </div>

                    <div className={styles.terminalBody} onClick={focusPrompt}>
                        <div
                            className={`${styles.terminalLines} ${styles.terminalScrollArea}`}
                            ref={scrollAreaRef}
                        >
                            {terminalLines.map((line) => (
                                <div key={line.id} className={styles.terminalLine}>
                                    {renderLineContent(line)}
                                </div>
                            ))}
                        </div>

                        <form
                            className={styles.promptLine}
                            onSubmit={handleSubmitCommand}
                        >
                            <span className={styles.promptLabel}>root#</span>
                            <div className={styles.promptInputWrapper}>
                                <span
                                    className={`${styles.inputMirror} ${
                                        commandInput ? styles.command : ""
                                    }`}
                                >
                                    {commandInput || ""}
                                </span>
                                {!commandInput && <span className={styles.cursor} />}
                                <input
                                    ref={commandInputRef}
                                    type="text"
                                    className={styles.promptInput}
                                    value={commandInput}
                                    onChange={(event) =>
                                        setCommandInput(event.target.value)
                                    }
                                    aria-label="Введите команду"
                                    autoComplete="off"
                                    spellCheck="false"
                                />
                            </div>
                        </form>
                    </div>
                </div>

                {/* Правая часть — настройки темы */}
                <aside className={styles.sidebarColumn}>
                    <div className={styles.sidebarTopIcons}>
                        <button
                            className={`${styles.roundIcon} ${styles.paletteButton}`}
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
                        <button
                            className={styles.roundIcon}
                            aria-label="Очистить терминал"
                            title="Очистить терминал"
                            onClick={handleClearTerminal}
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
                                        onClick={() => setActiveTheme(theme.id)}
                                    >
                                        <div
                                            className={`${
                                                styles.themePreview
                                            } ${
                                                theme.id === "light"
                                                    ? styles.themePreviewLight
                                                    : ""
                                            }`}
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
                                            <div className={styles.themeName}>
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
    );
};

export default Terminal;
