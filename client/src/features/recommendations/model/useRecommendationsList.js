import { useCallback, useEffect, useMemo } from "react";
import { useSelector, useDispatch } from "react-redux";
import {
    addBlock,
    addSection,
    deleteBlock,
    deleteSection,
    fetchRecommendations,
    updateBlock,
} from "../../../entities/recommendation";
import { useAuthContext } from "../../../shared/context/AuthContext";
import { useAlertContext } from "../../../shared/context/AlertContext";
import {
    loadStart,
    loadSuccess,
    setSections,
    setCurrentPage,
    startEdit,
    cancelEdit,
    setEditingBlock,
    setPendingFlag,
    setBlockPending,
    clearBlockPending,
    updateBlockInSection,
} from "./recommendationsSlice";

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

export const useRecommendationsList = () => {
    const { isAdmin } = useAuthContext();
    const { showAlert } = useAlertContext();
    const reduxDispatch = useDispatch();

    const state = useSelector((s) => s.recommendations);

    const currentSection = useMemo(
        () => state.sections[state.currentPage] || null,
        [state.sections, state.currentPage]
    );

    const isBlockPending = (blockId) =>
        Boolean(state.pending.blocks[blockId]);

    const applyRecommendations = useCallback(
        (list, preferredType = "") => {
            const grouped = groupRecommendations(list);
            reduxDispatch(setSections(grouped));

            if (grouped.length === 0) {
                reduxDispatch(setCurrentPage(0));
                return;
            }

            if (preferredType) {
                const idx = grouped.findIndex(
                    (section) =>
                        section.recommendationType === preferredType
                );
                if (idx !== -1) {
                    reduxDispatch(setCurrentPage(idx));
                    return;
                }
            }

            reduxDispatch(
                setCurrentPage(
                    Math.min(
                        state.currentPage,
                        Math.max(grouped.length - 1, 0)
                    )
                )
            );
        },
        [state.currentPage, reduxDispatch]
    );

    const loadRecommendations = useCallback(async () => {
        reduxDispatch(loadStart());
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
            reduxDispatch(setSections([]));
        } finally {
            reduxDispatch(loadSuccess());
        }
    }, [applyRecommendations, showAlert, reduxDispatch]);

    useEffect(() => {
        loadRecommendations();
    }, [loadRecommendations]);

    const handleStartEdit = (block, sectionType) => {
        if (!isAdmin) return;
        reduxDispatch(
            startEdit({
                id: block.id,
                sectionType,
                draftText: block.recommendationText,
                draftMode: block.textMode || "base",
            })
        );
    };

    const handleSaveBlock = async () => {
        if (!state.editingBlock) return;

        const text = state.editingBlock.draftText.trim();
        if (!text) {
            showAlert?.("error", "Заполните текст рекомендации");
            return;
        }

        reduxDispatch(setBlockPending(state.editingBlock.id));
        try {
            const { data } = await updateBlock({
                id: state.editingBlock.id,
                recommendationText: text,
                textMode: state.editingBlock.draftMode || "base",
            });

            const updated = normalizeRecommendation(
                data?.block || {
                    _id: state.editingBlock.id,
                    recommendationText: text,
                    recommendationType: state.editingBlock.sectionType,
                    textMode: state.editingBlock.draftMode,
                }
            );

            reduxDispatch(
                updateBlockInSection({
                    sectionType: state.editingBlock.sectionType,
                    blockId: state.editingBlock.id,
                    updates: {
                        recommendationText: updated.recommendationText,
                        textMode: updated.textMode,
                    },
                })
            );

            showAlert?.("success", data?.message || "Блок обновлен");
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                "Не удалось обновить блок";
            showAlert?.("error", message);
        } finally {
            reduxDispatch(clearBlockPending(state.editingBlock.id));
            reduxDispatch(cancelEdit());
        }
    };

    const handleAddSection = async () => {
        if (!isAdmin) return;

        reduxDispatch(
            setPendingFlag({ field: "addSection", value: true })
        );
        try {
            const { data } = await addSection();
            const list = Array.isArray(data?.recommendations)
                ? data.recommendations
                : [];
            const newType =
                data?.block?.recommendationType ||
                `Страница ${state.sections.length + 1}`;
            applyRecommendations(list, newType);
            showAlert?.("success", data?.message || "Добавлен новый раздел");
        } catch (error) {
            const message =
                error?.response?.data?.message ||
                "Не удалось добавить раздел";
            showAlert?.("error", message);
        } finally {
            reduxDispatch(
                setPendingFlag({ field: "addSection", value: false })
            );
            reduxDispatch(cancelEdit());
        }
    };

    const handleAddBlock = async () => {
        if (!isAdmin || !currentSection) return;

        reduxDispatch(
            setPendingFlag({ field: "addBlock", value: true })
        );
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
            reduxDispatch(
                setPendingFlag({ field: "addBlock", value: false })
            );
            reduxDispatch(cancelEdit());
        }
    };

    const handleDeleteBlock = async (block, sectionType) => {
        if (!isAdmin || !block?.id) return;

        reduxDispatch(setBlockPending(block.id));
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
            reduxDispatch(clearBlockPending(block.id));
            if (state.editingBlock?.id === block.id) {
                reduxDispatch(cancelEdit());
            }
        }
    };

    const handleDeleteSection = async () => {
        if (!isAdmin || !currentSection) return;

        reduxDispatch(
            setPendingFlag({ field: "deleteSection", value: true })
        );
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
            reduxDispatch(
                setPendingFlag({ field: "deleteSection", value: false })
            );
            reduxDispatch(cancelEdit());
        }
    };

    const handleCancelEdit = () => reduxDispatch(cancelEdit());

    const handlePageChange = (index) => {
        reduxDispatch(cancelEdit());
        reduxDispatch(setCurrentPage(index));
    };

    return {
        DEFAULT_BLOCK_TEXT,
        TEXT_MODE_OPTIONS,
        blockStyle,
        currentPage: state.currentPage,
        currentSection,
        editingBlock: state.editingBlock,
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
        isLoading: state.isLoading,
        pending: state.pending,
        sections: state.sections,
        setEditingBlock: (payload) => reduxDispatch(setEditingBlock(payload)),
    };
};
