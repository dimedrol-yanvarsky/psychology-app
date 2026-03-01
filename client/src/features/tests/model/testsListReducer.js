import { createQuestion } from "../../../entities/test";

// Начальное состояние модалки прохождения теста.
const getInitialAttemptModal = () => ({
    open: false,
    loading: false,
    error: "",
    test: null,
    questions: [],
    selected: {},
    submitting: false,
});

// Начальное состояние модалки редактирования теста.
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

// Начальное состояние модалки добавления теста.
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

// Начальное состояние списка тестов, объединяющее все состояния и модалки.
export const getInitialTestsState = () => ({
    // Список тестов
    tests: [],
    isListLoading: true,
    listError: "",

    // Флаги удаления тестов
    deletingTests: {},

    // Модальные окна
    attemptModal: getInitialAttemptModal(),
    editModal: getInitialEditModal(),
    addModal: getInitialAddModal(),
});

// Редьюсер списка тестов, обрабатывающий все переходы состояний.
export const testsListReducer = (state, action) => {
    switch (action.type) {
        // --- Загрузка списка тестов ---
        case "FETCH_TESTS_START":
            return {
                ...state,
                isListLoading: true,
                listError: "",
            };
        case "FETCH_TESTS_SUCCESS":
            return {
                ...state,
                isListLoading: false,
                tests: action.payload,
            };
        case "FETCH_TESTS_ERROR":
            return {
                ...state,
                isListLoading: false,
                tests: [],
                listError: action.payload,
            };

        // --- Удаление тестов ---
        case "SET_DELETING_FLAG":
            return {
                ...state,
                deletingTests: {
                    ...state.deletingTests,
                    [action.payload]: true,
                },
            };
        case "CLEAR_DELETING_FLAG": {
            const updated = { ...state.deletingTests };
            delete updated[action.payload];
            return {
                ...state,
                deletingTests: updated,
            };
        }
        case "REMOVE_TEST":
            return {
                ...state,
                tests: state.tests.filter((test) => test.id !== action.payload),
            };

        // --- Обновление тестов ---
        case "UPDATE_TEST":
            return {
                ...state,
                tests: state.tests.map((test) =>
                    test.id === action.payload.id
                        ? { ...test, ...action.payload.updates }
                        : test
                ),
            };
        case "MARK_TEST_COMPLETED":
            return {
                ...state,
                tests: state.tests.map((test) =>
                    test.id === action.payload
                        ? { ...test, isCompleted: true }
                        : test
                ),
            };

        // --- Модалка прохождения теста ---
        case "OPEN_ATTEMPT_MODAL":
            return {
                ...state,
                attemptModal: {
                    open: true,
                    loading: true,
                    error: "",
                    test: action.payload,
                    questions: [],
                    selected: {},
                    submitting: false,
                },
            };
        case "CLOSE_ATTEMPT_MODAL":
            return {
                ...state,
                attemptModal: getInitialAttemptModal(),
            };
        case "SET_ATTEMPT_QUESTIONS":
            return {
                ...state,
                attemptModal: {
                    ...state.attemptModal,
                    loading: false,
                    questions: action.payload,
                },
            };
        case "SET_ATTEMPT_ERROR":
            return {
                ...state,
                attemptModal: {
                    ...state.attemptModal,
                    loading: false,
                    error: action.payload,
                },
            };
        case "TOGGLE_ANSWER": {
            const { questionId, optionIndex, isSingle } = action.payload;
            const selected = { ...state.attemptModal.selected };
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

            return {
                ...state,
                attemptModal: {
                    ...state.attemptModal,
                    selected,
                },
            };
        }
        case "SUBMIT_ATTEMPT_START":
            return {
                ...state,
                attemptModal: {
                    ...state.attemptModal,
                    submitting: true,
                },
            };
        case "SUBMIT_ATTEMPT_END":
            return {
                ...state,
                attemptModal: {
                    ...state.attemptModal,
                    submitting: false,
                },
            };

        // --- Модалка редактирования теста ---
        case "OPEN_EDIT_MODAL":
            return {
                ...state,
                editModal: {
                    open: true,
                    loading: true,
                    saving: false,
                    error: "",
                    testId: action.payload.testId,
                    testName: "",
                    authorsInput: "",
                    description: "",
                    questions: [],
                },
            };
        case "CLOSE_EDIT_MODAL":
            return {
                ...state,
                editModal: getInitialEditModal(),
            };
        case "SET_EDIT_DATA":
            return {
                ...state,
                editModal: {
                    ...state.editModal,
                    loading: false,
                    ...action.payload,
                },
            };
        case "SET_EDIT_ERROR":
            return {
                ...state,
                editModal: {
                    ...state.editModal,
                    loading: false,
                    error: action.payload,
                },
            };
        case "SET_EDIT_SAVING":
            return {
                ...state,
                editModal: {
                    ...state.editModal,
                    saving: action.payload,
                },
            };
        case "UPDATE_EDIT_QUESTIONS": {
            const { index, updater } = action.payload;
            const questions = [...state.editModal.questions];
            const current = questions[index] || createQuestion(index + 1);
            questions[index] =
                typeof updater === "function"
                    ? updater(current)
                    : { ...current, ...updater };
            return {
                ...state,
                editModal: {
                    ...state.editModal,
                    questions,
                },
            };
        }
        case "ADD_EDIT_QUESTION": {
            const nextId =
                (state.editModal.questions[state.editModal.questions.length - 1]
                    ?.id || 0) + 1;
            return {
                ...state,
                editModal: {
                    ...state.editModal,
                    questions: [
                        ...state.editModal.questions,
                        createQuestion(nextId),
                    ],
                },
            };
        }
        case "ADD_EDIT_OPTION": {
            const qIndex = action.payload;
            const questions = [...state.editModal.questions];
            const question = questions[qIndex] || createQuestion(qIndex + 1);
            questions[qIndex] = {
                ...question,
                answerOptions: [...(question.answerOptions || []), ""],
            };
            return {
                ...state,
                editModal: {
                    ...state.editModal,
                    questions,
                },
            };
        }
        case "SET_EDIT_MODAL":
            return {
                ...state,
                editModal: {
                    ...state.editModal,
                    ...action.payload,
                },
            };

        // --- Модалка добавления теста ---
        case "OPEN_ADD_MODAL":
            return {
                ...state,
                addModal: {
                    open: true,
                    saving: false,
                    error: "",
                    form: {
                        testName: "",
                        authorsInput: "",
                        description: "",
                        questions: [createQuestion(1)],
                    },
                },
            };
        case "CLOSE_ADD_MODAL":
            return {
                ...state,
                addModal: getInitialAddModal(),
            };
        case "SET_ADD_SAVING":
            return {
                ...state,
                addModal: {
                    ...state.addModal,
                    saving: action.payload,
                },
            };
        case "UPDATE_ADD_QUESTIONS": {
            const { index, updater } = action.payload;
            const questions = [...state.addModal.form.questions];
            const current = questions[index] || createQuestion(index + 1);
            questions[index] =
                typeof updater === "function"
                    ? updater(current)
                    : { ...current, ...updater };
            return {
                ...state,
                addModal: {
                    ...state.addModal,
                    form: {
                        ...state.addModal.form,
                        questions,
                    },
                },
            };
        }
        case "ADD_QUESTION_TO_FORM": {
            const nextId =
                (state.addModal.form.questions[
                    state.addModal.form.questions.length - 1
                ]?.id || 0) + 1;
            return {
                ...state,
                addModal: {
                    ...state.addModal,
                    form: {
                        ...state.addModal.form,
                        questions: [
                            ...state.addModal.form.questions,
                            createQuestion(nextId),
                        ],
                    },
                },
            };
        }
        case "ADD_OPTION_TO_FORM": {
            const qIndex = action.payload;
            const questions = [...state.addModal.form.questions];
            const question = questions[qIndex] || createQuestion(qIndex + 1);
            questions[qIndex] = {
                ...question,
                answerOptions: [...(question.answerOptions || []), ""],
            };
            return {
                ...state,
                addModal: {
                    ...state.addModal,
                    form: {
                        ...state.addModal.form,
                        questions,
                    },
                },
            };
        }
        case "SET_ADD_MODAL":
            return {
                ...state,
                addModal: {
                    ...state.addModal,
                    ...action.payload,
                },
            };

        default:
            return state;
    }
};
