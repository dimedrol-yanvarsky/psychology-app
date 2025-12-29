import React from "react";
import styles from "./RecommendationsPage.module.css";
import { useRecommendationsPage } from "../model/useRecommendationsPage";

const RecommendationsPage = ({ isAdmin, showAlert }) => {
    const {
        TEXT_MODE_OPTIONS,
        blockStyle,
        currentPage,
        currentSection,
        editingBlock,
        getPageNumber,
        handleAddBlock,
        handleAddSection,
        handleCancelEdit,
        handleDeleteBlock,
        handleDeleteSection,
        handlePageChange,
        handleSaveBlock,
        handleStartEdit,
        isBlockPending,
        isLoading,
        pending,
        sections,
        setEditingBlock,
    } = useRecommendationsPage({ isAdmin, showAlert });

    const renderBlocks = () => {
        // Рендер списка блоков с учетом состояния загрузки и роли.
        if (isLoading) {
            return (
                <div className={styles.emptyState}>
                    Загружаем рекомендации...
                </div>
            );
        }

        if (!currentSection) {
            return (
                <div className={styles.emptyState}>
                    <p>Рекомендаций пока нет.</p>
                    {isAdmin && (
                        <button
                            type="button"
                            className={`${styles.actionButton} ${styles.secondaryAction}`}
                            onClick={handleAddSection}
                            disabled={pending.addSection}
                        >
                            Добавить раздел
                        </button>
                    )}
                </div>
            );
        }

        return (
            <div className={styles.blocksGrid}>
                {currentSection.blocks.map((block) => {
                    const isEditing =
                        editingBlock?.id === block.id &&
                        editingBlock?.sectionType ===
                            currentSection.recommendationType;

                    return (
                        <article
                            key={block.id || block.recommendationText}
                            className={`${styles.blockCard} ${
                                isEditing ? styles.blockEditing : ""
                            }`}
                        >
                            {!isEditing ? (
                                <>
                                    <p
                                        className={styles.blockText}
                                        style={blockStyle(block.textMode)}
                                    >
                                        {block.recommendationText}
                                    </p>
                                    {isAdmin && (
                                        <div className={styles.blockActions}>
                                            <button
                                                type="button"
                                                className={`${styles.inlineButton} ${styles.secondaryAction}`}
                                                onClick={() =>
                                                    handleStartEdit(
                                                        block,
                                                        currentSection.recommendationType
                                                    )
                                                }
                                                disabled={isBlockPending(
                                                    block.id
                                                )}
                                            >
                                                Редактировать
                                            </button>
                                            <button
                                                type="button"
                                                className={`${styles.inlineButton} ${styles.dangerAction}`}
                                                onClick={() =>
                                                    handleDeleteBlock(
                                                        block,
                                                        currentSection.recommendationType
                                                    )
                                                }
                                                disabled={isBlockPending(
                                                    block.id
                                                )}
                                            >
                                                Удалить
                                            </button>
                                        </div>
                                    )}
                                </>
                            ) : (
                                <div className={styles.editorWrapper}>
                                    <div className={styles.modePicker}>
                                        {TEXT_MODE_OPTIONS.map((option) => (
                                            <button
                                                type="button"
                                                key={option.value}
                                                className={`${styles.toolButton} ${
                                                    editingBlock.draftMode ===
                                                    option.value
                                                        ? styles.toolActive
                                                        : ""
                                                }`}
                                                onClick={() =>
                                                    setEditingBlock((prev) => ({
                                                        ...prev,
                                                        draftMode:
                                                            option.value,
                                                    }))
                                                }
                                                title={option.title}
                                                aria-label={option.title}
                                            >
                                                {option.label}
                                            </button>
                                        ))}
                                    </div>
                                    <textarea
                                        className={styles.editorInput}
                                        value={editingBlock.draftText}
                                        onChange={(e) =>
                                            setEditingBlock((prev) => ({
                                                ...prev,
                                                draftText: e.target.value,
                                            }))
                                        }
                                    />
                                    <div className={styles.editorActions}>
                                        <button
                                            type="button"
                                            className={`${styles.actionButton} ${styles.secondaryAction}`}
                                            onClick={handleSaveBlock}
                                            disabled={isBlockPending(
                                                block.id
                                            )}
                                        >
                                            Сохранить
                                        </button>
                                        <button
                                            type="button"
                                            className={`${styles.actionButton} ${styles.mutedAction}`}
                                            onClick={handleCancelEdit}
                                        >
                                            Отмена
                                        </button>
                                    </div>
                                </div>
                            )}
                        </article>
                    );
                })}
            </div>
        );
    };

    return (
        <div className={styles.page}>
            {/* Верхняя панель управления разделами */}
            <header className={styles.topBar}>
                <div className={styles.breadcrumb}>
                    <span className={styles.overline}>Recommendations</span>
                    <h1 className={styles.pageTitle}>Рекомендации</h1>
                    <p className={styles.pageSubtitle}>
                        Блоки рекомендаций с возможностью редактирования,
                        выбора стиля и управления разделами.
                    </p>
                </div>

                {isAdmin && (
                    <div className={styles.topBarActions}>
                        <button
                            type="button"
                            className={`${styles.actionButton} ${styles.secondaryAction}`}
                            onClick={handleAddSection}
                            disabled={pending.addSection}
                        >
                            Добавить раздел
                        </button>
                        <button
                            type="button"
                            className={`${styles.actionButton} ${styles.secondaryAction}`}
                            onClick={handleAddBlock}
                            disabled={pending.addBlock || !currentSection}
                        >
                            Добавить блок
                        </button>
                        <button
                            type="button"
                            className={`${styles.actionButton} ${styles.dangerAction}`}
                            onClick={handleDeleteSection}
                            disabled={
                                pending.deleteSection || !currentSection
                            }
                        >
                            Удалить раздел
                        </button>
                    </div>
                )}
            </header>

            {/* Контент текущего раздела */}
            <main className={styles.main}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionBadge}>
                        {currentSection
                            ? currentSection.recommendationType
                            : "Раздел отсутствует"}
                    </div>
                    <h2 className={styles.sectionTitle}>
                        {currentSection
                            ? "Список рекомендаций"
                            : "Добавьте раздел, чтобы начать"}
                    </h2>
                </div>

                {renderBlocks()}
            </main>

            {/* Навигация по разделам */}
            <footer className={styles.pagination}>
                {sections.map((section, index) => (
                    <button
                        key={section.recommendationType}
                        type="button"
                        className={`${styles.pageButton} ${
                            index === currentPage ? styles.pageButtonActive : ""
                        }`}
                        onClick={() => handlePageChange(index)}
                    >
                        {getPageNumber(section.recommendationType) || index + 1}
                    </button>
                ))}
            </footer>
        </div>
    );
};

export default RecommendationsPage;
