import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    sections: [],
    currentPage: 0,
    editingBlock: null,
    isLoading: true,
    pending: {
        addSection: false,
        addBlock: false,
        deleteSection: false,
        blocks: {},
    },
};

// Слайс рекомендаций: секции, блоки, редактирование, пагинация.
const recommendationsSlice = createSlice({
    name: "recommendations",
    initialState,
    reducers: {
        // Загрузка
        loadStart(state) {
            state.isLoading = true;
        },
        loadSuccess(state) {
            state.isLoading = false;
        },
        loadError(state) {
            state.isLoading = false;
        },

        // Секции и пагинация
        setSections(state, action) {
            state.sections = action.payload;
        },
        setCurrentPage(state, action) {
            state.currentPage = action.payload;
        },

        // Редактирование блока
        startEdit(state, action) {
            state.editingBlock = {
                id: action.payload.id,
                sectionType: action.payload.sectionType,
                draftText: action.payload.draftText,
                draftMode: action.payload.draftMode,
            };
        },
        cancelEdit(state) {
            state.editingBlock = null;
        },
        setEditingBlock(state, action) {
            state.editingBlock = action.payload;
        },

        // Флаги pending
        setPendingFlag(state, action) {
            state.pending[action.payload.field] = action.payload.value;
        },
        setBlockPending(state, action) {
            state.pending.blocks[action.payload] = true;
        },
        clearBlockPending(state, action) {
            delete state.pending.blocks[action.payload];
        },

        // Обновление блока в секции
        updateBlockInSection(state, action) {
            const { sectionType, blockId, updates } = action.payload;
            const section = state.sections.find(
                (s) => s.recommendationType === sectionType
            );
            if (section) {
                const block = section.blocks.find((b) => b.id === blockId);
                if (block) {
                    Object.assign(block, updates);
                }
            }
        },
    },
});

export const {
    loadStart,
    loadSuccess,
    loadError,
    setSections,
    setCurrentPage,
    startEdit,
    cancelEdit,
    setEditingBlock,
    setPendingFlag,
    setBlockPending,
    clearBlockPending,
    updateBlockInSection,
} = recommendationsSlice.actions;

export default recommendationsSlice.reducer;
