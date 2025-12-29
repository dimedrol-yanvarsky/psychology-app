import React from "react";
import styles from "../DashboardPage.module.css";

const AnswersQuestion = ({ question, questionIndex, selectedAnswers }) => {
    const selectType = question.selectype || question.selectType || "";
    const isSingle = selectType === "one";
    const questionNumber = question.id || question.number || questionIndex + 1;
    const selected = selectedAnswers.get(Number(questionNumber)) || new Set();

    const rawQuestionText =
        question.questionBody || question.question || question.body || "";
    const questionText =
        typeof rawQuestionText === "string"
            ? rawQuestionText
            : rawQuestionText?.body ||
              rawQuestionText?.text ||
              rawQuestionText?.title ||
              "Вопрос";

    const options = Array.isArray(question.answerOptions)
        ? question.answerOptions
        : Array.isArray(question.answers)
        ? question.answers
        : [];

    return (
        <div className={styles.questionCard}>
            <div className={styles.questionTitle}>
                {questionNumber}. {questionText}
            </div>
            <div className={styles.optionsList}>
                {options.map((option, optionIndex) => {
                    const optionNumber = optionIndex + 1;
                    const optionLabel =
                        typeof option === "string"
                            ? option
                            : option?.body ||
                              option?.text ||
                              option?.title ||
                              option?.label ||
                              String(optionNumber);
                    const isChecked =
                        selected instanceof Set &&
                        selected.has(optionNumber);
                    const inputType = isSingle ? "radio" : "checkbox";

                    return (
                        <label
                            key={`${question.id || questionIndex}-${
                                option.id || optionIndex
                            }`}
                            className={`${styles.optionItem} ${
                                isChecked ? styles.optionSelected : ""
                            }`}
                        >
                            <input
                                type={inputType}
                                name={`question-${questionNumber}`}
                                checked={isChecked}
                                disabled
                                readOnly
                            />
                            <span>{optionLabel}</span>
                        </label>
                    );
                })}
            </div>
        </div>
    );
};

export default AnswersQuestion;
