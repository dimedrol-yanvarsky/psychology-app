// Начальное состояние рекомендаций, объединяющее секции, редактирование и флаги загрузки.
export const getInitialRecommendationsState = () => ({
    // Секции рекомендаций
    sections: [],
    currentPage: 0,

    // Редактирование блока
    editingBlock: null, // { id, sectionType, draftText, draftMode }

    // Загрузка данных
    isLoading: true,

    // Флаги pending для операций
    pending: {
        addSection: false,
        addBlock: false,
        deleteSection: false,
        blocks: {}, // { [blockId]: true/false }
    },
});

// Редьюсер рекомендаций, обрабатывающий все переходы состояний.
export const recommendationsReducer = (state, action) => {
    switch (action.type) {
        // --- Загрузка данных ---
        case "LOAD_START":
            return { ...state, isLoading: true };
        case "LOAD_SUCCESS":
            return { ...state, isLoading: false };
        case "LOAD_ERROR":
            return { ...state, isLoading: false };

        // --- Секции и пагинация ---
        case "SET_SECTIONS":
            return { ...state, sections: action.payload };
        case "SET_CURRENT_PAGE":
            return { ...state, currentPage: action.payload };

        // --- Редактирование блока ---
        case "START_EDIT":
            return {
                ...state,
                editingBlock: {
                    id: action.payload.id,
                    sectionType: action.payload.sectionType,
                    draftText: action.payload.draftText,
                    draftMode: action.payload.draftMode,
                },
            };
        case "CANCEL_EDIT":
            return { ...state, editingBlock: null };
        case "SET_EDITING_BLOCK":
            return { ...state, editingBlock: action.payload };

        // --- Флаги pending ---
        case "SET_PENDING_FLAG":
            return {
                ...state,
                pending: {
                    ...state.pending,
                    [action.payload.field]: action.payload.value,
                },
            };
        case "SET_BLOCK_PENDING":
            return {
                ...state,
                pending: {
                    ...state.pending,
                    blocks: {
                        ...state.pending.blocks,
                        [action.payload]: true,
                    },
                },
            };
        case "CLEAR_BLOCK_PENDING": {
            const updatedBlocks = { ...state.pending.blocks };
            delete updatedBlocks[action.payload];
            return {
                ...state,
                pending: {
                    ...state.pending,
                    blocks: updatedBlocks,
                },
            };
        }

        // --- Обновление блока в секции ---
        case "UPDATE_BLOCK_IN_SECTION":
            return {
                ...state,
                sections: state.sections.map((section) =>
                    section.recommendationType === action.payload.sectionType
                        ? {
                              ...section,
                              blocks: section.blocks.map((block) =>
                                  block.id === action.payload.blockId
                                      ? { ...block, ...action.payload.updates }
                                      : block
                              ),
                          }
                        : section
                ),
            };

        default:
            return state;
    }
};
