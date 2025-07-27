/**
 * Redux Store Configuration
 */

import { configureStore } from '@reduxjs/toolkit';
import authSlice from './slices/authSlice';
import eventsSlice from './slices/eventsSlice';
import contentSlice from './slices/contentSlice';

export const store = configureStore({
  reducer: {
    auth: authSlice,
    events: eventsSlice,
    content: contentSlice,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;