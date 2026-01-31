export const createQuestion = (id = 1) => ({
    id,
    questionBody: "",
    answerOptions: ["", ""],
    selectType: "one",
});

export const getOptionValue = (option) =>
    typeof option === "string"
        ? option
        : option?.body || option?.text || option?.title || "";

const getOptionId = (option) => {
    if (option && typeof option === "object") {
        const candidate =
            option.id || option.ID || option.optionId || option.optionID;
        const id = Number(candidate);
        if (Number.isFinite(id) && id > 0) {
            return id;
        }
    }
    return null;
};

export const normalizeAnswerOptions = (options) => {
    if (!Array.isArray(options)) {
        return ["", ""];
    }
    const normalized = options.map(getOptionValue);
    return normalized.length ? normalized : ["", ""];
};

const normalizeOptionsForSave = (options) => {
    if (!Array.isArray(options)) {
        return [];
    }

    const normalized = [];
    let nextId = 1;

    for (let index = 0; index < options.length; index += 1) {
        const option = options[index];
        const body = getOptionValue(option).trim();
        if (!body) {
            continue;
        }

        const incomingId = getOptionId(option);
        const id = incomingId || nextId;
        nextId = Math.max(nextId, id + 1);

        normalized.push({ id, body });
    }

    return normalized;
};

export const normalizeAuthorsInput = (value) =>
    value
        .split(",")
        .map((item) => item.trim())
        .filter(Boolean);

export const normalizeQuestionsForSave = (questions) => {
    const normalized = [];

    for (let index = 0; index < questions.length; index += 1) {
        const question = questions[index] || {};
        const qBody = (question.questionBody || "").trim();

        if (!qBody) {
            return {
                error: "Заполните формулировку каждого вопроса",
                questions: [],
            };
        }

        const options = Array.isArray(question.answerOptions)
            ? normalizeOptionsForSave(question.answerOptions)
            : [];

        if (options.length === 0) {
            return {
                error: "Добавьте варианты ответа для каждого вопроса",
                questions: [],
            };
        }

        const selectType =
            (question.selectType || question.selectype || "").trim() || "one";
        const id =
            Number(question.id || question.number || index + 1) || index + 1;

        normalized.push({
            id,
            questionBody: qBody,
            answerOptions: options,
            selectType,
        });
    }

    return { questions: normalized, error: "" };
};

export const getQuestionNumber = (question, fallbackIndex = 0) =>
    Number(
        question?.id ||
            question?.number ||
            question?.questionNumber ||
            fallbackIndex + 1
    );
