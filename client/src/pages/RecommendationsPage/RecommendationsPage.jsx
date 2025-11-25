import React, { useMemo, useState } from "react";
import styles from "./RecommendationsPage.module.css";

const initialSections = [
    {
        id: "section-1",
        title: "Когда апатия только накрывает",
        blocks: [
            {
                id: "block-1-1",
                text: "Признайте состояние: апатия — не лень, а сигнал о перегрузке. Дайте себе право передохнуть без чувства вины.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-2",
                text: "Мини-движение: выберите действие на 5 минут — проветрить комнату, разложить стол, пройтись до кухни. Важно только начать.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-3",
                text: "Дыхание 4-7-8: вдох на 4, задержка на 7, выдох на 8. Повторите 4 раза, чтобы снизить внутреннее напряжение.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-4",
                text: "Разгрузка инфопотока: отключите уведомления на пару часов, оставьте только музыку или белый шум.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-5",
                text: "Проверка базовых нужд: вода, перекус, короткое растяжение — тело подскажет, если не хватает энергии.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-6",
                text: "Упражнение «3 звука»: закройте глаза и отметьте три звука вокруг, чтобы вернуть внимание к моменту.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-7",
                text: "Микро-список: запишите три микро-задачи и выполните одну прямо сейчас.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-8",
                text: "Точка опоры: прислонитесь спиной к стене на 30 секунд, почувствуйте поддержку.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-9",
                text: "Смена сцены: переместитесь в другое место комнаты, чтобы мозг ощутил новизну.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-10",
                text: "Тёплое тепло: согрейте ладони о чашку чая или тёплую воду — телесный комфорт снижает напряжение.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-11",
                text: "Правило «двух минут»: если действие занимает меньше двух минут, сделайте его сразу.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-12",
                text: "Освободите стол: уберите один предмет, чтобы зрительно стало легче.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-13",
                text: "Пауза на дыхание животом: 6 глубоких вдохов-выдохов медленно, без усилия.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-14",
                text: "Нейтральное движение: мягко покачайтесь из стороны в сторону 30 секунд, чтобы запустить расслабление.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-15",
                text: "Ограничьте сравнения: напомните себе, что темп восстановления — индивидуален.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-16",
                text: "Запланируйте отдых: поставьте таймер на 15 минут полного ничегонеделания, чтобы снять давление.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-17",
                text: "Границы для задач: выберите «достаточно хорошо», а не «идеально», чтобы не сгорать на старте.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-18",
                text: "Свет и воздух: откройте шторы или окно, чтобы добавить дневной свет и кислород.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-19",
                text: "Мягкая музыка: включите трек без слов для фона, чтобы снизить внутренний шум.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-1-20",
                text: "Отметьте успех: после каждого шага фиксируйте «сделано» — даже если шаг крошечный.",
                style: { bold: false, italic: false, underline: false },
            },
        ],
    },
    {
        id: "section-2",
        title: "Возврат энергии и интереса",
        blocks: [
            {
                id: "block-2-1",
                text: "Режим «одно дело в час»: выбирайте одно посильное действие и отмечайте его выполненным. Маленькие победы собирают ощущение контроля.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-2",
                text: "Микро-радости: добавьте 10 минут приятного — любимый плейлист, тёплый душ, вкусный напиток. Пусть тело вспомнит, что ему хорошо.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-3",
                text: "Контакт: напишите одному близкому человеку короткое сообщение без обязательств. Простое «как ты?» помогает вернуть ощущение связи.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-4",
                text: "Щадящий план на неделю: 3 главных задачи и 3 задачи «если будет ресурс». Остальное — на потом. Так мозг понимает, что вы в безопасности.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-5",
                text: "Карта ресурсов: выпишите людей и места, где вам обычно становится чуть легче.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-6",
                text: "Ритуал утра: один повторяемый шаг — стакан воды, умыться, потянуться — закрепляет привычку двигаться.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-7",
                text: "Мини-прогулка: 7–10 минут на улице без телефона, только смотреть вокруг.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-8",
                text: "Техника «5-4-3-2-1»: отметьте 5 предметов, 4 ощущения, 3 звука, 2 запаха, 1 вкус — возвращает в тело.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-9",
                text: "Переключение: если задача стопорится 10 минут, смените её на другую посильную.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-10",
                text: "Поддержка тела: лёгкий углеводный перекус и вода помогают мозгу включиться.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-11",
                text: "Тонкий порядок: разберите один ящик или полку — визуальная ясность снижает внутренний шум.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-12",
                text: "Напоминание «я в процессе»: повесьте карточку с этой фразой, чтобы снять ожидание мгновенного результата.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-13",
                text: "Дневной план: одно важное, одно приятное, одно восстанавливающее — три столпа дня.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-14",
                text: "Гибкий таймбокс: 20 минут работы, 10 минут отдыха — без самокритики, если не успели.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-15",
                text: "Мягкий спорт: растяжка, йога-нидра или прогулка — без задач «сжечь калории».",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-16",
                text: "Социальный мостик: договоритесь о совместном молчаливом созвоне с другом для параллельной работы.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-17",
                text: "Вечерний выдох: запишите три вещи, за которые благодарны телу сегодня.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-18",
                text: "Ритуал завершения дня: свет приглушённый, тёплый душ, отключить уведомления за час до сна.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-19",
                text: "Мета-уровень: замечайте, что апатия — часть процесса восстановления, а не ваша характеристика.",
                style: { bold: false, italic: false, underline: false },
            },
            {
                id: "block-2-20",
                text: "План на завтра: одна посильная цель и время подъёма — больше не нужно.",
                style: { bold: false, italic: false, underline: false },
            },
        ],
    },
];

