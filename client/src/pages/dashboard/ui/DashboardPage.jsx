import React from "react";
import logoYandex from "../../../shared/assets/images/yandex-logo.png";
import logoGoogle from "../../../shared/assets/images/google-logo.png";
import styles from "./DashboardPage.module.css";
import AdminPanel from "./components/AdminPanel";
import AnswersModal from "./components/AnswersModal";
import CompletedTestsSection from "./components/CompletedTestsSection";
import DashboardHero from "./components/DashboardHero";
import EmotionGraphSection from "./components/EmotionGraphSection";
import ProfileSection from "./components/ProfileSection";
import PsychotypeSection from "./components/PsychotypeSection";
import TerminalModal from "./components/TerminalModal";
import TestModal from "./components/TestModal";
import { useDashboardPage } from "../model/useDashboardPage";

const DashboardPage = ({
    showAlert,
    setIsAuth,
    setIsAdmin,
    isAdmin,
    profileData,
    setProfileData,
}) => {
    const {
        adminAccounts,
        adminListError,
        answersModal,
        blockingUsers,
        closeAnswersModal,
        closeTerminal,
        closeTestModal,
        completedTests,
        completedTestsError,
        deletingUsers,
        emotionData,
        handleAdminAction,
        handleBlockUser,
        handleChangePassword,
        handleDeleteAccount,
        handleDeleteUser,
        handleFieldChange,
        handleLinkProvider,
        handleLogout,
        handleProfileSave,
        handleStartTesting,
        handleToggleTerminal,
        hasGoogle,
        hasYandex,
        isAdmin: isAdminView,
        isAdminListLoading,
        isCompletedTestsLoading,
        isDeletingAccount,
        isSavingProfile,
        isTerminalOpen,
        isTestModalOpen,
        openAnswersModal,
        openTestModal,
        selectedAnswers,
        setTerminalOpen,
        showLinkButtons,
    } = useDashboardPage({
        showAlert,
        setIsAuth,
        setIsAdmin,
        isAdmin,
        profileData,
        setProfileData,
    });

    return (
        <div className={styles.page}>
            <div className={styles.layout}>
                {/* Основные секции личного кабинета */}
                <DashboardHero />

                <ProfileSection
                    profileData={profileData}
                    isSavingProfile={isSavingProfile}
                    isDeletingAccount={isDeletingAccount}
                    onSave={handleProfileSave}
                    onFieldChange={handleFieldChange}
                    onChangePassword={handleChangePassword}
                    onLogout={handleLogout}
                    onDeleteAccount={handleDeleteAccount}
                    showLinkButtons={showLinkButtons}
                    hasGoogle={hasGoogle}
                    hasYandex={hasYandex}
                    onLinkProvider={handleLinkProvider}
                    logoGoogle={logoGoogle}
                    logoYandex={logoYandex}
                />

                <PsychotypeSection
                    profileData={profileData}
                    onOpenTestModal={openTestModal}
                />

                <CompletedTestsSection
                    completedTests={completedTests}
                    isLoading={isCompletedTestsLoading}
                    error={completedTestsError}
                    onOpenAnswers={openAnswersModal}
                />

                <EmotionGraphSection
                    profileData={profileData}
                    emotionData={emotionData}
                    onOpenTestModal={openTestModal}
                />

                {/* Блоки только для администратора */}
                {isAdminView && (
                    <AdminPanel
                        isTerminalOpen={isTerminalOpen}
                        onToggleTerminal={handleToggleTerminal}
                        adminAccounts={adminAccounts}
                        isAdminListLoading={isAdminListLoading}
                        adminListError={adminListError}
                        blockingUsers={blockingUsers}
                        deletingUsers={deletingUsers}
                        onAdminAction={handleAdminAction}
                        onBlockUser={handleBlockUser}
                        onDeleteUser={handleDeleteUser}
                    />
                )}
            </div>

            {/* Модальные окна страницы */}
            <TerminalModal
                isOpen={isTerminalOpen}
                onClose={closeTerminal}
                profileData={profileData}
                setIsTerminalOpen={setTerminalOpen}
            />

            <TestModal
                isOpen={isTestModalOpen}
                onClose={closeTestModal}
                onStartTesting={handleStartTesting}
            />

            <AnswersModal
                answersModal={answersModal}
                selectedAnswers={selectedAnswers}
                onClose={closeAnswersModal}
            />
        </div>
    );
};

export default DashboardPage;
