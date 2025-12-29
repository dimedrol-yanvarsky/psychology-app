import React from "react";
import { Button } from "../../../../shared/ui/button";
import AdminUserCard from "./AdminUserCard";
import styles from "../DashboardPage.module.css";

const AdminPanel = ({
    isTerminalOpen,
    onToggleTerminal,
    adminAccounts,
    isAdminListLoading,
    adminListError,
    blockingUsers,
    deletingUsers,
    onAdminAction,
    onBlockUser,
    onDeleteUser,
}) => {
    const renderListContent = () => {
        if (isAdminListLoading) {
            return (
                <div className={styles.adminLoader}>
                    <div className={styles.loaderSpinner} aria-hidden="true" />
                    <span className={styles.loaderText}>Загрузка данных...</span>
                </div>
            );
        }

        if (adminListError) {
            return <div className={styles.adminError}>{adminListError}</div>;
        }

        if (adminAccounts.length === 0) {
            return (
                <div className={styles.adminEmpty}>Пользователи не найдены</div>
            );
        }

        return adminAccounts.map((account) => (
            <AdminUserCard
                key={account.id}
                account={account}
                isBlocking={Boolean(blockingUsers[account.id])}
                isDeleting={Boolean(deletingUsers[account.id])}
                onAdminAction={onAdminAction}
                onBlockUser={onBlockUser}
                onDeleteUser={onDeleteUser}
            />
        ));
    };

    return (
        <section className={styles.adminPanel}>
            <div className={styles.sectionHead}>
                <div>
                    <span className={styles.badge}>Панель администратора</span>
                    <h3 className={styles.cardTitle}>Управление пользователями</h3>
                    <p className={styles.cardSubtitle}>
                        Просматривайте зарегистрированные аккаунты, блокируйте
                        или удаляйте доступ, открывайте тестирования и дерево
                        эмоций.
                    </p>
                </div>
                <Button
                    type="button"
                    className={`${styles.cardActionButton} ${styles.cardEditButton}`}
                    onClick={onToggleTerminal}
                >
                    {isTerminalOpen ? "Скрыть терминал" : "Открыть терминал"}
                </Button>
            </div>

            <div className={styles.adminList}>{renderListContent()}</div>
        </section>
    );
};

export default AdminPanel;
