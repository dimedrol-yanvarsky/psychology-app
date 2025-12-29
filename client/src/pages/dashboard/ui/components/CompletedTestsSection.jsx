import React from "react";
import { Button } from "../../../../shared/ui/button";
import styles from "../DashboardPage.module.css";

const CompletedTestsSection = ({
    completedTests,
    isLoading,
    error,
    onOpenAnswers,
}) => {
    const renderContent = () => {
        if (isLoading) {
            return (
                <div className={styles.testsLoader}>
                    <div className={styles.loaderSpinner} aria-hidden="true" />
                    <span className={styles.loaderText}>Загрузка...</span>
                </div>
            );
        }

        if (error) {
            return <div className={styles.testsError}>{error}</div>;
        }

        if (completedTests.length === 0) {
            return (
                <div className={styles.psychotypeEmpty}>
                    <div>
                        <h4 className={styles.emptyTitle}>
                            Пока нет пройденных тестов
                        </h4>
                        <p className={styles.cardSubtitle}>
                            Пройдите тестирование, чтобы мы могли показать
                            историю результатов.
                        </p>
                    </div>
                    <Button
                        as="link"
                        to="/tests"
                        className={styles.primaryButtonLink}
                    >
                        Перейти к тестам
                    </Button>
                </div>
            );
        }

        return (
            <div className={styles.testsList}>
                {completedTests.map((test) => (
                    <div
                        key={test.testId || test.id}
                        className={styles.testRow}
                    >
                        <div className={styles.testInfo}>
                            <div className={styles.testTitle}>{test.testName}</div>
                            <div className={styles.testMeta}>
                                {test.date || "Дата не указана"} · {" "}
                                {test.result || "Без результата"}
                            </div>
                        </div>
                        <Button
                            type="button"
                            className={styles.linkButton}
                            onClick={() => onOpenAnswers(test)}
                        >
                            Открыть
                        </Button>
                    </div>
                ))}
            </div>
        );
    };

    return (
        <section className={styles.psychotypeCard}>
            <div className={styles.sectionHead}>
                <div>
                    <span className={styles.badge}>Мои тестирования</span>
                    <h3 className={styles.cardTitle}>Пройденные тесты</h3>
                    <p className={styles.cardSubtitle}>
                        Отслеживайте завершенные тестирования и возвращайтесь к
                        результатам при необходимости.
                    </p>
                </div>
            </div>

            {renderContent()}
        </section>
    );
};

export default CompletedTestsSection;
