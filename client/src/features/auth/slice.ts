import { createSlice } from '@reduxjs/toolkit';

interface LoggedInState {
  loggedIn: boolean;
}

const initialState: LoggedInState = {
  loggedIn: false,
};

export const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    loggedIn(state) {
      state.loggedIn = true;
    }
  }  
});

export const { loggedIn } = authSlice.actions;
export default authSlice.reducer;
