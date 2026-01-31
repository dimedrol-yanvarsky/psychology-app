import { useCallback, useEffect, useMemo, useState } from "react";
import { useLockBodyScroll } from "../../../shared/lib/hooks/useLockBodyScroll";
import {
    addTest,
    changeTest,
    deleteTest,
    fetchQuestions,
    fetchTests as fetchTestsApi,
    submitAttempt,
} from "../../../entities/test";
import {
    createQuestion,
    normalizeAnswerOptions,
    normalizeAuthorsInput,
    normalizeQuestionsForSave,
    getQuestionNumber,
} from "../../../entities/test";

const getInitialAttemptModal = () => ({
    open: false,
    loading: false,
    error: "",
    test: null,
    questions: [],
    selected: {},
    submitting: false,
});

const getInitialEditModal = () => ({
    open: false,
    loading: false,
    saving: false,
    error: "",
    testId: "",
    testName: "",
    authorsInput: "",
    description: "",
    questions: [],
});

const getInitialAddModal = () => ({
    open: false,
    saving: false,
    error: "",
    form: {
        testName: "",
        authorsInput: "",
        description: "",
        questions: [createQuestion(1)],
    },
});

export const useTestsList = ({
    showAlert,
    profileData = {},
}) => {
    const userId = (profileData?.id || "").trim();

    const [tests, setTests] = useState([]);
    const [isListLoading, setIsListLoading] = useState(true);
    const [listError, setListError] = useState("");
    const [deletingTests, setDeletingTests] = useState({});

    const [attemptModal, setAttemptModal] = useState(getInitialAttemptModal);

    const [editModal, setEditModal] = useState(getInitialEditModal);

    const [addModal, setAddModal] = useState(getInitialAddModal);

    useLockBodyScroll(attemptModal.open || editModal.open || addModal.open);

    const fetchTests = useCallback(async () => {
        setIsListLoading(true);
        setListError("");

        try {
            const { data } = await fetchTestsApi({ userId });
            const loadedTests = Array.isArray(data?.tests) ? data.tests : [];
            setTests(loadedTests);
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось загрузить тестирования";
            setListError(message);
            setTests([]);
            showAlert?.("error", message);
        } finally {
            setIsListLoading(false);
        }
    }, [showAlert, userId]);

    useEffect(() => {
        fetchTests();
    }, [fetchTests]);

    const emptyState = useMemo(
        () => !isListLoading && !listError && tests.length === 0,
        [isListLoading, listError, tests]
    );

    const setDeletingFlag = (testId, value) => {
        setDeletingTests((prev) => {
            const updated = { ...prev };
            if (value) {
                updated[testId] = true;
            } else {
                delete updated[testId];
            }
            return updated;
        });
    };

    const closeAttemptModal = () => {
        setAttemptModal(getInitialAttemptModal());
    };

    const closeEditModal = () => {
        setEditModal(getInitialEditModal());
    };

    const closeAddModal = () => {
        setAddModal(getInitialAddModal());
    };

    const handleDeleteTest = async (test) => {
        const testId = test?.id;
        if (!testId) {
            return;
        }

        setDeletingFlag(testId, true);
        try {
            const { data } = await deleteTest({ testId });

            if (data?.status === "success") {
                setTests((prev) => prev.filter((item) => item.id !== testId));
                showAlert?.("success", data?.message || "Тестирование удалено");
            }
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось удалить тестирование";
            showAlert?.("error", message);
        } finally {
            setDeletingFlag(testId, false);
        }
    };

    const handleOpenAttempt = async (test) => {
        setAttemptModal({
            open: true,
            loading: true,
            error: "",
            test,
            questions: [],
            selected: {},
            submitting: false,
        });

        try {
            const { data } = await fetchQuestions({
                testId: test.id,
            });

            const questions = Array.isArray(data?.questions)
                ? data.questions
                : [];

            setAttemptModal((prev) => ({
                ...prev,
                loading: false,
                questions,
            }));
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось загрузить вопросы";
            setAttemptModal((prev) => ({
                ...prev,
                loading: false,
                error: message,
            }));
            showAlert?.("error", message);
        }
    };

    const toggleAnswer = (questionId, optionIndex, isSingle) => {
        setAttemptModal((prev) => {
            const selected = { ...prev.selected };
            const key = String(questionId);
            const current = new Set(selected[key] || []);

            if (isSingle) {
                current.clear();
                current.add(optionIndex);
            } else if (current.has(optionIndex)) {
                current.delete(optionIndex);
            } else {
                current.add(optionIndex);
            }

            selected[key] = Array.from(current).sort((a, b) => a - b);

            return { ...prev, selected };
        });
    };

    const allQuestionsAnswered =
        attemptModal.questions.length > 0 &&
        attemptModal.questions.every((question, index) => {
            const questionNumber = getQuestionNumber(question, index);
            return (attemptModal.selected[String(questionNumber)] || []).length;
        });

    const handleSubmitAttempt = async () => {
        if (!attemptModal.test || attemptModal.submitting) {
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

        setAttemptModal((prev) => ({ ...prev, submitting: true }));

        const answersPayload = attemptModal.questions.map((question, index) => {
            const questionNumber = getQuestionNumber(question, index);
            const selected =
                attemptModal.selected[String(questionNumber)] || [];
            return [questionNumber, ...selected];
        });

        try {
            const { data } = await submitAttempt({
                testId: attemptModal.test.id,
                userId,
                answers: answersPayload,
            });

            showAlert?.(
                "success",
                data?.message || "Тестирование завершено успешно"
            );

            setTests((prev) =>
                prev.map((item) =>
                    item.id === attemptModal.test.id
                        ? { ...item, isCompleted: true }
                        : item
                )
            );

            closeAttemptModal();
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось сохранить результаты";
            showAlert?.("error", message);
        } finally {
            setAttemptModal((prev) => ({ ...prev, submitting: false }));
        }
    };

    const handleOpenEdit = async (test) => {
        if (!test?.id) {
            return;
        }

        setEditModal({
            open: true,
            loading: true,
            saving: false,
            error: "",
            testId: test.id,
            testName: "",
            authorsInput: "",
            description: "",
            questions: [],
        });

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

            setEditModal((prev) => ({
                ...prev,
                loading: false,
                testId: data?.testId || test.id,
                testName: data?.test?.testName || test.testName || "",
                authorsInput: authors,
                description: data?.test?.description || test.description || "",
                questions: questions.length ? questions : [createQuestion(1)],
            }));
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                error?.message ||
                "Не удалось загрузить тестирование";
            setEditModal((prev) => ({
                ...prev,
                loading: false,
                error: message,
            }));
            showAlert?.("error", message);
        }
    };

    const updateEditQuestion = (index, updater) => {
        setEditModal((prev) => {
            const questions = [...prev.questions];
            const current = questions[index] || createQuestion(index + 1);
            questions[index] =
                typeof updater === "function"
                    ? updater(current)
                    : { ...current, ...updater };
            return { ...prev, questions };
        });
    };

    const addEditQuestion = () => {
        setEditModal((prev) => {
            const nextId =
                (prev.questions[prev.questions.length - 1]?.id || 0) + 1;
            return {
                ...prev,
                questions: [...prev.questions, createQuestion(nextId)],
            };
        });
    };

    const addEditOption = (qIndex) => {
        updateEditQuestion(qIndex, (question) => ({
            ...question,
            answerOptions: [...(question.answerOptions || []), ""],
        }));
    };

    const handleSaveEdit = async () => {
        if (editModal.saving) {
            return;
        }

        const testName = editModal.testName.trim();
        const description = editModal.description.trim();
        const authors = normalizeAuthorsInput(editModal.authorsInput || "");

        if (!testName || !description || authors.length === 0) {
            showAlert?.(
                "error",
                "Заполните название, описание и авторов теста"
            );
            return;
        }

        const { questions, error } = normalizeQuestionsForSave(
            editModal.questions
        );
        if (error) {
            showAlert?.("error", error);
            return;
        }

        setEditModal((prev) => ({ ...prev, saving: true }));

        try {
            const { data } = await changeTest({
                action: "update",
                testId: editModal.testId,
                testName,
                description,
                authorsName: authors,
                questions,
            });

            showAlert?.("success", data?.message || "Изменения сохранены");

            if (data?.test?.id) {
                setTests((prev) =>
                    prev.map((item) =>
                        item.id === data.test.id
                            ? {
                                  ...item,
                                  testName: data.test.testName,
                                  description: data.test.description,
                                  questionCount: data.test.questionCount,
                              }
                            : item
                    )
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
            setEditModal((prev) => ({ ...prev, saving: false }));
        }
    };

    const openAddModal = () => {
        setAddModal({
            open: true,
            saving: false,
            error: "",
            form: {
                testName: "",
                authorsInput: "",
                description: "",
                questions: [createQuestion(1)],
            },
        });
    };

    const updateAddQuestion = (index, updater) => {
        setAddModal((prev) => {
            const questions = [...prev.form.questions];
            const current = questions[index] || createQuestion(index + 1);
            questions[index] =
                typeof updater === "function"
                    ? updater(current)
                    : { ...current, ...updater };
            return { ...prev, form: { ...prev.form, questions } };
        });
    };

    const addQuestionToForm = () => {
        setAddModal((prev) => {
            const nextId =
                (prev.form.questions[prev.form.questions.length - 1]?.id || 0) +
                1;
            return {
                ...prev,
                form: {
                    ...prev.form,
                    questions: [...prev.form.questions, createQuestion(nextId)],
                },
            };
        });
    };

    const addOptionToForm = (qIndex) => {
        updateAddQuestion(qIndex, (question) => ({
            ...question,
            answerOptions: [...(question.answerOptions || []), ""],
        }));
    };

    const handleAddTest = async () => {
        if (addModal.saving) {
            return;
        }

        if (!userId) {
            showAlert?.("error", "Авторизуйтесь, чтобы добавить тестирование");
            return;
        }

        const testName = addModal.form.testName.trim();
        const description = addModal.form.description.trim();
        const authors = normalizeAuthorsInput(addModal.form.authorsInput || "");

        if (!testName || !description || authors.length === 0) {
            showAlert?.(
                "error",
                "Заполните название, описание и авторов теста"
            );
            return;
        }

        const { questions, error } = normalizeQuestionsForSave(
            addModal.form.questions
        );
        if (error) {
            showAlert?.("error", error);
            return;
        }

        setAddModal((prev) => ({ ...prev, saving: true }));

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
            setAddModal((prev) => ({ ...prev, saving: false }));
        }
    };

    return {
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
        fetchTests,
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
    };
};
