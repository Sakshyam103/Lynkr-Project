import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface Event {
  id: string;
  name: string;
  description: string;
  startDate: string;
  endDate: string;
  location: {
    latitude: number;
    longitude: number;
    address: string;
  };
  brandId: string;
}

interface EventsState {
  events: Event[];
  currentEvent: Event | null;
  loading: boolean;
}

const initialState: EventsState = {
  events: [],
  currentEvent: null,
  loading: false,
};

const eventsSlice = createSlice({
  name: 'events',
  initialState,
  reducers: {
    setEvents: (state, action: PayloadAction<Event[]>) => {
      state.events = action.payload;
    },
    setCurrentEvent: (state, action: PayloadAction<Event>) => {
      state.currentEvent = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
  },
});

export const { setEvents, setCurrentEvent, setLoading } = eventsSlice.actions;
export default eventsSlice.reducer;