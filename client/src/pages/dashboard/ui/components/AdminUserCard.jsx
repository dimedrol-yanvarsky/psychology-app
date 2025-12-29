import React from "react";
import { Button } from "../../../../shared/ui/button";
import styles from "../DashboardPage.module.css";

const AdminUserCard = ({
    account,
    isBlocking,
    isDeleting,
    onAdminAction,
    onBlockUser,
    onDeleteUser,
}) => {
    const status = (account.status || "").trim();
    const isBlocked = status === "Заблокирован";
    const isDeleted = status === "Удален";
    const blockButtonClass = `${styles.secondaryButton} ${
        isBlocked || isBlocking || isDeleted || isDeleting
            ? styles.blockedButton
            : ""
    }`;
    const deleteButtonClass = `${styles.cardActionButton} ${
        isDeleted ? styles.deletedButton : styles.cardDeleteButton
    }`;

    return (
        <div className={styles.adminUserCard}>
            <div className={styles.adminUserInfo}>
                <div className={styles.adminUserName}>
                    {account.firstName || "Без имени"}
                    {account.lastName ? ` ${account.lastName}` : ""}
                </div>
                <div className={styles.adminUserEmail}>{account.email || "—"}</div>
            </div>
            <div className={styles.adminActions}>
                <Button
                    type="button"
                    className={`${styles.cardActionButton} ${styles.cardPrimaryButton}`}
                    onClick={() =>
                        onAdminAction("Просмотр тестирований", account)
                    }
                >
                    Просмотреть тестирования
                </Button>
                <Button
                    type="button"
                    className={`${styles.cardActionButton} ${styles.cardPrimaryButton}`}
                    onClick={() => onAdminAction("Дерево эмоций", account)}
                >
                    Дерево эмоций
                </Button>
                <Button
                    type="button"
                    className={blockButtonClass}
                    onClick={() => onBlockUser(account)}
                    disabled={isBlocked || isBlocking || isDeleted || isDeleting}
                >
                    {isDeleted
                        ? "Удален"
                        : isBlocked
                        ? "Заблокирован"
                        : isBlocking
                        ? "Блокируем..."
                        : "Заблокировать аккаунт"}
                </Button>
                <Button
                    type="button"
                    className={deleteButtonClass}
                    onClick={() => onDeleteUser(account)}
                    disabled={isDeleted || isDeleting}
                >
                    {isDeleted
                        ? "Удален"
                        : isDeleting
                        ? "Удаляем..."
                        : "Удалить аккаунт"}
                </Button>
            </div>
        </div>
    );
};

export default AdminUserCard;
