import React from "react";
import { Button } from "../../../../shared/ui/button";
import { Modal } from "../../../../shared/ui/modal";
import styles from "../LoginPage.module.css";

const RecoveryModal = ({
    isOpen,
    onClose,
    recoveryStep,
    recoveryLogin,
    recoveryError,
    onRecoveryLoginChange,
    onSubmit,
}) => {
    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            overlayClassName={styles.modalOverlay}
            contentClassName={styles.modal}
            closeButtonClassName={styles.modalClose}
            closeLabel="Закрыть модальное окно"
        >
            <h4 className={styles.modalTitle}>Восстановление пароля</h4>
            <div className={styles.stepper}>
                <div className={styles.stepLine}>
                    <div
                        className={styles.stepLineFill}
                        style={{
                            width: recoveryStep === 2 ? "100%" : "50%",
                        }}
                    />
                </div>
                <div className={styles.stepDots}>
                    {[1, 2].map((step) => (
                        <div
                            key={step}
                            className={`${styles.stepDot} ${
                                recoveryStep >= step
                                    ? styles.stepDotActive
                                    : ""
                            }`}
                        >
                            <span>{step}</span>
                            <p className={styles.stepLabel}>
                                {step === 1 ? "Логин" : "Подтверждение"}
                            </p>
                        </div>
                    ))}
                </div>
            </div>

            <form className={styles.modalForm} onSubmit={onSubmit}>
                <label className={styles.label} htmlFor="recovery-login">
                    Введите логин или email
                </label>
                <input
                    id="recovery-login"
                    className={styles.input}
                    type="text"
                    value={recoveryLogin}
                    onChange={(event) =>
                        onRecoveryLoginChange(event.target.value)
                    }
                    placeholder="Ваш логин"
                />
                {recoveryError && (
                    <p className={styles.modalError}>{recoveryError}</p>
                )}
                <p className={styles.modalHint}>
                    {recoveryStep === 1
                        ? "Отправим ссылку на восстановление пароля на указанный адрес."
                        : "Проверьте почту и следуйте инструкции для обновления пароля."}
                </p>
                <Button type="submit" className={styles.primaryButton}>
                    Восстановить
                </Button>
            </form>
        </Modal>
    );
};

export default RecoveryModal;
