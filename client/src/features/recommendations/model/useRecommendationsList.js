import { useCallback, useEffect, useMemo, useState } from "react";
import {
    addBlock,
    addSection,
    deleteBlock,
    deleteSection,
    fetchRecommendations,
    updateBlock,
} from "../../../entities/recommendation";

const TEXT_MODE_OPTIONS = [
    { value: "base", label: "Aa", title: "Обычный стиль" },
    { value: "bold", label: "Ж", title: "Жирный текст" },
    { value: "line", label: "⎯", title: "Подчеркнутый текст" },
    {
        value: "bold-italics-line",
        label: "ЖК⎯",
        title: "Жирный курсив с линией",
    },
];

const DEFAULT_BLOCK_TEXT =
    "Новый текстовый блок — добавьте конкретное действие или мысль поддержки.";

const getPageNumber = (recommendationType = "") => {
    const match = `${recommendationType}`.match(/\d+/);
    return match ? Number(match[0]) : 0;
};

const normalizeRecommendation = (item = {}) => {
    const text =
        (item.recommendationText || item.text || DEFAULT_BLOCK_TEXT).trim() ||
        DEFAULT_BLOCK_TEXT;
    const recommendationType =
        (item.recommendationType || item.type || "Страница 1").trim() ||
        "Страница 1";
    const textMode = (item.textMode || "base").trim() || "base";
    const id = (item._id || item.id || "").toString();

    return {
        id,
        recommendationText: text,
        textMode,
        recommendationType,
    };
};

const groupRecommendations = (list = []) => {
    const map = new Map();

    list.forEach((item) => {
        const rec = normalizeRecommendation(item);
        const key = rec.recommendationType;

        if (!map.has(key)) {
            map.set(key, { recommendationType: key, blocks: [] });
        }
        map.get(key).blocks.push(rec);
    });

    const sections = Array.from(map.values());

    sections.sort((a, b) => {
        const diff =
            getPageNumber(a.recommendationType) -
            getPageNumber(b.recommendationType);
        return diff !== 0
            ? diff
            : a.recommendationType.localeCompare(b.recommendationType);
    });

    sections.forEach((section) =>
        section.blocks.sort((a, b) => a.id.localeCompare(b.id))
    );

    return sections;
};

const blockStyle = (mode = "base") => {
    const trimmed = mode.trim();
    switch (trimmed) {
        case "bold":
            return { fontWeight: 700 };
        case "line":
            return { textDecoration: "underline" };
        case "bold-italics-line":
            return {
                fontWeight: 700,
                fontStyle: "italic",
                textDecoration: "underline",
                color: "var(--green-primary)",
            };
        default:
            return {};
    }
};

