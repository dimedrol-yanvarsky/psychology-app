import React, { useMemo, useState } from "react";
import styles from "./TestsPage.module.css";

const baseTests = [
    {
        id: "t-1",
        title: "Шкала депрессии Бека",
        duration: "10–15 минут",
        author: "Аарон Бек",
        summary:
            "Самостоятельная оценка выраженности депрессивной симптоматики. Подходит для первичного скрининга состояния.",
    },
    {
        id: "t-2",
        title: "Опросник Спилбергера—Ханина",
        duration: "7–10 минут",
        author: "Ч. Д. Спилбергер, Ю. Л. Ханин",
        summary:
            "Измеряет ситуативную и личностную тревожность, помогает отследить динамику в ходе терапии.",
    },
    {
        id: "t-3",
        title: "Шкала позитивного и негативного аффекта (PANAS)",
        duration: "5–7 минут",
        author: "Д. Уотсон, Л. Кларк, Э. Теллеген",
        summary:
            "Определяет баланс позитивных и негативных эмоций, даёт быстрый срез эмоционального состояния.",
    },
];

const TestsPage = (props) => {
    const [tests, setTests] = useState(baseTests);
    const emptyState = useMemo(() => tests.length === 0, [tests]);

    const handleDelete = (id) => {
        setTests((prev) => prev.filter((item) => item.id !== id));
    };

    return (
        <div className={styles.page}>
            <header className={styles.topBar}>
                <h1 className={styles.title}>Психологическое тестирование</h1>
                {props.isAdmin === true ? (
                    <button
                        type="button"
                        className={`${styles.actionButton} ${styles.primaryAction}`}
                    >
                        Добавить тестирование
                    </button>
                ) : (
                    <></>
                )}
            </header>

            <main className={styles.main}>
                {emptyState ? (
                    <div className={styles.emptyStateBanner}>
                        Пока нет доступных тестов. Добавьте первый тест.
                    </div>
                ) : (
                    <div className={styles.testsList}>
                        {tests.map((test) => (
                            <article key={test.id} className={styles.testCard}>
                                <header className={styles.cardHeader}>
                                    <div className={styles.cardTitleBlock}>
                                        <h2 className={styles.testTitle}>
                                            {test.title}
                                        </h2>
                                        <span className={styles.testAuthor}>
                                            Автор: {test.author}
                                        </span>
                                    </div>
                                    <span className={styles.durationBadge}>
                                        {test.duration}
                                    </span>
                                </header>

                                <p className={styles.testSummary}>
                                    {test.summary}
                                </p>

                                <div className={styles.cardActions}>
                                    <button
                                        type="button"
                                        className={`${styles.cardButton} ${styles.primaryButton}`}
                                    >
                                        Пройти тест
                                    </button>
                                    {props.isAdmin === true ? (
                                        <div className={styles.adminActions}>
                                            <button
                                                type="button"
                                                className={`${styles.cardButton} ${styles.editButton}`}
                                            >
                                                Изменить тестирование
                                            </button>
                                            <button
                                                type="button"
                                                className={`${styles.cardButton} ${styles.deleteButton}`}
                                                onClick={() =>
                                                    handleDelete(test.id)
                                                }
                                            >
                                                Удалить тестирование
                                            </button>
                                        </div>
                                    ) : (
                                        <></>
                                    )}
                                </div>
                            </article>
                        ))}
                    </div>
                )}
            </main>
        </div>
    );
};

export default TestsPage;
