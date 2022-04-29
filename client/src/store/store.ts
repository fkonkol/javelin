import create from 'zustand';
import { devtools, persist } from 'zustand/middleware';

interface AuthState {
  loggedIn: boolean;
  login: () => void;
}

const useStore = create<AuthState>()(devtools(persist((set) => ({
  loggedIn: false,
  login: () => set({ loggedIn: true }),
}))));

export default useStore;