const defaultStyle = { bold: false, italic: false, underline: false };

const RecommendationsPage = ({ isAdmin, showAlert }) => {
    const [sections, setSections] = useState(initialSections);
    const [currentPage, setCurrentPage] = useState(0);
    const [editingBlock, setEditingBlock] = useState(null);

    const currentSection = useMemo(
        () => sections[currentPage] || sections[0],
        [sections, currentPage]
    );

    const createId = (prefix) =>
        `${prefix}-${Date.now().toString(36)}-${Math.random()
            .toString(36)
            .slice(2, 6)}`;

    const handleBlockClick = (sectionId, block) => {
        if (!isAdmin) return;
        setEditingBlock({
            sectionId,
            blockId: block.id,
            draftText: block.text,
            draftStyle: block.style || defaultStyle,
        });
    };

    const handleToggleStyle = (key) => {
        if (!editingBlock) return;
        setEditingBlock((prev) => ({
            ...prev,
            draftStyle: {
                ...prev.draftStyle,
                [key]: !prev.draftStyle?.[key],
            },
        }));
    };

    const handleSaveBlock = () => {
        if (!editingBlock) return;

        setSections((prev) =>
            prev.map((section) =>
                section.id === editingBlock.sectionId
                    ? {
                          ...section,
                          blocks: section.blocks.map((block) =>
                              block.id === editingBlock.blockId
                                  ? {
                                        ...block,
                                        text: editingBlock.draftText,
                                        style: editingBlock.draftStyle,
                                    }
                                  : block
                          ),
                      }
                    : section
            )
        );
        setEditingBlock(null);
        if (showAlert) {
            showAlert("success", "Блок обновлён");
        }
    };

    const handleAddSection = () => {
        const newSection = {
            id: createId("section"),
            title: "Новый раздел рекомендаций",
            blocks: [
                {
                    id: createId("block"),
                    text: "Начните с описания нового шага или практики, полезной при апатии.",
                    style: defaultStyle,
                },
            ],
        };

        setSections((prev) => [...prev, newSection]);
        setCurrentPage(sections.length);
        setEditingBlock(null);
        if (showAlert) {
            showAlert("success", "Добавлен новый раздел");
        }
    };

    const handleAddBlock = () => {
        const sectionId = currentSection?.id;
        if (!sectionId) return;
        const newBlock = {
            id: createId("block"),
            text: "Новый текстовый блок — добавьте конкретное действие или мысль поддержки.",
            style: defaultStyle,
        };

        setSections((prev) =>
            prev.map((section) =>
                section.id === sectionId
                    ? { ...section, blocks: [...section.blocks, newBlock] }
                    : section
            )
        );
        setEditingBlock(null);
        if (showAlert) {
            showAlert("success", "Блок добавлен");
        }
    };

    const handleRemoveBlock = () => {
        const sectionId = currentSection?.id;
        if (!sectionId) return;

        setSections((prev) =>
            prev.map((section) => {
                if (section.id !== sectionId) return section;
                if (section.blocks.length === 0) return section;
                const updatedBlocks = section.blocks.slice(
                    0,
                    section.blocks.length - 1
                );
                return { ...section, blocks: updatedBlocks };
            })
        );
        setEditingBlock(null);
    };

    const handleRemoveSection = () => {
        if (sections.length <= 1) return;
        const updated = sections.filter((_, index) => index !== currentPage);
        const nextIndex = Math.max(0, currentPage - 1);
        setSections(updated);
        setCurrentPage(nextIndex);
        setEditingBlock(null);
    };

    const blockStyle = (style) => ({
        fontWeight: style?.bold ? 700 : 500,
        fontStyle: style?.italic ? "italic" : "normal",
        textDecoration: style?.underline ? "underline" : "none",
    });

    return (
        <div className={styles.page}>
            <header className={styles.topBar}>
                <div className={styles.breadcrumb}>
                    <span className={styles.overline}>Recommendations</span>
                    <h1 className={styles.pageTitle}>Рекомендации</h1>
                    <p className={styles.pageSubtitle}>
                        Два набора коротких шагов для человека, который
                        застрял в апатии.
                    </p>
                </div>

                {isAdmin && (
                    <div className={styles.topBarActions}>
                        <button
                            type="button"
                            className={`${styles.actionButton} ${styles.secondaryAction}`}
                            onClick={handleAddSection}
                        >
                            Добавить раздел
                        </button>
                        <button
                            type="button"
                            className={`${styles.actionButton} ${styles.secondaryAction}`}
                            onClick={handleAddBlock}
                        >
                            Добавить блок
                        </button>
                        <button
                            type="button"
                            className={`${styles.actionButton} ${styles.dangerAction}`}
                            onClick={handleRemoveSection}
                            disabled={sections.length <= 1}
                        >
                            Удалить раздел
                        </button>
                        <button
                            type="button"
                            className={`${styles.actionButton} ${styles.dangerAction}`}
                            onClick={handleRemoveBlock}
                            disabled={!currentSection || currentSection.blocks.length === 0}
                        >
                            Удалить блок
                        </button>
                    </div>
                )}
            </header>

            <main className={styles.main}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionBadge}>Страница {currentPage + 1}</div>
                    <h2 className={styles.sectionTitle}>
                        {currentSection?.title || "Новый раздел"}
                    </h2>
                </div>

                <div className={styles.blocksGrid}>
                    {currentSection?.blocks?.map((block) => {
                        const isEditing =
                            editingBlock &&
                            editingBlock.blockId === block.id &&
                            editingBlock.sectionId === currentSection.id;
                        return (
                            <article
                                key={block.id}
                                className={`${styles.blockCard} ${
                                    isAdmin ? styles.blockCardAdmin : ""
                                } ${isEditing ? styles.blockEditing : ""}`}
                                onClick={() =>
                                    handleBlockClick(currentSection.id, block)
                                }
                            >
                                {!isEditing ? (
                                    <>
                                        <p
                                            className={styles.blockText}
                                            style={blockStyle(block.style)}
                                        >
                                            {block.text}
                                        </p>
                                        {isAdmin && (
                                            <span
                                                className={styles.editHint}
                                                aria-hidden="true"
                                            >
                                                редактировать
                                            </span>
                                        )}
                                    </>
                                ) : (
                                    <div
                                        className={styles.editorWrapper}
                                        onClick={(e) => e.stopPropagation()}
                                    >
                                        <div className={styles.toolbar}>
                                            <button
                                                type="button"
                                                className={`${styles.toolButton} ${
                                                    editingBlock.draftStyle?.bold
                                                        ? styles.toolActive
                                                        : ""
                                                }`}
                                                onClick={() =>
                                                    handleToggleStyle("bold")
                                                }
                                            >
                                                Ж
                                            </button>
                                            <button
                                                type="button"
                                                className={`${styles.toolButton} ${
                                                    editingBlock.draftStyle
                                                        ?.italic
                                                        ? styles.toolActive
                                                        : ""
                                                }`}
                                                onClick={() =>
                                                    handleToggleStyle("italic")
                                                }
                                            >
                                                К
                                            </button>
                                            <button
                                                type="button"
                                                className={`${styles.toolButton} ${
                                                    editingBlock.draftStyle
                                                        ?.underline
                                                        ? styles.toolActive
                                                        : ""
                                                }`}
                                                onClick={() =>
                                                    handleToggleStyle(
                                                        "underline"
                                                    )
                                                }
                                            >
                                                U
                                            </button>
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
                                            >
                                                Сохранить
                                            </button>
                                            <button
                                                type="button"
                                                className={`${styles.actionButton} ${styles.mutedAction}`}
                                                onClick={() =>
                                                    setEditingBlock(null)
                                                }
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
            </main>

            <footer className={styles.pagination}>
                {sections.map((section, index) => (
                    <button
                        key={section.id}
                        type="button"
                        className={`${styles.pageButton} ${
                            index === currentPage ? styles.pageButtonActive : ""
                        }`}
                        onClick={() => setCurrentPage(index)}
                    >
                        {index + 1}
                    </button>
                ))}
            </footer>
        </div>
    );
};

export default RecommendationsPage;
