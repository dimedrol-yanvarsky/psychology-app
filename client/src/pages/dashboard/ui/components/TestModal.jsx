import React from "react";
import { Button } from "../../../../shared/ui/button";
import { Modal } from "../../../../shared/ui/modal";
import styles from "../DashboardPage.module.css";

const TestModal = ({ isOpen, onClose, onStartTesting }) => {
    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            overlayClassName={styles.modalOverlay}
            contentClassName={styles.modal}
            closeButtonClassName={styles.modalClose}
            closeLabel="Закрыть модальное окно"
        >
            <p className={styles.modalOverline}>Психотип</p>
            <h4 className={styles.modalTitle}>Тестирование на психотип</h4>
            <p className={styles.modalText}>
                Краткий опрос помогает понять эмоциональные реакции и
                предпочтения. Результат влияет на рекомендации, дерево эмоций и
                подборку тестов.
            </p>
            <div className={styles.modalChips}>
                <span className={styles.pill}>15-20 минут</span>
                <span className={styles.pill}>30 вопросов</span>
                <span className={styles.pill}>Результат сразу</span>
            </div>
            <Button
                type="button"
                className={styles.primaryButton}
                onClick={onStartTesting}
            >
                Начать тестирование
            </Button>
        </Modal>
    );
};

export default TestModal;
