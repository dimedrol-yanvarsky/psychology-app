import { useCallback, useEffect, useMemo } from "react";
import { useSelector, useDispatch } from "react-redux";
import { useLockBodyScroll } from "../../../shared/lib/hooks/useLockBodyScroll";
import {
    addTest,
    changeTest,
    deleteTest,
    fetchQuestions,
    fetchTests as fetchTestsApi,
    submitAttempt,
} from "../../../entities/test";
import { useAuthContext } from "../../../shared/context/AuthContext";
import { useAlertContext } from "../../../shared/context/AlertContext";
import {
    createQuestion,
    normalizeAnswerOptions,
    normalizeAuthorsInput,
    normalizeQuestionsForSave,
    getQuestionNumber,
} from "../../../entities/test";
import {
    fetchTestsStart,
    fetchTestsSuccess,
    fetchTestsError,
    setDeletingFlag,
    clearDeletingFlag,
    removeTest,
    updateTest,
    markTestCompleted,
    openAttemptModal,
    closeAttemptModal as closeAttemptModalAction,
    setAttemptQuestions,
    setAttemptError,
    toggleAnswer as toggleAnswerAction,
    submitAttemptStart,
    submitAttemptEnd,
    openEditModal,
    closeEditModal as closeEditModalAction,
    setEditData,
    setEditError,
    setEditSaving,
    updateEditQuestions,
    addEditQuestion as addEditQuestionAction,
    addEditOption as addEditOptionAction,
    setEditModal,
    openAddModal as openAddModalAction,
    closeAddModal as closeAddModalAction,
    setAddSaving,
    updateAddQuestions,
    addQuestionToForm as addQuestionToFormAction,
    addOptionToForm as addOptionToFormAction,
    setAddModal,
} from "./testsSlice";

