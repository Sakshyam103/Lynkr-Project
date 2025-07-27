import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface Content {
  id: string;
  userId: string;
  eventId: string;
  mediaType: 'photo' | 'video';
  caption: string;
  tags: string[];
  createdAt: string;
}

interface ContentState {
  content: Content[];
  loading: boolean;
}

const initialState: ContentState = {
  content: [],
  loading: false,
};

const contentSlice = createSlice({
  name: 'content',
  initialState,
  reducers: {
    setContent: (state, action: PayloadAction<Content[]>) => {
      state.content = action.payload;
    },
    addContent: (state, action: PayloadAction<Content>) => {
      state.content.unshift(action.payload);
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
  },
});

export const { setContent, addContent, setLoading } = contentSlice.actions;
export default contentSlice.reducer;