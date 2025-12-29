import React from "react";
import { Button } from "../../../../shared/ui/button";
import styles from "../DashboardPage.module.css";

const EmotionGraphSection = ({ profileData, emotionData, onOpenTestModal }) => {
    const renderContent = () => {
        if (!profileData.psychotype) {
            return (
                <div className={styles.psychotypeEmpty}>
                    <div>
                        <h4 className={styles.emptyTitle}>
                            Нет данных для построения графика
                        </h4>
                        <p className={styles.cardSubtitle}>
                            Пройдите тест на психотип, чтобы увидеть динамику
                            эмоций.
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
        }

        return (
            <div className={styles.emotionGraph}>
                {emotionData.map((item) => (
                    <div key={item.id} className={styles.emotionRow}>
                        <span className={styles.emotionLabel}>{item.label}</span>
                        <div className={styles.emotionBarTrack}>
                            <div
                                className={styles.emotionBar}
                                style={{ width: `${item.value}%` }}
                            />
                        </div>
                        <span className={styles.emotionValue}>{item.value}%</span>
                    </div>
                ))}
            </div>
        );
    };

    return (
        <section className={styles.psychotypeCard}>
            <div className={styles.sectionHead}>
                <div>
                    <span className={styles.badge}>Мои эмоции</span>
                    <h3 className={styles.cardTitle}>
                        Граф эмоционального состояния
                    </h3>
                    <p className={styles.cardSubtitle}>
                        Сводка по тестированиям помогает понять динамику вашего
                        состояния.
                    </p>
                </div>
            </div>

            {renderContent()}
        </section>
    );
};

export default EmotionGraphSection;