export const useRecommendationsList = ({ isAdmin, showAlert }) => {
    const [sections, setSections] = useState([]);
    const [currentPage, setCurrentPage] = useState(0);
    const [editingBlock, setEditingBlock] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [pending, setPending] = useState({
        addSection: false,
        addBlock: false,
        deleteSection: false,
        blocks: {},
    });

    const currentSection = useMemo(
        () => sections[currentPage] || null,
        [sections, currentPage]
    );

    const setPendingFlag = (key, value) =>
        setPending((prev) => ({ ...prev, [key]: value }));

    const setBlockPending = useCallback((blockId, value) => {
        if (!blockId) return;
        setPending((prev) => {
            const nextBlocks = { ...prev.blocks };
            if (value) {
                nextBlocks[blockId] = true;
            } else {
                delete nextBlocks[blockId];
            }
            return { ...prev, blocks: nextBlocks };
        });
    }, []);

    const isBlockPending = (blockId) => Boolean(pending.blocks[blockId]);

    const applyRecommendations = useCallback((list, preferredType = "") => {
        const grouped = groupRecommendations(list);
        setSections(grouped);

        if (grouped.length === 0) {
            setCurrentPage(0);
            return;
        }

        if (preferredType) {
            const idx = grouped.findIndex(
                (section) => section.recommendationType === preferredType
            );
            if (idx !== -1) {
                setCurrentPage(idx);
                return;
            }
        }

        setCurrentPage((prev) =>
            Math.min(prev, Math.max(grouped.length - 1, 0))
        );
    }, []);

    const loadRecommendations = useCallback(async () => {
        setIsLoading(true);
        try {
            const { data } = await fetchRecommendations();
            const list = Array.isArray(data?.recommendations)
                ? data.recommendations
                : [];
            applyRecommendations(list);
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                "Не удалось загрузить рекомендации";
            showAlert?.("error", message);
            setSections([]);
        } finally {
            setIsLoading(false);
        }
    }, [applyRecommendations, showAlert]);

    useEffect(() => {
        loadRecommendations();
    }, [loadRecommendations]);

    const handleStartEdit = (block, sectionType) => {
        if (!isAdmin) return;
        setEditingBlock({
            id: block.id,
            sectionType,
            draftText: block.recommendationText,
            draftMode: block.textMode || "base",
        });
    };

    const handleSaveBlock = async () => {
        if (!editingBlock) return;

        const text = editingBlock.draftText.trim();
        if (!text) {
            showAlert?.("error", "Заполните текст рекомендации");
            return;
        }

        setBlockPending(editingBlock.id, true);
        try {
            const { data } = await updateBlock({
                id: editingBlock.id,
                recommendationText: text,
                textMode: editingBlock.draftMode || "base",
            });

            const updated = normalizeRecommendation(
                data?.block || {
                    _id: editingBlock.id,
                    recommendationText: text,
                    recommendationType: editingBlock.sectionType,
                    textMode: editingBlock.draftMode,
                }
            );

            setSections((prev) =>
                prev.map((section) =>
                    section.recommendationType === editingBlock.sectionType
                        ? {
                              ...section,
                              blocks: section.blocks.map((block) =>
                                  block.id === editingBlock.id
                                      ? {
                                            ...block,
                                            recommendationText:
                                                updated.recommendationText,
                                            textMode: updated.textMode,
                                        }
                                      : block
                              ),
                          }
                        : section
                )
            );

            showAlert?.("success", data?.message || "Блок обновлен");
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                "Не удалось обновить блок";
            showAlert?.("error", message);
        } finally {
            setBlockPending(editingBlock.id, false);
            setEditingBlock(null);
        }
    };

    const handleAddSection = async () => {
        if (!isAdmin) return;

        setPendingFlag("addSection", true);
        try {
            const { data } = await addSection();
            const list = Array.isArray(data?.recommendations)
                ? data.recommendations
                : [];
            const newType =
                data?.block?.recommendationType ||
                `Страница ${sections.length + 1}`;
            applyRecommendations(list, newType);
            showAlert?.("success", data?.message || "Добавлен новый раздел");
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                "Не удалось добавить раздел";
            showAlert?.("error", message);
        } finally {
            setPendingFlag("addSection", false);
            setEditingBlock(null);
        }
    };

    const handleAddBlock = async () => {
        if (!isAdmin || !currentSection) return;

        setPendingFlag("addBlock", true);
        try {
            const { data } = await addBlock({
                recommendationType: currentSection.recommendationType,
                recommendationText: DEFAULT_BLOCK_TEXT,
                textMode: "base",
            });

            const list = Array.isArray(data?.recommendations)
                ? data.recommendations
                : [];
            applyRecommendations(list, currentSection.recommendationType);
            showAlert?.("success", data?.message || "Блок добавлен");
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                "Не удалось добавить блок";
            showAlert?.("error", message);
        } finally {
            setPendingFlag("addBlock", false);
            setEditingBlock(null);
        }
    };

    const handleDeleteBlock = async (block, sectionType) => {
        if (!isAdmin || !block?.id) return;

        setBlockPending(block.id, true);
        try {
            const { data } = await deleteBlock({
                id: block.id,
            });

            const list = Array.isArray(data?.recommendations)
                ? data.recommendations
                : [];
            applyRecommendations(list, sectionType);
            showAlert?.("success", data?.message || "Блок удален");
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                "Не удалось удалить блок";
            showAlert?.("error", message);
        } finally {
            setBlockPending(block.id, false);
            if (editingBlock?.id === block.id) {
                setEditingBlock(null);
            }
        }
    };

    const handleDeleteSection = async () => {
        if (!isAdmin || !currentSection) return;

        setPendingFlag("deleteSection", true);
        const nextPageNumber = Math.max(
            1,
            getPageNumber(currentSection.recommendationType) - 1
        );
        const preferredType = `Страница ${nextPageNumber}`;

        try {
            const { data } = await deleteSection({
                recommendationType: currentSection.recommendationType,
            });

            const list = Array.isArray(data?.recommendations)
                ? data.recommendations
                : [];
            applyRecommendations(list, preferredType);
            showAlert?.("success", data?.message || "Раздел удален");
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                "Не удалось удалить раздел";
            showAlert?.("error", message);
        } finally {
            setPendingFlag("deleteSection", false);
            setEditingBlock(null);
        }
    };

    const handleCancelEdit = () => setEditingBlock(null);

    const handlePageChange = (index) => {
        setEditingBlock(null);
        setCurrentPage(index);
    };

    return {
        DEFAULT_BLOCK_TEXT,
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
    };
};
