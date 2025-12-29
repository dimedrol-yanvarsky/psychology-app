import React from "react";
import styles from "./TestsPage.module.css";
import { useTestsPage } from "../model/useTestsPage";
import { getOptionValue, getQuestionNumber } from "../lib/testUtils";

const TestsPage = ({
    showAlert,
    isAdmin = false,
    isAuth,
    profileData = {},
}) => {
    const {
        addModal,
        addEditOption,
        addEditQuestion,
        addOptionToForm,
        addQuestionToForm,
        allQuestionsAnswered,
        attemptModal,
        closeAddModal,
        closeAttemptModal,
        closeEditModal,
        deletingTests,
        editModal,
        emptyState,
        handleAddTest,
        handleDeleteTest,
        handleOpenAttempt,
        handleOpenEdit,
        handleSaveEdit,
        handleSubmitAttempt,
        isListLoading,
        listError,
        openAddModal,
        setAddModal,
        setEditModal,
        tests,
        toggleAnswer,
        updateAddQuestion,
        updateEditQuestion,
    } = useTestsPage({ showAlert, profileData });

    // Рендер списка тестов и состояния загрузки.
    const renderTestsContent = () => {
        if (isListLoading) {
            return (
                <div className={styles.testsList}>
                    {[1, 2, 3].map((item) => (
                        <div
                            key={item}
                            className={styles.skeletonCard}
                            aria-hidden="true"
                        >
                            <div className={styles.skeletonHeader}>
                                <div className={styles.skeletonTitle} />
                                <div className={styles.skeletonBadge} />
                            </div>
                            <div className={styles.skeletonLine} />
                            <div className={styles.skeletonLine} />
                            <div className={styles.skeletonActions} />
                        </div>
                    ))}
                </div>
            );
        }

        if (listError) {
            return <div className={styles.emptyStateBanner}>{listError}</div>;
        }

        if (emptyState) {
            return (
                <div className={styles.emptyStateBanner}>
                    Пока нет доступных тестов. Добавьте первый тест.
                </div>
            );
        }

        return (
            <div className={styles.testsList}>
                {tests.map((test) => {
                    const questionCount =
                        test.questionCount || test.duration || 0;
                    const authors = Array.isArray(test.authorsName)
                        ? test.authorsName.join(", ")
                        : test.author || "";
                    const isCompleted = Boolean(test.isCompleted);
                    return (
                        <article key={test.id} className={styles.testCard}>
                            <header className={styles.cardHeader}>
                                <div className={styles.cardTitleBlock}>
                                    <h2 className={styles.testTitle}>
                                        {test.testName || test.title}
                                    </h2>
                                    {authors ? (
                                        <span className={styles.testAuthor}>
                                            Автор: {authors}
                                        </span>
                                    ) : null}
                                </div>
                                <span className={styles.durationBadge}>
                                    Количество вопросов: {questionCount}
                                </span>
                            </header>

                            <p className={styles.testSummary}>
                                {test.description || test.summary}
                            </p>

                            {isCompleted ? (
                                <div className={styles.completedBadge}>
                                    Тестирование уже пройдено
                                </div>
                            ) : (
                                <>
                                    {isAuth || isAdmin ? (
                                        <div className={styles.cardActions}>
                                            {isAuth ? (
                                                <button
                                                    type="button"
                                                    className={`${styles.cardButton} ${styles.primaryButton}`}
                                                    onClick={() =>
                                                        handleOpenAttempt(test)
                                                    }
                                                >
                                                    Пройти тест
                                                </button>
                                            ) : null}
                                            {isAdmin ? (
                                                <div
                                                    className={
                                                        styles.adminActions
                                                    }
                                                >
                                                    <button
                                                        type="button"
                                                        className={`${styles.cardButton} ${styles.editButton}`}
                                                        onClick={() =>
                                                            handleOpenEdit(test)
                                                        }
                                                    >
                                                        Изменить тестирование
                                                    </button>
                                                    <button
                                                        type="button"
                                                        className={`${styles.cardButton} ${styles.deleteButton}`}
                                                        onClick={() =>
                                                            handleDeleteTest(
                                                                test
                                                            )
                                                        }
                                                        disabled={Boolean(
                                                            deletingTests[
                                                                test.id
                                                            ]
                                                        )}
                                                    >
                                                        {deletingTests[
                                                            test.id
                                                        ]
                                                            ? "Удаляем..."
                                                            : "Удалить тестирование"}
                                                    </button>
                                                </div>
                                            ) : null}
                                        </div>
                                    ) : null}
                                </>
                            )}
                        </article>
                    );
                })}
            </div>
        );
    };

    // Рендер карточки вопроса в модальном окне прохождения.
    const renderAttemptQuestion = (question, questionIndex) => {
        const questionNumber = getQuestionNumber(question, questionIndex);
        const selectType =
            (question.selectType || question.selectype || "").toLowerCase() ||
            "one";
        const isSingle = selectType === "one";
        const options = Array.isArray(question.answerOptions)
            ? question.answerOptions
            : [];
        const selected =
            attemptModal.selected[String(questionNumber)] || [];

        return (
            <div key={questionNumber} className={styles.questionCard}>
                <div className={styles.questionTitle}>
                    {questionNumber}.{" "}
                    {question.questionBody ||
                        question.question ||
                        question.body ||
                        "Вопрос"}
                </div>
                <div className={styles.optionsList}>
                    {options.map((option, optionIndex) => {
                        const optionLabel =
                            typeof option === "string"
                                ? option
                                : option?.body ||
                                  option?.text ||
                                  option?.title ||
                                  `Вариант ${optionIndex + 1}`;
                        const optionNumber = optionIndex + 1;
                        const isChecked = selected.includes(optionNumber);

                        return (
                            <label
                                key={`${questionNumber}-${optionNumber}`}
                                className={`${styles.optionItem} ${
                                    isChecked ? styles.optionSelected : ""
                                }`}
                            >
                                <input
                                    type={isSingle ? "radio" : "checkbox"}
                                    name={`question-${questionNumber}`}
                                    checked={isChecked}
                                    onChange={() =>
                                        toggleAnswer(
                                            questionNumber,
                                            optionNumber,
                                            isSingle
                                        )
                                    }
                                />
                                <span>{optionLabel}</span>
                            </label>
                        );
                    })}
                </div>
            </div>
        );
    };

    // Контент модального окна прохождения теста.
    const renderAttemptModalContent = () => {
        if (attemptModal.loading) {
            return (
                <div className={styles.testsLoader}>
                    <div className={styles.loaderSpinner} aria-hidden="true" />
                    <span className={styles.loaderText}>
                        Загружаем вопросы...
                    </span>
                </div>
            );
        }

        if (attemptModal.error) {
            return <div className={styles.testsError}>{attemptModal.error}</div>;
        }

        return (
            <>
                <div className={styles.questionsList}>
                    {attemptModal.questions.map(renderAttemptQuestion)}
                </div>
                <div className={styles.modalActions}>
                    <button
                        type="button"
                        className={`${styles.cardButton} ${styles.secondaryButton}`}
                        onClick={closeAttemptModal}
                    >
                        Отменить
                    </button>
                    <button
                        type="button"
                        className={`${styles.cardButton} ${styles.primaryButton}`}
                        onClick={handleSubmitAttempt}
                        disabled={
                            attemptModal.submitting || !allQuestionsAnswered
                        }
                    >
                        {attemptModal.submitting
                            ? "Сохраняем..."
                            : "Завершить тест"}
                    </button>
                </div>
            </>
        );
    };

    return (
        <div className={styles.page}>
            {/* Верхняя панель страницы */}
            <header className={styles.topBar}>
                <h1 className={styles.title}>Психологическое тестирование</h1>
                {isAdmin ? (
                    <button
                        type="button"
                        className={`${styles.actionButton} ${styles.primaryAction}`}
                        onClick={openAddModal}
                    >
                        Добавить тестирование
                    </button>
                ) : null}
            </header>

            {/* Список тестов и состояния */}
            <main className={styles.main}>{renderTestsContent()}</main>

            {/* Модальное окно прохождения теста */}
            {attemptModal.open && (
                <div
                    className={styles.modalOverlay}
                    onClick={closeAttemptModal}
                >
                    <div
                        className={styles.answersModal}
                        onClick={(event) => event.stopPropagation()}
                        role="dialog"
                        aria-modal="true"
                    >
                        <button
                            type="button"
                            className={styles.modalClose}
                            onClick={closeAttemptModal}
                            aria-label="Закрыть модальное окно"
                        >
                            ×
                        </button>
                        <p className={styles.modalOverline}>Тестирование</p>
                        <h4 className={styles.modalTitle}>
                            {attemptModal.test?.testName ||
                                attemptModal.test?.title ||
                                "Прохождение теста"}
                        </h4>
                        {renderAttemptModalContent()}
                    </div>
                </div>
            )}

            {/* Модальное окно редактирования теста */}
            {editModal.open && (
                <div className={styles.modalOverlay} onClick={closeEditModal}>
                    <div
                        className={styles.formModal}
                        onClick={(event) => event.stopPropagation()}
                        role="dialog"
                        aria-modal="true"
                    >
                        <button
                            type="button"
                            className={styles.modalClose}
                            onClick={closeEditModal}
                            aria-label="Закрыть модальное окно"
                        >
                            ×
                        </button>
                        <p className={styles.modalOverline}>Редактирование</p>
                        <h4 className={styles.modalTitle}>
                            Изменить тестирование
                        </h4>

                        {editModal.loading ? (
                            <div className={styles.testsLoader}>
                                <div
                                    className={styles.loaderSpinner}
                                    aria-hidden="true"
                                />
                                <span className={styles.loaderText}>
                                    Загружаем данные...
                                </span>
                            </div>
                        ) : editModal.error ? (
                            <div className={styles.testsError}>
                                {editModal.error}
                            </div>
                        ) : (
                            <div className={styles.formContent}>
                                <label className={styles.fieldLabel}>
                                    Название тестирования
                                    <input
                                        type="text"
                                        className={styles.input}
                                        value={editModal.testName}
                                        onChange={(event) =>
                                            setEditModal((prev) => ({
                                                ...prev,
                                                testName: event.target.value,
                                            }))
                                        }
                                        placeholder="Введите название"
                                    />
                                </label>

                                <label className={styles.fieldLabel}>
                                    Авторы (через запятую)
                                    <input
                                        type="text"
                                        className={styles.input}
                                        value={editModal.authorsInput}
                                        onChange={(event) =>
                                            setEditModal((prev) => ({
                                                ...prev,
                                                authorsInput:
                                                    event.target.value,
                                            }))
                                        }
                                        placeholder="Например: Ч. Д. Спилбергер, Ю. Л. Ханин"
                                    />
                                </label>

                                <label className={styles.fieldLabel}>
                                    Описание тестирования
                                    <textarea
                                        className={styles.textarea}
                                        value={editModal.description}
                                        onChange={(event) =>
                                            setEditModal((prev) => ({
                                                ...prev,
                                                description: event.target.value,
                                            }))
                                        }
                                        placeholder="Расскажите, что оценивает тест"
                                    />
                                </label>

                                <div className={styles.questionsEditor}>
                                    <div className={styles.editorHeader}>
                                        <span>Вопросы</span>
                                        <button
                                            type="button"
                                            className={`${styles.cardButton} ${styles.secondaryButton}`}
                                            onClick={addEditQuestion}
                                        >
                                            Добавить вопрос
                                        </button>
                                    </div>

                                    {editModal.questions.map(
                                        (question, qIndex) => (
                                            <div
                                                key={question.id || qIndex}
                                                className={
                                                    styles.questionEditor
                                                }
                                            >
                                                <div
                                                    className={
                                                        styles.questionRow
                                                    }
                                                >
                                                    <span
                                                        className={
                                                            styles.questionBadge
                                                        }
                                                    >
                                                        Вопрос {qIndex + 1}
                                                    </span>
                                                    <select
                                                        className={
                                                            styles.select
                                                        }
                                                        value={
                                                            question.selectType ||
                                                            "one"
                                                        }
                                                        onChange={(event) =>
                                                            updateEditQuestion(
                                                                qIndex,
                                                                {
                                                                    selectType:
                                                                        event
                                                                            .target
                                                                            .value,
                                                                }
                                                            )
                                                        }
                                                    >
                                                        <option value="one">
                                                            Один вариант
                                                        </option>
                                                        <option value="couple">
                                                            Несколько вариантов
                                                        </option>
                                                    </select>
                                                </div>
                                                <input
                                                    type="text"
                                                    className={styles.input}
                                                    value={
                                                        question.questionBody ||
                                                        ""
                                                    }
                                                    onChange={(event) =>
                                                        updateEditQuestion(
                                                            qIndex,
                                                            {
                                                                questionBody:
                                                                    event.target
                                                                        .value,
                                                            }
                                                        )
                                                    }
                                                    placeholder="Формулировка вопроса"
                                                />

                                                <div
                                                    className={
                                                        styles.optionsEditor
                                                    }
                                                >
                                                    {Array.isArray(
                                                        question.answerOptions
                                                    ) &&
                                                        question.answerOptions.map(
                                                            (
                                                                option,
                                                                optionIndex
                                                            ) => (
                                                                <div
                                                                    key={`${qIndex}-${optionIndex}`}
                                                                    className={
                                                                        styles.optionRow
                                                                    }
                                                                >
                                                                    <input
                                                                        type="text"
                                                                        className={
                                                                            styles.optionInput
                                                                        }
                                                                        value={
                                                                            getOptionValue(
                                                                                option
                                                                            )
                                                                        }
                                                                        onChange={(
                                                                            event
                                                                        ) =>
                                                                            updateEditQuestion(
                                                                                qIndex,
                                                                                (
                                                                                    current
                                                                                ) => {
                                                                                    const nextOptions = [
                                                                                        ...(current.answerOptions ||
                                                                                            []),
                                                                                    ];
                                                                                    nextOptions[optionIndex] =
                                                                                        event
                                                                                            .target
                                                                                            .value;
                                                                                    return {
                                                                                        ...current,
                                                                                        answerOptions:
                                                                                            nextOptions,
                                                                                    };
                                                                                }
                                                                            )
                                                                        }
                                                                        placeholder={`Вариант ответа ${
                                                                            optionIndex +
                                                                            1
                                                                        }`}
                                                                    />
                                                                </div>
                                                            )
                                                        )}

                                                    <button
                                                        type="button"
                                                        className={`${styles.cardButton} ${styles.ghostButton}`}
                                                        onClick={() =>
                                                            addEditOption(
                                                                qIndex
                                                            )
                                                        }
                                                    >
                                                        Добавить вариант
                                                    </button>
                                                </div>
                                            </div>
                                        )
                                    )}
                                </div>

                                <div className={styles.modalActions}>
                                    <button
                                        type="button"
                                        className={`${styles.cardButton} ${styles.secondaryButton}`}
                                        onClick={closeEditModal}
                                    >
                                        Отменить
                                    </button>
                                    <button
                                        type="button"
                                        className={`${styles.cardButton} ${styles.primaryButton}`}
                                        onClick={handleSaveEdit}
                                        disabled={editModal.saving}
                                    >
                                        {editModal.saving
                                            ? "Сохраняем..."
                                            : "Сохранить изменения"}
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            )}

            {/* Модальное окно добавления теста */}
            {addModal.open && (
                <div className={styles.modalOverlay} onClick={closeAddModal}>
                    <div
                        className={styles.formModal}
                        onClick={(event) => event.stopPropagation()}
                        role="dialog"
                        aria-modal="true"
                    >
                        <button
                            type="button"
                            className={styles.modalClose}
                            onClick={closeAddModal}
                            aria-label="Закрыть модальное окно"
                        >
                            ×
                        </button>
                        <p className={styles.modalOverline}>
                            Новое тестирование
                        </p>
                        <h4 className={styles.modalTitle}>
                            Добавить тестирование
                        </h4>

                        <div className={styles.formContent}>
                            <label className={styles.fieldLabel}>
                                Название тестирования
                                <input
                                    type="text"
                                    className={styles.input}
                                    value={addModal.form.testName}
                                    onChange={(event) =>
                                        setAddModal((prev) => ({
                                            ...prev,
                                            form: {
                                                ...prev.form,
                                                testName: event.target.value,
                                            },
                                        }))
                                    }
                                    placeholder="Введите название"
                                />
                            </label>

                            <label className={styles.fieldLabel}>
                                Авторы (через запятую)
                                <input
                                    type="text"
                                    className={styles.input}
                                    value={addModal.form.authorsInput}
                                    onChange={(event) =>
                                        setAddModal((prev) => ({
                                            ...prev,
                                            form: {
                                                ...prev.form,
                                                authorsInput:
                                                    event.target.value,
                                            },
                                        }))
                                    }
                                    placeholder="Например: Ч. Д. Спилбергер"
                                />
                            </label>

                            <label className={styles.fieldLabel}>
                                Описание тестирования
                                <textarea
                                    className={styles.textarea}
                                    value={addModal.form.description}
                                    onChange={(event) =>
                                        setAddModal((prev) => ({
                                            ...prev,
                                            form: {
                                                ...prev.form,
                                                description: event.target.value,
                                            },
                                        }))
                                    }
                                    placeholder="Расскажите, что оценивает тест"
                                />
                            </label>

                            <div className={styles.questionsEditor}>
                                <div className={styles.editorHeader}>
                                    <span>Вопросы</span>
                                    <button
                                        type="button"
                                        className={`${styles.cardButton} ${styles.secondaryButton}`}
                                        onClick={addQuestionToForm}
                                    >
                                        Добавить вопрос
                                    </button>
                                </div>

                                {addModal.form.questions.map(
                                    (question, qIndex) => (
                                        <div
                                            key={question.id || qIndex}
                                            className={styles.questionEditor}
                                        >
                                            <div className={styles.questionRow}>
                                                <span
                                                    className={
                                                        styles.questionBadge
                                                    }
                                                >
                                                    Вопрос {qIndex + 1}
                                                </span>
                                                <select
                                                    className={styles.select}
                                                    value={
                                                        question.selectType ||
                                                        "one"
                                                    }
                                                    onChange={(event) =>
                                                        updateAddQuestion(
                                                            qIndex,
                                                            {
                                                                selectType:
                                                                    event.target
                                                                        .value,
                                                            }
                                                        )
                                                    }
                                                >
                                                    <option value="one">
                                                        Один вариант
                                                    </option>
                                                    <option value="couple">
                                                        Несколько вариантов
                                                    </option>
                                                </select>
                                            </div>
                                            <input
                                                type="text"
                                                className={styles.input}
                                                value={
                                                    question.questionBody || ""
                                                }
                                                onChange={(event) =>
                                                    updateAddQuestion(qIndex, {
                                                        questionBody:
                                                            event.target.value,
                                                    })
                                                }
                                                placeholder="Формулировка вопроса"
                                            />

                                            <div
                                                className={styles.optionsEditor}
                                            >
                                                {Array.isArray(
                                                    question.answerOptions
                                                ) &&
                                                    question.answerOptions.map(
                                                        (
                                                            option,
                                                            optionIndex
                                                        ) => (
                                                            <div
                                                                key={`${qIndex}-${optionIndex}`}
                                                                className={
                                                                    styles.optionRow
                                                                }
                                                            >
                                                                <input
                                                                    type="text"
                                                                    className={
                                                                        styles.optionInput
                                                                    }
                                                                    value={
                                                                        getOptionValue(
                                                                            option
                                                                        )
                                                                    }
                                                                    onChange={(
                                                                        event
                                                                    ) =>
                                                                        updateAddQuestion(
                                                                            qIndex,
                                                                            (
                                                                                current
                                                                            ) => {
                                                                                const nextOptions = [
                                                                                    ...(current.answerOptions ||
                                                                                        []),
                                                                                ];
                                                                                nextOptions[
                                                                                    optionIndex
                                                                                ] =
                                                                                    event.target.value;
                                                                                return {
                                                                                    ...current,
                                                                                    answerOptions:
                                                                                        nextOptions,
                                                                                };
                                                                            }
                                                                        )
                                                                    }
                                                                    placeholder={`Вариант ответа ${
                                                                        optionIndex +
                                                                        1
                                                                    }`}
                                                                />
                                                            </div>
                                                        )
                                                    )}
                                                <button
                                                    type="button"
                                                    className={`${styles.cardButton} ${styles.ghostButton}`}
                                                    onClick={() =>
                                                        addOptionToForm(qIndex)
                                                    }
                                                >
                                                    Добавить вариант
                                                </button>
                                            </div>
                                        </div>
                                    )
                                )}
                            </div>

                            <div className={styles.modalActions}>
                                <button
                                    type="button"
                                    className={`${styles.cardButton} ${styles.secondaryButton}`}
                                    onClick={closeAddModal}
                                >
                                    Отменить
                                </button>
                                <button
                                    type="button"
                                    className={`${styles.cardButton} ${styles.primaryButton}`}
                                    onClick={handleAddTest}
                                    disabled={addModal.saving}
                                >
                                    {addModal.saving
                                        ? "Сохраняем..."
                                        : "Сохранить тестирование"}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default TestsPage;
