import React from "react";
import { Button } from "../../../../shared/ui/button";
import styles from "../DashboardPage.module.css";

const ProfileSection = ({
    profileData,
    isSavingProfile,
    isDeletingAccount,
    onSave,
    onFieldChange,
    onChangePassword,
    onLogout,
    onDeleteAccount,
    showLinkButtons,
    hasGoogle,
    hasYandex,
    onLinkProvider,
    logoGoogle,
    logoYandex,
}) => {
    return (
        <section className={styles.profileCard}>
            <div className={styles.sectionHead}>
                <div>
                    <span className={styles.badge}>Профиль</span>
                    <h3 className={styles.cardTitle}>Персональные данные</h3>
                    <p className={styles.cardSubtitle}>
                        Обновляйте имя, а почта закреплена за аккаунтом. Мы
                        бережно храним изменения и используем их в
                        рекомендациях.
                    </p>
                </div>
            </div>

            <div className={styles.profileGrid}>
                <form className={styles.form} onSubmit={onSave}>
                    <label className={styles.label} htmlFor="first-name">
                        Имя
                    </label>
                    <input
                        id="first-name"
                        className={styles.input}
                        type="text"
                        value={profileData.firstName}
                        onChange={(event) =>
                            onFieldChange("firstName", event.target.value)
                        }
                        placeholder="Имя"
                    />

                    <label className={styles.label} htmlFor="email">
                        Почтовый адрес
                    </label>
                    <input
                        id="email"
                        className={styles.input}
                        type="email"
                        value={profileData.email}
                        placeholder="example@domain.com"
                        readOnly
                        aria-readonly="true"
                    />

                    <Button
                        type="submit"
                        className={styles.primaryButton}
                        disabled={isSavingProfile}
                    >
                        {isSavingProfile ? "Сохраняем..." : "Сохранить изменения"}
                    </Button>
                </form>

                <div className={styles.infoPanel}>
                    <div className={styles.infoRow}>
                        <span className={styles.infoLabel}>Статус</span>
                        <span className={styles.infoValue}>
                            {profileData.status}
                        </span>
                    </div>
                    <div className={styles.actionsRow}>
                        <Button
                            type="button"
                            className={`${styles.cardActionButton} ${styles.cardEditButton}`}
                            onClick={onChangePassword}
                        >
                            Сменить пароль
                        </Button>
                        <Button
                            type="button"
                            className={styles.secondaryButton}
                            onClick={onLogout}
                        >
                            Выйти из аккаунта
                        </Button>
                        <Button
                            type="button"
                            className={`${styles.cardActionButton} ${styles.cardDeleteButton}`}
                            onClick={onDeleteAccount}
                            disabled={isDeletingAccount}
                        >
                            {isDeletingAccount ? "Удаляем..." : "Удалить аккаунт"}
                        </Button>
                    </div>
                </div>

                {showLinkButtons && (
                    <div className={styles.infoPanel}>
                        <div className={styles.infoRow}>
                            <span className={styles.infoLabel}>
                                Привязать аккаунт
                            </span>
                        </div>
                        <div className={styles.connectRow}>
                            {!hasGoogle && (
                                <Button
                                    type="button"
                                    className={styles.oauthButton}
                                    onClick={() => onLinkProvider("Google")}
                                >
                                    <img src={logoGoogle} alt="Google" />
                                    <span>Привязать аккаунт Google</span>
                                </Button>
                            )}
                            {!hasYandex && (
                                <Button
                                    type="button"
                                    className={styles.oauthButton}
                                    onClick={() => onLinkProvider("Яндекс")}
                                >
                                    <img src={logoYandex} alt="Yandex" />
                                    <span>Привязать аккаунт Яндекс</span>
                                </Button>
                            )}
                        </div>
                    </div>
                )}
            </div>
        </section>
    );
};

export default ProfileSection;
