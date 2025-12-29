import React from "react";
import { Modal } from "../../../../shared/ui/modal";
import AnswersQuestion from "./AnswersQuestion";
import styles from "../DashboardPage.module.css";

const AnswersModal = ({ answersModal, selectedAnswers, onClose }) => {
    const renderContent = () => {
        if (answersModal.loading) {
            return (
                <div className={styles.testsLoader}>
                    <div className={styles.loaderSpinner} aria-hidden="true" />
                    <span className={styles.loaderText}>Загрузка...</span>
                </div>
            );
        }

        if (answersModal.error) {
            return <div className={styles.testsError}>{answersModal.error}</div>;
        }

        return (
            <div className={styles.questionsList}>
                {answersModal.questions.map((question, index) => (
                    <AnswersQuestion
                        key={question.id || index}
                        question={question}
                        questionIndex={index}
                        selectedAnswers={selectedAnswers}
                    />
                ))}
            </div>
        );
    };

    return (
        <Modal
            isOpen={answersModal.open}
            onClose={onClose}
            overlayClassName={styles.modalOverlay}
            contentClassName={styles.answersModal}
            closeButtonClassName={styles.modalClose}
            closeLabel="Закрыть модальное окно"
        >
            <p className={styles.modalOverline}>Пройденный тест</p>
            <h4 className={styles.modalTitle}>
                {answersModal.title || "Результаты теста"}
            </h4>
            {renderContent()}
        </Modal>
    );
};

export default AnswersModal;
