import { createSlice } from "@reduxjs/toolkit";
import { createQuestion } from "../../../entities/test";

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

const initialState = {
    tests: [],
    isListLoading: true,
    listError: "",
    deletingTests: {},
    attemptModal: getInitialAttemptModal(),
    editModal: getInitialEditModal(),
    addModal: getInitialAddModal(),
};

// Слайс тестирований: список, прохождение, редактирование, добавление.
const testsSlice = createSlice({
    name: "tests",
    initialState,
    reducers: {
        // Загрузка списка
        fetchTestsStart(state) {
            state.isListLoading = true;
            state.listError = "";
        },
        fetchTestsSuccess(state, action) {
            state.isListLoading = false;
            state.tests = action.payload;
        },
        fetchTestsError(state, action) {
            state.isListLoading = false;
            state.tests = [];
            state.listError = action.payload;
        },

        // Удаление
        setDeletingFlag(state, action) {
            state.deletingTests[action.payload] = true;
        },
        clearDeletingFlag(state, action) {
            delete state.deletingTests[action.payload];
        },
        removeTest(state, action) {
            state.tests = state.tests.filter(
                (test) => test.id !== action.payload
            );
        },

        // Обновление
        updateTest(state, action) {
            const { id, updates } = action.payload;
            const test = state.tests.find((t) => t.id === id);
            if (test) {
                Object.assign(test, updates);
            }
        },
        markTestCompleted(state, action) {
            const test = state.tests.find((t) => t.id === action.payload);
            if (test) {
                test.isCompleted = true;
            }
        },

        // Модалка прохождения
        openAttemptModal(state, action) {
            state.attemptModal = {
                open: true,
                loading: true,
                error: "",
                test: action.payload,
                questions: [],
                selected: {},
                submitting: false,
            };
        },
        closeAttemptModal(state) {
            state.attemptModal = getInitialAttemptModal();
        },
        setAttemptQuestions(state, action) {
            state.attemptModal.loading = false;
            state.attemptModal.questions = action.payload;
        },
        setAttemptError(state, action) {
            state.attemptModal.loading = false;
            state.attemptModal.error = action.payload;
        },
        toggleAnswer(state, action) {
            const { questionId, optionIndex, isSingle } = action.payload;
            const key = String(questionId);
            const current = new Set(state.attemptModal.selected[key] || []);

            if (isSingle) {
                current.clear();
                current.add(optionIndex);
            } else if (current.has(optionIndex)) {
                current.delete(optionIndex);
            } else {
                current.add(optionIndex);
            }

            state.attemptModal.selected[key] = Array.from(current).sort(
                (a, b) => a - b
            );
        },
        submitAttemptStart(state) {
            state.attemptModal.submitting = true;
        },
        submitAttemptEnd(state) {
            state.attemptModal.submitting = false;
        },

        // Модалка редактирования
        openEditModal(state, action) {
            state.editModal = {
                open: true,
                loading: true,
                saving: false,
                error: "",
                testId: action.payload.testId,
                testName: "",
                authorsInput: "",
                description: "",
                questions: [],
            };
        },
        closeEditModal(state) {
            state.editModal = getInitialEditModal();
        },
        setEditData(state, action) {
            state.editModal.loading = false;
            Object.assign(state.editModal, action.payload);
        },
        setEditError(state, action) {
            state.editModal.loading = false;
            state.editModal.error = action.payload;
        },
        setEditSaving(state, action) {
            state.editModal.saving = action.payload;
        },
        updateEditQuestions(state, action) {
            const { index, updater } = action.payload;
            const current =
                state.editModal.questions[index] || createQuestion(index + 1);
            state.editModal.questions[index] =
                typeof updater === "function"
                    ? updater(current)
                    : { ...current, ...updater };
        },
        addEditQuestion(state) {
            const lastId =
                state.editModal.questions[
                    state.editModal.questions.length - 1
                ]?.id || 0;
            state.editModal.questions.push(createQuestion(lastId + 1));
        },
        addEditOption(state, action) {
            const qIndex = action.payload;
            const question =
                state.editModal.questions[qIndex] ||
                createQuestion(qIndex + 1);
            if (!question.answerOptions) {
                question.answerOptions = [];
            }
            question.answerOptions.push("");
            state.editModal.questions[qIndex] = question;
        },
        setEditModal(state, action) {
            Object.assign(state.editModal, action.payload);
        },

        // Модалка добавления
        openAddModal(state) {
            state.addModal = {
                open: true,
                saving: false,
                error: "",
                form: {
                    testName: "",
                    authorsInput: "",
                    description: "",
                    questions: [createQuestion(1)],
                },
            };
        },
        closeAddModal(state) {
            state.addModal = getInitialAddModal();
        },
        setAddSaving(state, action) {
            state.addModal.saving = action.payload;
        },
        updateAddQuestions(state, action) {
            const { index, updater } = action.payload;
            const current =
                state.addModal.form.questions[index] ||
                createQuestion(index + 1);
            state.addModal.form.questions[index] =
                typeof updater === "function"
                    ? updater(current)
                    : { ...current, ...updater };
        },
        addQuestionToForm(state) {
            const lastId =
                state.addModal.form.questions[
                    state.addModal.form.questions.length - 1
                ]?.id || 0;
            state.addModal.form.questions.push(createQuestion(lastId + 1));
        },
        addOptionToForm(state, action) {
            const qIndex = action.payload;
            const question =
                state.addModal.form.questions[qIndex] ||
                createQuestion(qIndex + 1);
            if (!question.answerOptions) {
                question.answerOptions = [];
            }
            question.answerOptions.push("");
            state.addModal.form.questions[qIndex] = question;
        },
        setAddModal(state, action) {
            Object.assign(state.addModal, action.payload);
        },
    },
});

export const {
    fetchTestsStart,
    fetchTestsSuccess,
    fetchTestsError,
    setDeletingFlag,
    clearDeletingFlag,
    removeTest,
    updateTest,
    markTestCompleted,
    openAttemptModal,
    closeAttemptModal,
    setAttemptQuestions,
    setAttemptError,
    toggleAnswer,
    submitAttemptStart,
    submitAttemptEnd,
    openEditModal,
    closeEditModal,
    setEditData,
    setEditError,
    setEditSaving,
    updateEditQuestions,
    addEditQuestion,
    addEditOption,
    setEditModal,
    openAddModal,
    closeAddModal,
    setAddSaving,
    updateAddQuestions,
    addQuestionToForm,
    addOptionToForm,
    setAddModal,
} = testsSlice.actions;

export default testsSlice.reducer;
