import React from "react";
import { Button } from "../../../../shared/ui/button";
import styles from "../DashboardPage.module.css";

const PsychotypeSection = ({ profileData, onOpenTestModal }) => {
    const renderContent = () => {
        if (profileData.psychoType) {
            return (
                <div className={styles.psychotypeContent}>
                    <div className={styles.psychotypeTag}>
                        {profileData.psychoType}
                    </div>
                    <p className={styles.subtitle}>
                        Вы всегда можете обновить результат, пройдя тест
                        повторно.
                    </p>
                    <div className={styles.retakeRow}>
                        <Button
                            type="button"
                            className={styles.secondaryButton}
                            onClick={onOpenTestModal}
                        >
                            Пройти заново
                        </Button>
                        <Button
                            as="link"
                            to="/tests"
                            className={styles.linkButton}
                        >
                            Открыть все тесты
                        </Button>
                    </div>
                </div>
            );
        }

        return (
            <div className={styles.psychotypeEmpty}>
                <div>
                    <h4 className={styles.emptyTitle}>Нет сохраненного психотипа</h4>
                    <p className={styles.cardSubtitle}>
                        Пройдите короткое тестирование, чтобы получить
                        персональные рекомендации и дерево эмоций.
                    </p>
                </div>
                <Button
                    type="button"
                    className={styles.primaryButton}
                    onClick={onOpenTestModal}
                >
                    Пройти тестирование
                </Button>
            </div>
        );
    };

    return (
        <section className={styles.psychotypeCard}>
            <div className={styles.sectionHead}>
                <div>
                    <span className={styles.badge}>Психотип</span>
                    <h3 className={styles.cardTitle}>
                        Результаты тестирования на психотип
                    </h3>
                    <p className={styles.cardSubtitle}>
                        Мы покажем Ваш психотип или предложим пройти
                        тестирование, если Вы еще этого не сделали.
                    </p>
                </div>
            </div>

            {renderContent()}
        </section>
    );
};

export default PsychotypeSection;