export const useTestsList = () => {
    const { profileData } = useAuthContext();
    const { showAlert } = useAlertContext();
    const reduxDispatch = useDispatch();

    const userId = (profileData?.id || "").trim();

    const state = useSelector((s) => s.tests);

    useLockBodyScroll(
        state.attemptModal.open || state.editModal.open || state.addModal.open
    );

    const fetchTests = useCallback(async () => {
        reduxDispatch(fetchTestsStart());

        try {
            const { data } = await fetchTestsApi({ userId });
            const loadedTests = Array.isArray(data?.tests) ? data.tests : [];
            reduxDispatch(fetchTestsSuccess(loadedTests));
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось загрузить тестирования";
            reduxDispatch(fetchTestsError(message));
            showAlert?.("error", message);
        }
    }, [showAlert, userId, reduxDispatch]);

    useEffect(() => {
        fetchTests();
    }, [fetchTests]);

    const emptyState = useMemo(
        () =>
            !state.isListLoading &&
            !state.listError &&
            state.tests.length === 0,
        [state.isListLoading, state.listError, state.tests]
    );

    const closeAttemptModal = () => {
        reduxDispatch(closeAttemptModalAction());
    };

    const closeEditModal = () => {
        reduxDispatch(closeEditModalAction());
    };

    const closeAddModal = () => {
        reduxDispatch(closeAddModalAction());
    };

    const handleDeleteTest = async (test) => {
        const testId = test?.id;
        if (!testId) {
            return;
        }

        reduxDispatch(setDeletingFlag(testId));
        try {
            const { data } = await deleteTest({ testId });

            if (data?.status === "success") {
                reduxDispatch(removeTest(testId));
                showAlert?.("success", data?.message || "Тестирование удалено");
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось удалить тестирование";
            showAlert?.("error", message);
        } finally {
            reduxDispatch(clearDeletingFlag(testId));
        }
    };

    const handleOpenAttempt = async (test) => {
        reduxDispatch(openAttemptModal(test));

        try {
            const { data } = await fetchQuestions({
                testId: test.id,
            });

            const questions = Array.isArray(data?.questions)
                ? data.questions
                : [];

            reduxDispatch(setAttemptQuestions(questions));
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось загрузить вопросы";
            reduxDispatch(setAttemptError(message));
            showAlert?.("error", message);
        }
    };

    const toggleAnswer = (questionId, optionIndex, isSingle) => {
        reduxDispatch(
            toggleAnswerAction({ questionId, optionIndex, isSingle })
        );
    };

    const allQuestionsAnswered =
        state.attemptModal.questions.length > 0 &&
        state.attemptModal.questions.every((question, index) => {
            const questionNumber = getQuestionNumber(question, index);
            return (state.attemptModal.selected[String(questionNumber)] || [])
                .length;
        });

    const handleSubmitAttempt = async () => {
        if (!state.attemptModal.test || state.attemptModal.submitting) {
            return;
        }

        if (!userId) {
            showAlert?.("error", "Авторизуйтесь, чтобы пройти тест");
            return;
        }

        if (!allQuestionsAnswered) {
            showAlert?.("error", "Выберите ответы на все вопросы");
            return;
        }

        reduxDispatch(submitAttemptStart());

        const answersPayload = state.attemptModal.questions.map(
            (question, index) => {
                const questionNumber = getQuestionNumber(question, index);
                const selected =
                    state.attemptModal.selected[String(questionNumber)] || [];
                return [questionNumber, ...selected];
            }
        );

        try {
            const { data } = await submitAttempt({
                testId: state.attemptModal.test.id,
                userId,
                answers: answersPayload,
            });

            showAlert?.(
                "success",
                data?.message || "Тестирование завершено успешно"
            );

            reduxDispatch(markTestCompleted(state.attemptModal.test.id));

            closeAttemptModal();
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось сохранить результаты";
            showAlert?.("error", message);
        } finally {
            reduxDispatch(submitAttemptEnd());
        }
    };

    const handleOpenEdit = async (test) => {
        if (!test?.id) {
            return;
        }

        reduxDispatch(openEditModal({ testId: test.id }));

        try {
            const { data } = await changeTest({
                testId: test.id,
                action: "load",
            });

            const questions = Array.isArray(data?.questions)
                ? data.questions.map((question) => ({
                      id: question.id,
                      questionBody:
                          question.questionBody ||
                          question.question ||
                          question.body ||
                          "",
                      answerOptions: normalizeAnswerOptions(
                          question.answerOptions
                      ),
                      selectType:
                          question.selectType || question.selectype || "one",
                  }))
                : [];

            const authors =
                data?.test?.authorsName && Array.isArray(data.test.authorsName)
                    ? data.test.authorsName.join(", ")
                    : "";

            reduxDispatch(
                setEditData({
                    testId: data?.testId || test.id,
                    testName: data?.test?.testName || test.testName || "",
                    authorsInput: authors,
                    description:
                        data?.test?.description || test.description || "",
                    questions: questions.length
                        ? questions
                        : [createQuestion(1)],
                })
            );
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось загрузить тестирование";
            reduxDispatch(setEditError(message));
            showAlert?.("error", message);
        }
    };

    const handleUpdateEditQuestion = (index, updater) => {
        reduxDispatch(updateEditQuestions({ index, updater }));
    };

    const addEditQuestion = () => {
        reduxDispatch(addEditQuestionAction());
    };

    const addEditOption = (qIndex) => {
        reduxDispatch(addEditOptionAction(qIndex));
    };

    const handleSaveEdit = async () => {
        if (state.editModal.saving) {
            return;
        }

        const testName = state.editModal.testName.trim();
        const description = state.editModal.description.trim();
        const authors = normalizeAuthorsInput(
            state.editModal.authorsInput || ""
        );

        if (!testName || !description || authors.length === 0) {
            showAlert?.(
                "error",
                "Заполните название, описание и авторов теста"
            );
            return;
        }

        const { questions, error } = normalizeQuestionsForSave(
            state.editModal.questions
        );
        if (error) {
            showAlert?.("error", error);
            return;
        }

        reduxDispatch(setEditSaving(true));

        try {
            const { data } = await changeTest({
                action: "update",
                testId: state.editModal.testId,
                testName,
                description,
                authorsName: authors,
                questions,
            });

            showAlert?.("success", data?.message || "Изменения сохранены");

            if (data?.test?.id) {
                reduxDispatch(
                    updateTest({
                        id: data.test.id,
                        updates: {
                            testName: data.test.testName,
                            description: data.test.description,
                            questionCount: data.test.questionCount,
                        },
                    })
                );
            }

            closeEditModal();
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось сохранить изменения";
            showAlert?.("error", message);
        } finally {
            reduxDispatch(setEditSaving(false));
        }
    };

    const openAddModal = () => {
        reduxDispatch(openAddModalAction());
    };

    const updateAddQuestion = (index, updater) => {
        reduxDispatch(updateAddQuestions({ index, updater }));
    };

    const addQuestionToForm = () => {
        reduxDispatch(addQuestionToFormAction());
    };

    const addOptionToForm = (qIndex) => {
        reduxDispatch(addOptionToFormAction(qIndex));
    };

    const handleAddTest = async () => {
        if (state.addModal.saving) {
            return;
        }

        if (!userId) {
            showAlert?.("error", "Авторизуйтесь, чтобы добавить тестирование");
            return;
        }

        const testName = state.addModal.form.testName.trim();
        const description = state.addModal.form.description.trim();
        const authors = normalizeAuthorsInput(
            state.addModal.form.authorsInput || ""
        );

        if (!testName || !description || authors.length === 0) {
            showAlert?.(
                "error",
                "Заполните название, описание и авторов теста"
            );
            return;
        }

        const { questions, error } = normalizeQuestionsForSave(
            state.addModal.form.questions
        );
        if (error) {
            showAlert?.("error", error);
            return;
        }

        reduxDispatch(setAddSaving(true));

        try {
            const { data } = await addTest({
                userId,
                testName,
                description,
                authorsName: authors,
                questions,
            });

            showAlert?.("success", data?.message || "Тестирование добавлено");

            fetchTests();
            closeAddModal();
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось добавить тестирование";
            showAlert?.("error", message);
        } finally {
            reduxDispatch(setAddSaving(false));
        }
    };

    return {
        addModal: state.addModal,
        addEditOption,
        addEditQuestion,
        addOptionToForm,
        addQuestionToForm,
        allQuestionsAnswered,
        attemptModal: state.attemptModal,
        closeAddModal,
        closeAttemptModal,
        closeEditModal,
        deletingTests: state.deletingTests,
        editModal: state.editModal,
        emptyState,
        fetchTests,
        handleAddTest,
        handleDeleteTest,
        handleOpenAttempt,
        handleOpenEdit,
        handleSaveEdit,
        handleSubmitAttempt,
        isListLoading: state.isListLoading,
        listError: state.listError,
        openAddModal,
        setAddModal: (payload) => reduxDispatch(setAddModal(payload)),
        setEditModal: (payload) => reduxDispatch(setEditModal(payload)),
        tests: state.tests,
        toggleAnswer,
        updateAddQuestion,
        updateEditQuestion: handleUpdateEditQuestion,
    };
};
